package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pay/util"
	"go.gin.order/src/config/messagequeue"
	"go.gin.order/src/internal/controller"
	"go.gin.order/src/internal/middleware"
	"log"
)

func initApi(r *gin.Engine) {

	api := r.Group("/api/v1")
	{
		publiccontroller := controller.NewPublicController()
		api.GET("/public/captcha", publiccontroller.Captcha, middleware.LimitCaptcha())
		api.POST("/public/upload", publiccontroller.File)
		api.GET("/public/test", func(c *gin.Context) {
			message := util.RandomString(8)
			pro, err := messagequeue.NewProducer("approval")
			if err != nil {
				c.Error(err)
				return
			}
			log.Println(pro)
			err = pro.Publish(message)
			if err != nil {
				c.Error(err)
				return
			}
			//defer pro.Close()
			//defer pro.Close()
			//var l = make(chan string)
			//con, _ := messagequeue.NewConsumer("approval")
			//defer con.Close()
			//go con.Consume(func(message string) {
			//	l <- message
			//})
			//for i := range l {
			//	log.Println(i, "发送的信息，接收到是")
			//}
			c.Set("res", message)
		})
		auth := api.Group("/auth")
		logincontroller := controller.NewLoginController()
		auth.POST("/login", logincontroller.Login)
		auth.POST("/register", logincontroller.Register)
		auth.GET("/info", logincontroller.Info)
		auth.POST("/logout", logincontroller.Logout)
		auth.GET("/account/list", logincontroller.AllAccounts)

		auth.GET("/merchant/list", logincontroller.MerchantList)
		auth.GET("/merchant/id", logincontroller.MerchantProducts)
		system := api.Group("/system")
		systemController := controller.NewSystemController()
		system.GET("/role/list", systemController.RoleList)
		system.GET("/role/authUser/allocatedList", systemController.AuthorizedAccount)
		system.GET("/role/authUser/unallocatedList", systemController.UnAuthorizedAccount)

		approval := api.Group("/approval")
		approvalcontroller := controller.NewApprovalController()
		approval.POST("/create", approvalcontroller.CreateApproval)
		//	auth.POST("/register", login.Register)
		//	auth.POST("/login", login.Login)
		//	auth.GET("/info", login.Info)
	}
	//ws:=r.Group("/v1/api/ws")
	//static:=r.Group("/v1/api/static")

}
