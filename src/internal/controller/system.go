package controller

import (
	"github.com/gin-gonic/gin"
	"go.gin.order/src/internal/service"
)

type SystemController struct {
	systemService *service.SystemService
}

func NewSystemController() *SystemController {
	return &SystemController{
		systemService: service.NewSystemService(),
	}
}

func (s *SystemController) RoleList(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("limit", "10")
	list, err := s.systemService.RoleList(page, limit)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("res", list)
}

func (s *SystemController) AuthorizedAccount(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("limit", "10")
	//s.systemService.AuthorizedAccount(page, limit)
	list, total, err := s.systemService.AuthorizedAccount(limit, page)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("res", gin.H{"total": total, "list": list})
}
func (s *SystemController) UnAuthorizedAccount(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("limit", "10")
	list, total, err := s.systemService.UnAuthorizedAccount(page, limit)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("res", gin.H{"total": total, "list": list})
}
