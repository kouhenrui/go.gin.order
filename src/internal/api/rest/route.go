package rest

import (
	"github.com/gin-gonic/gin"
	"go.gin.order/src/config"
	"go.gin.order/src/config/ws"
	"go.gin.order/src/internal/middleware"
	"log"
)

func InitHttp() {
	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.MaxMultipartMemory = 20 << 20

	r.Use(middleware.Cors())
	//r.Use(middleware.TokenMiddleware()) //token检测
	//r.Use(middleware.DataEncrypr())                //解密中间件，将请求体解密给日志存放了reqbody参数
	r.Use(middleware.VerifyCookie())               //cookie验证
	r.Use(middleware.LoggerMiddleWare())           //日志捕捉
	r.Use(middleware.GlobalErrorMiddleware())      //错误捕捉
	r.Use(middleware.UnifiedResponseMiddleware())  //全局统一返回格式，添加了rsa
	initApi(r)                                     //挂载请求路径
	r.NoRoute(middleware.NotFoundHandler)          //404
	r.NoMethod(middleware.MethodNotAllowedHandler) //405，方法为找到
	//fmt.Println("中间加载结束", config.Port)
	hub := ws.NewHub()
	go hub.Run()
	//ws.NewConsumerHub(hub)
	r.GET("/ws", func(c *gin.Context) {
		ws.WsInit(hub, c.Writer, c.Request)
	})
	//http.ListenAndServe(":9999", nil)
	log.Println("ws server on 9999")
	//err := r.RunTLS(config.Port, "../../config/https/certificate.crt", "../../config/https/private.key")
	err := r.Run(config.Port)
	//logger.Println(su)
	if err != nil {
		log.Print(err)
		panic("端口启动错误")
	}
}
