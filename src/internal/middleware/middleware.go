package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.gin.order/pkg/encrypt"
	"go.gin.order/pkg/msg"
	"go.gin.order/pkg/token"
	"go.gin.order/pkg/util"
	"go.gin.order/src/config"
	casbin2 "go.gin.order/src/config/casbin"
	"go.gin.order/src/config/dto"
	"go.gin.order/src/config/logger"
	"go.gin.order/src/config/redis"
	"golang.org/x/time/rate"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"runtime/debug"
	"time"
)

var (
	tokens = token.Token{}
	rsa    = encrypt.NewRSAEncrypt()
	casbin = casbin2.CasbinService{}
)

/**
 * @Author Khr
 * @Description //casbin鉴权
 * @Date 15:51 2024/7/23
 * @Param
 * @return
 **/
func CasbinEnforce() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestUrl := c.Request.URL.String()
		method := c.Request.Method
		usertype := c.GetString("user_type")
		permission, err := casbin.CheckPermission(usertype, requestUrl, method)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "message": msg.CASBINFOUNDERROR})
			return
		}
		if permission {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "message": msg.CASBINFOUNDERROR})
			return
		}
		c.Next()
	}
}

// 鉴权中间件
func TokenMiddleware() gin.HandlerFunc {
	// 在此处进行鉴权逻辑
	// 如果鉴权失败，则返回 401 Unauthorized
	return func(c *gin.Context) {
		fmt.Println("token认证开始执行")
		//t := time.Now()
		requestUrl := c.Request.URL.String()
		//路径模糊匹配
		if !util.FuzzyMatch(requestUrl, config.WhiteUrl) {
			//请求头是否携带token
			tokenExist := c.GetHeader("Authorization")
			if len(tokenExist) < 0 || tokenExist == "" {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "message": msg.FoundTokenError})
				return
			}
			claims := tokens.ParseToken(c.GetHeader("Authorization"))
			c.Set("claims", claims)
			c.Set("user_type", claims.Role)
		}
		c.Next()
		//fmt.Println("中间件执行结束")
	}
}

func VerifyCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestUrl := c.Request.URL.String()
		log.Println(requestUrl, "请求路径")
		if !util.FuzzyMatch(requestUrl, config.WhiteUrl) {
			log.Println("cookie验证")
			cookie, err := c.Cookie("name")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "message": msg.COOKIEEXPRITIMEERROR})
					return
				}
			}
			c.Set("user", cookie)
			c.SetCookie("user", cookie, 3600, "/", "localhost", false, true)
			//log.Println(cookie, "打印cookie")
			c.Next()
		}
	}
}

func LimitCaptcha() gin.HandlerFunc {
	return func(c *gin.Context) {
		redisStore := redis.NewRedisClient()
		ip := c.ClientIP()
		key := fmt.Sprintf("limiter:%s", ip)
		count, _ := redisStore.AutoInc(key)
		//log.Println(e, "自增错误")
		if count > 5 {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusTooManyRequests, "message": msg.CAPTCHAERROR})
			return
		}
		redisStore.ExpireRedis(key, time.Minute)
		c.Next()
	}
}

/**
 * @Author Khr
 * @Description // 熔断中间件
 * @Date 14:53 2024/2/21
 * @Param
 * @return
 **/
func CircuitBreakerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		hystrix.ConfigureCommand("service1", hystrix.CommandConfig{
			Timeout:                1000,
			MaxConcurrentRequests:  100,
			ErrorPercentThreshold:  25,
			RequestVolumeThreshold: 10,
			SleepWindow:            5000,
		})
		err := hystrix.Do("service1", func() error {
			// 在此处调用后端服务
			// 如果调用失败，则返回错误
			return nil
		}, nil)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "message": "Service unavailable"})
			return
		}
	}
}

/*
* @MethodName
* @Description 其他异常捕捉
* @Author khr
* @Date 2023/7/31 15:21
 */

func GlobalErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			//fmt.Println("程序捕捉错误")
			if r := recover(); r != nil {
				fmt.Println("打印错误信息:", r)
				// 打印错误堆栈信息
				debug.PrintStack()
				err := c.Errors.Last()
				errorMessage := err.Err.Error()
				c.JSON(http.StatusInternalServerError, &dto.Res{
					Code:    http.StatusInternalServerError,
					Message: errorMessage,
					Data:    nil,
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}

/*
* @MethodName UnifiedResponseMiddleware
* @Description 统一返回正确和错误格式
* @Author khr
* @Date 2023/7/29 9:45
 */

func UnifiedResponseMiddleware() gin.HandlerFunc {
	//log.Println("使用格式化中间件")
	return func(c *gin.Context) {

		// 继续处理请求
		c.Next()
		//log.Println("进入统一格式化中间件，处理返回结构体")
		//logger.Println("************** 处理请求之后")
		if len(c.Errors) > 0 {
			fmt.Println("出现错误", c.Errors)
			err := c.Errors.Last()
			errorMessage := err.Err.Error()
			fmt.Println(c.Writer.Status(), "错误类型")
			c.JSON(http.StatusOK, &dto.Res{
				Code:    c.Writer.Status(),
				Message: errorMessage,
				Data:    "",
			})
			return

		}
		//log.Println("------------------")
		if c.Writer.Status() == http.StatusOK {
			data, exists := c.Get("res")
			if exists {
				c.JSON(http.StatusOK, &dto.Res{
					Code:    http.StatusOK,
					Message: "success",
					Data:    data,
				})
			} else {
				// 如果没有设置自定义响应数据，则处理默认响应
				if !c.Writer.Written() {
					c.JSON(http.StatusOK, &dto.Res{
						Code:    c.Writer.Status(),
						Message: "success",
						Data:    nil,
					})
				}
			}
		}

	}
}

// 限流中间件
func RateLimiterMiddleware() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second), 100) // 每秒最多处理100个请求

	return func(c *gin.Context) {
		if limiter.Allow() == false {
			c.AbortWithError(http.StatusTooManyRequests, errors.New("too Many Requests"))
			//http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}

// 限流中间件
//func rateLimitMiddleware() gin.HandlerFunc {
//	// 在此处进行限流逻辑
//	return func(c *gin.Context) {
//		ip := c.ClientIP()
//		if ip == "" {
//			ip = c.Request.RemoteAddr
//		}
//		if utils.ValidateExist(ip, global.IpAccess) {
//			c.Next()
//		}
//		path := c.Request.URL.Path
//		//fmt.Println(ip, path)
//
//		// 组合出 key
//		key := fmt.Sprintf("request:%s:%s", ip, path)
//		//fmt.Print("key", key)
//		// 将请求次数 +1，并设置过期时间
//		err := global.AutoInc(key)
//
//		if err != nil {
//			// 记录日志
//			fmt.Println("incr error:", err)
//			c.AbortWithError(http.StatusInternalServerError, err)
//			return
//		}
//		if err = global.ExpireRedis(key, time.Hour); err != nil {
//			logger.Printf("redis缓存失败：%s", err)
//			c.AbortWithError(http.StatusInternalServerError, err)
//			return
//		}
//
//		// 获取当前IP在 path 上的请求次数
//		accessTime := global.GetLimitRedis(key)
//
//		if err != nil {
//			// 记录日志
//			fmt.Println("get error:", err)
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//		//ip一小时内访问路径超过次数限制，拒绝访问
//		if accessTime > 60 {
//			requestLimit := fmt.Sprintf("request:%s:%s", ip, path)
//			if err = global.RpushRedis(global.InterceptPrefix, requestLimit); err != nil {
//				c.AbortWithStatusJSON(http.StatusServiceUnavailable, err)
//				return
//			}
//			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
//			return
//		}
//		mu.Lock()
//		_, ok := visitorMap[ip]
//		var limiter = rate.NewLimiter(1, 10) // 设置限制为1个请求/秒，最多允许10个并发请求
//		// 如果该IP地址不存在，则创建一个速率限制器
//		if !ok {
//			visitorMap[ip] = limiter
//		}
//		mu.Unlock()
//		// 尝试获取令牌，如果没有可用的令牌则阻塞
//		if !limiter.Allow() {
//			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
//			return
//		}
//	}
//}

/*
* @MethodName
* @Description 日志中间件
* @Author khr
* @Date 2023/7/31 15:19
 */

func LoggerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		//requestBody, _ := c.Get("reqbody")
		//rbody, _ := requestBody.(string)
		query := c.Request.URL.RawQuery
		c.Next() // 调用该请求的剩余处理程序
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds()/1000000))))
		//logBody := &dto.LogBody{
		//	SpendTime: spendTime,
		//	Path:      c.Request.RequestURI,
		//	Method:    c.Request.Method,
		//	Status:    c.Writer.Status(),
		//	Proto:     c.Request.Proto,
		//	Ip:        c.ClientIP(),
		//	//Body:      rbody,
		//	Query:   query,
		//	Message: c.Errors,
		//}
		//log.Print("*************************************************")
		//logger.LoggerWithFields(logBody, c)
		logMap := logrus.Fields{
			"SpendTime": spendTime,            //接口花费时间
			"path":      c.Request.RequestURI, //请求路径
			"Method":    c.Request.Method,     //请求方法
			"Status":    c.Writer.Status(),    //接口返回状态
			"Proto":     c.Request.Proto,      //http请求版本
			"Ip":        c.ClientIP(),         //IP地址
			//"Body":      body,      //请求体
			"Query":   query,    //请求query
			"Message": c.Errors, //返回错误信息
		}
		//log.Println(logMap, "日志打印")
		Log := logger.Logger.WithFields(logMap)
		if len(c.Errors) > 0 { // 矿建内部错误
			Log.Error(c.Errors.ByType(gin.ErrorTypePrivate))
		}
		if c.Writer.Status() > 200 {
			Log.Error()
		} else {
			Log.Info()
		}
	}
}

func MethodNotAllowedHandler(c *gin.Context) {
	fmt.Println("405不允许")
	c.Error(errors.New("405 Method Not Allowed"))
	return
	//c.AbortWithError(http.StatusMethodNotAllowed, errors.New("405 Method Not Allowed"))
	//return
}
func NotFoundHandler(c *gin.Context) {
	fmt.Println("404未找到")
	c.Error(errors.New("404 Not Found"))
	return
	//c.AbortWithError(http.StatusNotFound, errors.New("404 Not Found"))
	//return
}

/*
* @MethodName Cors
* @Description 跨域，限制请求方法，限制请求头
* @Author khr
* @Date 2023/7/29 9:52
 */

func Cors() gin.HandlerFunc {
	//log.Println("跨域中间件")
	return func(c *gin.Context) {
		//log.Println("进入中间件处理")
		//logger.Println("cors 进入中间件 - 处理请求之前")

		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Origin, X-CSRF-Token,X-Requested-With,Accept, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		// 允许放行OPTIONS请求
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
		//log.Println("中间件处理请求完毕")
	}
}

func DataEncrypr() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody, _ := c.GetRawData()
		privateKey, err := rsa.DecryptWithPrivateKey(string(requestBody))
		if err != nil {
			c.AbortWithStatusJSON(600, msg.EncryptError)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(privateKey))
		c.Set("reqbody", string(privateKey))
		c.Next()
	}
}
