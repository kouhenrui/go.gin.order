package service

import (
	"encoding/json"
	"errors"
	"go.gin.order/pkg/encrypt"
	"go.gin.order/pkg/msg"
	"go.gin.order/pkg/token"
	"go.gin.order/src/config"
	"go.gin.order/src/config/dto"
	"go.gin.order/src/config/redis"
	"go.gin.order/src/internal/pojo"
	"go.gin.order/src/internal/repository"
	"log"
	"strconv"
	"time"
)

var (
	crypt     = encrypt.RSAEncrypt{}
	makeToken = token.Token{}
	err       error
)

type LoginService struct {
	accountrep  *repository.Account
	merchantep  *repository.Mertant
	redisClient *redis.NewRedis
	//crypt *encrypt.RSAEncrypt
	//pservice    *PublicService
}

func NewLoginService() *LoginService {
	return &LoginService{
		accountrep:  repository.NewAccount(),
		redisClient: redis.NewRedisClient(),
		merchantep:  repository.NewMerchat(),
		//pservice:    NewPublicService(),
	}
}

func (l *LoginService) Login(body dto.LoginBody) (*dto.TokenAndExp, *dto.Cookie, error) {
	//校验验证码
	//publiceservice := NewPublicService()
	//if !publiceservice.VerifyCaptcha(body.Id, body.Content) {
	//	return nil, nil, errors.New(msg.CaptchaNotFound)
	//}
	//校验账户
	account, err := l.validateaccount(body.Email, body.Password)
	if err != nil {
		return nil, nil, err
	}
	var expTime time.Duration
	var maxAge int
	switch account.Type {
	case "admin":
		expTime = config.UserExp
		maxAge = 6 * 3600
		break
	case "user":
		expTime = config.AdminExp
		maxAge = 12 * 3600
		break
	default:
		return nil, nil, errors.New(msg.AccountTypeError)
	}
	var cookvalue = &dto.User{
		Id:       account.Id,
		Name:     account.Name,
		Gender:   account.Gender,
		Role:     account.Role,
		RoleName: account.RoleName,
		Type:     account.Type,
	}
	data, _ := json.Marshal(cookvalue)
	var cook = &dto.Cookie{
		Name:   "user",
		Value:  string(data),
		MaxAge: maxAge,
		Path:   "/",
		Domain: "",
		Https:  false,
		Http:   true,
	}
	signToken := &dto.TokenAndExp{}
	if len(account.AccessToken) > 0 {
		e := l.redisClient.ExistRedis(account.AccessToken)
		if e {
			t := l.redisClient.GetRedis(account.AccessToken)
			json.Unmarshal([]byte(t), &signToken)
			return signToken, cook, nil
		}
	}
	signToken, err = l.maketoken(account, expTime)
	if err != nil {
		return nil, nil, err
	}
	return signToken, cook, err
}
func (l *LoginService) validateaccount(email, password string) (*repository.AccountRes, error) {
	account, err := l.accountrep.FindEmail(email)
	if err != nil {
		return nil, err
	}
	pwd, err := crypt.DePwdCode(account.Password, account.Salt)
	if err != nil {
		return nil, err
	}
	if pwd != password {
		return nil, errors.New(msg.PasswordError)
	}
	return account, nil
}
func (l *LoginService) maketoken(account *repository.AccountRes, exp_time time.Duration) (*dto.TokenAndExp, error) {
	tokenClaim := dto.TokenClaims{
		Id:    account.Id,
		Name:  account.Name,
		Phone: account.Phone,
		Type:  account.Type,
		Email: account.Email,
		Role:  account.Role,
	}
	signToken := &dto.TokenAndExp{}
	signToken, err = makeToken.SignToken(tokenClaim, exp_time)
	access_token := crypt.MakeSalt()
	err = l.accountrep.SaveAccessToken(account.Id, access_token)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(signToken)
	err := l.redisClient.SetRedis(access_token, data, exp_time)
	if err != nil {
		return nil, err
	}
	return signToken, nil
}
func (l *LoginService) Register(body dto.RegisterBody) error {
	_, err = l.accountrep.FindAccountByName(body.UserName)
	log.Println(err, "222")
	if err != nil {
		return err
	}
	salt := crypt.MakeSalt()
	password, _ := crypt.EnPwdCode(body.Password, salt)
	var acc = &pojo.Account{
		Name:     body.UserName,
		Password: password,
		Salt:     salt,
		Phone:    body.Phone,
		Status:   1,
		Email:    body.Email,
		Type:     body.Type,
		Role:     body.Role,
	}
	log.Println(acc.Name)
	return l.accountrep.AddAccount(*acc)
}
func (l *LoginService) MerchantList(p, li, name string) (*[]repository.MerchantListRes, error) {
	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(li)
	return l.merchantep.MerctantList(page, limit, name)
	//list, err := l.merchantep.MertantList(page, limit, name)
	//if err != nil {
	//	return nil, err
	//}
	//return list, nil
}
func (l *LoginService) Info(user dto.User) (*repository.InfoRes, error) {
	log.Println(user)
	return l.accountrep.InfoById(user.Id)
}

func (l *LoginService) MerchantProduct(id string) (*[]repository.ProductRes, error) {
	merchatid, _ := strconv.Atoi(id)
	return l.merchantep.MerchantProducts(uint(merchatid))
}
func (l *LoginService) AllAccounts(p, lm string) (*[]repository.AllAccountRes, error) {
	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(lm)
	return l.accountrep.AllAccount(page, limit)
}
