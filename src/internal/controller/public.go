package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.gin.order/src/config/dto"
	"go.gin.order/src/internal/service"
	"log"
	"net/http"
	"path/filepath"
)

type PublicController struct {
	publiceService *service.PublicService
}

func NewPublicController() *PublicController {
	return &PublicController{publiceService: service.NewPublicService()}
}
func (p *PublicController) Captcha(c *gin.Context) {
	ca := &dto.Captcha{}
	ca, err := p.publiceService.MakeCaptcha()
	//ca, err = loginService.MakeCaptchaNoRedis()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("res", ca)
}

func (p *PublicController) File(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	// 获取文件名并将其保存到指定目录
	filename := filepath.Base(file.Filename)
	savePath := filepath.Join("uploads", filename)
	log.Println(savePath, "************")
	// 保存文件到服务器
	if err = c.SaveUploadedFile(file, savePath); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Set("res", fmt.Sprintf("'%s' uploaded successfully!", file.Filename))
}
