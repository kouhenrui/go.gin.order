package dto

type LoginBody struct {
	Email    string `json:"email,omitempty" binding:"required" `
	Password string `json:"password,omitempty" binding:"required"`
	//Captcha  Captcha `binding:"required"`
	Id      string `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
}
type RegisterBody struct {
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email" validate:"email"`
	Salt     string
	Phone    string
	Status   int
	Type     string
	Role     int32
	Captcha  Captcha
}
