package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.gin.order/src/config/casbin"
	"go.gin.order/src/config/database"
	"go.gin.order/src/config/logger"
	"go.gin.order/src/config/messagequeue"
	"go.gin.order/src/config/redis"
	"go.gin.order/src/config/translate"
	"go.gin.order/src/internal/pojo"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	v   *viper.Viper
	err error
)

var (
	Port     string
	Tcp      string
	language string
	JWTKey   string

	Captcha       pojo.Captcha
	mysqlClient   pojo.MysqlConf
	postgreClient pojo.PostGreConf
	redisClient   pojo.RedisConf
	casbinClient  pojo.CabinConf //权限连接实例
	logConf       pojo.LogCof    //连接日志实例化参数
	mqConf        pojo.RabbitmqConf
	mongoClient   pojo.MongoConf
	etcdArry      = []string{}
	WhiteUrl      = []string{"/api/v1/auth/login", "/api/v1/auth/register", "/api/v1/public/.*", "/api/v1/approval/.*", "/api/v1/approval/.*", "/ws"}
)

const (
	AdminExp = time.Hour * 6
	UserExp  = time.Hour * 3
)

func InitConfig() {
	log.Println("实例化配置文件")
	v = viper.New()         // 构建 Viper 实例
	v.SetConfigType("yaml") //指定文件后缀
	if runtime.GOOS == "windows" {
		v.SetConfigName("conf.development") // 配置文件名称(无扩展名)
		v.AddConfigPath(".")                // Adjust this path as necessary
	} else {
		v.SetConfigName("conf.production") // 配置文件名称(无扩展名)
		v.AddConfigPath("root/data/go/data/")
	}
	//logger.Println(v.)
	env := os.Getenv("GIN_ENV")
	if env == "" {
		env = "development" // Default environment
	}
	log.Println(env, "本机环境")

	if err = v.ReadInConfig(); err != nil { // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	} // 查找并读取配置文件
	viperLoadConf()
	v.WatchConfig() //开启监听
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file updated.")
		viperLoadConf() // 加载配置的方法
	})

}

func viperLoadConf() {
	Port = v.GetString("server.port")         //读取http请求端口号
	Tcp = v.GetString("server.tcp")           //读取tcp请求端口号
	language = v.GetString("server.language") //读取语言包
	JWTKey = v.GetString("server.jwt")        //读取jwt密钥
	//JwtExpT = v.GetInt("server.jwtexp")
	logConfig := v.GetStringMap("logger") //日志路径及名称设置
	captchaexp := v.GetInt("captcha.expire")
	captchaprefix := v.GetString("captcha.prefix")
	Captcha.Expired = time.Duration(captchaexp) * time.Minute
	Captcha.Prefix = captchaprefix
	//captchar := v.GetStringMap("captcha") //验证码缓存
	mysql := v.GetStringMap("mysql") //读取MySQL配置
	postgresql := v.GetStringMap("postgresql")
	rediss := v.GetStringMap("redis") //读取redis配置
	mq := v.GetStringMap("rabbitmq")  //读取rabbitmq配置
	cn := v.GetStringMap("casbin")    //读取casbin配置
	//ck := v.GetStringMap("click")    //读取click house配置
	mg := v.GetStringMap("mongo")
	//map转struct
	mapstructure.Decode(mysql, &mysqlClient)
	mapstructure.Decode(postgresql, &postgreClient)
	mapstructure.Decode(rediss, &redisClient)
	mapstructure.Decode(mq, &mqConf)
	mapstructure.Decode(logConfig, &logConf)
	mapstructure.Decode(cn, &casbinClient)
	mapstructure.Decode(mg, &mongoClient)
	//mapstructure.Decode(ck, &ClickConfig)
	//log.Println("casbinClient.", casbinClient)
	//mapstructure.Decode(captchar, &Captcha)
	//log.Println("Captcha.", Captcha)
	//////t := Captcha.Expired
	//Captcha.Expired = Captcha.Expired * time.Minute
	//log.Println("Captcha.Expired", Captcha.Expired)
	etcdConnect := v.GetStringSlice("etcd")
	//kafka := v.GetStringSlice("kafka")
	//oracle := v.GetStringSlice("oracle")
	etcdArry = append(etcdArry, etcdConnect...)
	log.Println("全局配置文件信息读取无误,开始载入")
	//Dbinit()         //mysql初始化
	//Redisinit() //redis初始化

	translate.TranslateInit(language) //i18语言设置

	logger.LogInit(&logConf) //日志初始化
	//etcd.EtcdInit(etcdArry)
	database.NewMysqlClient(&mysqlClient) //mysql连接
	database.InitPostgre(&postgreClient)
	redis.Redisinit(&redisClient)    //redis连接
	casbin.CasbinInit(&casbinClient) //casbin连接

	//etcd.EtcdInit(etcdConnect)
	messagequeue.Mqinit(&mqConf)
	//repository.TableInit()
}
