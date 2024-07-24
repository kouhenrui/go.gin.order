package pay

import (
	"github.com/go-pay/gopay/alipay"
)

var (
	client     *alipay.Client
	err        error
	AppId      = "9021000123610932" // 你的支付宝应用id
	PrivateKey = ""
)

type AliPay struct {
	client *alipay.Client
}
type AliPaytor interface {
}

func NewAliPay() *AliPay {
	client, err = alipay.NewClient(AppId, PrivateKey, false)
	client.SetAliPayPublicCertSN("publiccert")
	client.SetLocation(alipay.LocationShanghai)
	client.SetCharset(alipay.UTF8)                //设置字符编码，不设置默认 utf-8
	client.SetSignType(alipay.RSA2)               // 设置签名类型，不设置默认 RSA2
	client.SetReturnUrl("http://127.0.0.1:8080")  // 设置返回URL
	client.SetNotifyUrl("https://127.0.0.1:8080") // 设置异步通知URL
	//client.SetAppAuthToken()                    // 设置第三方应用授权
	return &AliPay{client}
}

func (a *AliPay) WebPageAlipay() {
	//a.client.

	//pay := alipay.TradeOrderPay{}
	//// 支付成功之后，支付宝将会重定向到该 URL
	//pay.ReturnURL = "http://localhost:8088/return"
	////支付标题
	//pay.Subject = "支付宝支付测试"
	////订单号，一个订单号只能支付一次
	//pay.OutTradeNo = time.Now().String()
	////销售产品码，与支付宝签约的产品码名称,目前仅支持FAST_INSTANT_TRADE_PAY
	//pay.ProductCode = "FAST_INSTANT_TRADE_PAY"
	////金额
	//pay.TotalAmount = "0.01"
}
