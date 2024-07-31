package msg

const (
	SUCCESS         = "成功"
	EmailNotFound   = "邮箱未找到"
	PwdError        = "密码错误"
	CaptchaNotFound = "验证码错误"
	MakeTokenError  = "token生成错误"
	FoundTokenError = "未找到请求头"

	EncryptError = "加密错误"

	AccountNotFoundError = "账号未找到"

	PasswordError = "密码错误"

	AccountTypeError     = "账号类型未找到"
	COOKIEERROR          = "cookie错误"
	COOKIEEXPRITIMEERROR = "cookie已失效"
	CAPTCHAERROR         = "验证码请求频繁，请稍后再试"

	CASBINFOUNDERROR = "权限未找到"

	APPROVALORDERERROR = "工单审核顺序错误，无法操作"
)
