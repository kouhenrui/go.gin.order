package service

import (
	"go.gin.order/pkg/captcha"
	"go.gin.order/src/config/dto"
)

type PublicService struct {
	cap   *captcha.Captcha
	capnr *captcha.Captchar
}

func NewPublicService() *PublicService {
	return &PublicService{
		cap:   captcha.NewCaptcha(),
		capnr: captcha.NewCaptchar(),
	}
}

type Publictor interface {
	MakeCaptchaNoRedis() (*dto.Captcha, error)
	MakeCaptcha() (*dto.Captcha, error)
}

func (p *PublicService) MakeCaptchaNoRedis() (*dto.Captcha, error) {
	err = p.capnr.CreateCaptchar()
	if err != nil {
		return nil, err
	}
	var ca = &dto.Captcha{}
	ca.Id = p.capnr.Id
	ca.Content = p.capnr.Content
	return ca, nil
}

func (p *PublicService) MakeCaptcha() (*dto.Captcha, error) {
	//captchator := p.cap.NewCaptchaInit()
	err = p.cap.CreateCaptcha()
	if err != nil {
		return nil, err
	}
	var ca = &dto.Captcha{}
	ca.Id = p.cap.Id
	ca.Content = p.cap.Content
	return ca, nil
}
func (p *PublicService) VerifyCaptcha(id string, content string) bool {
	p.cap.Id = id
	p.cap.Content = content
	return p.cap.VerifyCaptcha()
}
