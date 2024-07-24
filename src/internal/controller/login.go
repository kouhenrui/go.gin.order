package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.gin.order/pkg/msg"
	"go.gin.order/src/config/dto"
	"go.gin.order/src/internal/service"
)

type LoginController struct {
	loginService *service.LoginService
}

func NewLoginController() *LoginController {
	return &LoginController{loginService: service.NewLoginService()}
}
func (l *LoginController) Login(c *gin.Context) {
	var loginBody dto.LoginBody
	if err := c.ShouldBind(&loginBody); err != nil {
		c.Error(err)
		return
	}
	tokenExp, cook, err := l.loginService.Login(loginBody)
	if err != nil {
		c.Error(err)
		return
	}
	c.SetCookie("name", cook.Value, cook.MaxAge, cook.Path, cook.Domain, cook.Https, cook.Http)
	c.Set("res", tokenExp)
}
func (l *LoginController) Register(c *gin.Context) {
	var registerBody dto.RegisterBody
	if err := c.ShouldBind(&registerBody); err != nil {
		c.Error(err)
	}
	err := l.loginService.Register(registerBody)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("res", msg.SUCCESS)
}
func (l *LoginController) Info(c *gin.Context) {
	var user dto.User
	cookievalue := c.GetString("user")
	_ = json.Unmarshal([]byte(cookievalue), &user)
	info, err := l.loginService.Info(user)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("res", info)
}
func (l *LoginController) Logout(c *gin.Context) {
	c.SetCookie("user", "", -1, "/", "", false, true)
	c.Set("res", "success")
}

func (l *LoginController) MerchantList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	name := c.DefaultQuery("name", "")
	list, _ := l.loginService.MerchantList(page, limit, name)
	c.Set("res", list)
}
func (l *LoginController) MerchantProducts(c *gin.Context) {
	id := c.Query("id")
	product, err := l.loginService.MerchantProduct(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("res", product)
}
func (l *LoginController) AllAccounts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	list, err := l.loginService.AllAccounts(page, limit)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("res", list)
}
