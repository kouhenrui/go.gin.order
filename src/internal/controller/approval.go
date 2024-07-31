package controller

import (
	"github.com/gin-gonic/gin"
	"go.gin.order/pkg/msg"
	"go.gin.order/pkg/util"
	"go.gin.order/src/config/dto"
	"go.gin.order/src/internal/service"
	"net/http"
)

type ApprovalController struct {
	approvalService *service.ApprovalService
}

func NewApprovalController() *ApprovalController {
	return &ApprovalController{approvalService: service.NewApprovalService()}
}

func (a *ApprovalController) CreateApproval(c *gin.Context) {
	var newapproval dto.CreateApprovalReq
	if err := c.ShouldBindJSON(&newapproval); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.GetValidate(err, &newapproval))
		return
	}
	err := a.approvalService.CreateApprove(&newapproval)
	//log.Println(err, "err")
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("res", msg.SUCCESS)
}
