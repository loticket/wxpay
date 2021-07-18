package wxpay

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/loticket/wxpay/option"
	"github.com/loticket/wxpay/payorder"
	"github.com/loticket/wxpay/utils"
	"github.com/loticket/wxpay/client"
	"net/http"
	"time"
)

const (
	PAYJSAPI    = "jsapi"
	PAYNAVITIVE = "native"
	PAYAPP      = "app"
	PAYH5       = "h5"
)


type PayOrderWx struct {
	AppID       string `json:"app_id"`      //用户的appid
	MchID       string `json:"mch_id"`      //商户号
	SerialNo    string `json:"serial_no"`      //商户证书序列号
	PrivateKey  string `json:"private_key"` //加密私钥
	Certificate string `json:"certificate"` //证书
	NotifyURL   string `json:"notify_url"`  //通知地址
}

//下单支付 获取需要支付的信息
func (p *PayOrderWx) PlaceOrder(paytype string,out_trade_no string,amount int64,openid string,description string,attach string)(map[string]interface{},error)  {
	 var config payorder.PayOrder = payorder.PayOrder{
		 AppID:p.AppID,
		 MchID:p.MchID,
		 SerialNo:p.SerialNo,
		 PrivateKey:p.PrivateKey,
		 Certificate:p.Certificate,
		 NotifyURL:p.NotifyURL,
	 }

	 var pay payorder.Wxpay = p.GetPayType(paytype)
	 pay.ClientOption(config)
	 return pay.Pay(out_trade_no,amount,openid,description,attach)
}


//根据支付方式获取需要支付的类型
func (p *PayOrderWx)GetPayType(paytype string) payorder.Wxpay {
	switch paytype {
	case PAYJSAPI:
		return payorder.NewPayJsapiPay()
	case PAYNAVITIVE:
		return payorder.NewNativePay()
	case PAYAPP:
		return payorder.NewAppPay()
	case PAYH5:
		return payorder.NewH5Pay()
	default:
		return nil
	}
	return nil
}
//解析字符串样式的私钥
func (c *PayOrderWx) PaserPrivateKey() (*rsa.PrivateKey,error) {
	var keyTemplate string =
		`-----BEGIN PRIVATE KEY-----
%s
-----END PRIVATE KEY-----`;

	var key string = fmt.Sprintf(keyTemplate,c.PrivateKey)
	return utils.LoadPrivateKey(key)

}

func  (c *PayOrderWx) PaserCertificate() (*x509.Certificate,error) {
	var keyTemplate string   = `
-----BEGIN CERTIFICATE-----
%s
-----END CERTIFICATE-----`
	var key string = fmt.Sprintf(keyTemplate,c.Certificate)
	return utils.LoadCertificate(key)
}

//初始化请求客户端
func (c *PayOrderWx) NewClient()(*client.Client,error) {

	var (
		privateKey  *rsa.PrivateKey
		certificate *x509.Certificate
		err         error
	)

	//解析加密私钥
	if privateKey ,err = c.PaserPrivateKey();err != nil {
		return nil, err
	}
	//解析证书
	if certificate,err = c.PaserCertificate();err != nil {
		return nil, err
	}



	opts := []option.ClientOption{
		option.WithMerchant(c.MchID, c.SerialNo, privateKey),    // 设置商户信息，用于生成签名信息
		option.WithWechatPay([]*x509.Certificate{certificate}),  // 设置微信支付平台证书信息，对回包进行校验
		option.WithHTTPClient(&http.Client{}),                   // 可以不设置
		option.WithTimeout(2 * time.Second),                     // 自行进行超时时间配置
	}

	return client.NewClient(context.Background(),opts...)
}