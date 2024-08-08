package rest

import (
	"github.com/gin-gonic/gin"
	"go.gin.order/src/config/messagequeue"
	"go.gin.order/src/internal/controller"
	"go.gin.order/src/internal/middleware"
	"net/http"
)

func initApi(r *gin.Engine) {

	api := r.Group("/api/v1")
	{
		publiccontroller := controller.NewPublicController()
		api.GET("/public/captcha", publiccontroller.Captcha, middleware.LimitCaptcha())
		api.POST("/public/upload", publiccontroller.File)
		api.POST("/public/test", func(c *gin.Context) {
			body := c.PostForm("message")
			types := c.PostForm("type")
			//log.Println(body)
			mq, _ := messagequeue.NewRabbitMQ()
			var exchangname string
			switch types {
			case "direct":
				exchangname = "direct_exchange"
			case "broadcast":
				exchangname = "fanout_exchange"
			case "topic":
				exchangname = "topic_exchange"
			default:
				break
			}
			err := mq.Publish(exchangname, "", []byte(body))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.Set("res", body)
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
