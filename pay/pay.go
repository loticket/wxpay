package pay

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/loticket/wxpay/client"
	"github.com/loticket/wxpay/option"
	"github.com/loticket/wxpay/utils"
	"net/http"
	"time"
)

// Config config for pay
type Config struct {
	AppID       string `json:"app_id"`      //用户的appid
	MchID       string `json:"mch_id"`      //商户号
	SerialNo    string `json:"serial_no"`      //商户证书序列号
	PrivateKey  string `json:"private_key"` //加密私钥
	Certificate string `json:"certificate"` //证书
	NotifyURL   string `json:"notify_url"`  //通知地址

}

//解析字符串样式的私钥
func (c *Config) PaserPrivateKey() (*rsa.PrivateKey,error) {
	var keyTemplate string =
`-----BEGIN PRIVATE KEY-----
%s
-----END PRIVATE KEY-----`;

  var key string = fmt.Sprintf(keyTemplate,c.PrivateKey)
  return utils.LoadPrivateKey(key)

}

func  (c *Config) PaserCertificate() (*x509.Certificate,error) {
	var keyTemplate string   = `
-----BEGIN CERTIFICATE-----
%s
-----END CERTIFICATE-----`
	var key string = fmt.Sprintf(keyTemplate,c.Certificate)
    return utils.LoadCertificate(key)
}

//初始化请求客户端
func (c *Config) NewClient()(*client.Client,error) {

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
