package profitsharing

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/loticket/wxpay/client"
	"github.com/loticket/wxpay/option"
	"github.com/loticket/wxpay/utils"
	"io/ioutil"
	"net/http"
	"time"
)

//分账接口
type Profitshare struct {
	AppID       string `json:"app_id"`      //用户的appid
	MchID       string `json:"mch_id"`      //商户号
	SerialNo    string `json:"serial_no"`   //商户证书序列号
	PrivateKey  string `json:"private_key"` //加密私钥
	Certificate string `json:"certificate"` //证书
	NotifyURL   string `json:"notify_url"`  //通知地址
}

//创建请求分账
func (p *Profitshare) CreateOrder(req CreateOrderRequest) (*OrdersEntity, error) {
	client, err := p.NewClient()
	if err != nil {
		return nil, err
	}

	httpResp, errRes := client.Post(context.Background(), SHAREORDER, req)
	if errRes != nil {
		return nil, errRes
	}

	//获取返回参数
	resByte, errByte := ioutil.ReadAll(httpResp.Body)

	if errByte != nil {
		return nil, errByte
	}

	defer httpResp.Body.Close()

	var ordersEntity OrdersEntity
	if err = json.Unmarshal(resByte, &ordersEntity); err != nil {
		return nil, err
	}

	return &ordersEntity, nil
}

//查询分账结果
func (p *Profitshare) QueryOrder(req QueryOrderRequest) (*OrdersEntity, error) {
	if req.OutOrderNo == "" {
		return nil, errors.New("OutOrderNo is empty")
	}

	var reqUrl = fmt.Sprintf(SHAREORDERSELECT,req.OutOrderNo,req.TransactionId)

	client, err := p.NewClient()
	if err != nil {
		return nil, err
	}

	httpResp, errRes := client.Get(context.Background(), reqUrl)
	if errRes != nil {
		return nil, errRes
	}

	//获取返回参数
	resByte, errByte := ioutil.ReadAll(httpResp.Body)

	if errByte != nil {
		return nil, errByte
	}

	defer httpResp.Body.Close()

	var ordersEntity OrdersEntity
	if err = json.Unmarshal(resByte, &ordersEntity); err != nil {
		return nil, err
	}

	return &ordersEntity, nil
}

// UnfreezeOrder 解冻剩余资金API
func (p *Profitshare) UnfreezeOrder(req UnfreezeOrderRequest) (resp *OrdersEntity, err error) {
	client, err := p.NewClient()
	if err != nil {
		return nil, err
	}

	httpResp, errRes := client.Post(context.Background(), SHAREORDERUNFREEZE, req)
	if errRes != nil {
		return nil, errRes
	}

	//获取返回参数
	resByte, errByte := ioutil.ReadAll(httpResp.Body)

	if errByte != nil {
		return nil, errByte
	}

	defer httpResp.Body.Close()

	var ordersEntity OrdersEntity
	if err = json.Unmarshal(resByte, &ordersEntity); err != nil {
		return nil, err
	}

	return &ordersEntity, nil
}

//解析字符串样式的私钥
func (c *Profitshare) PaserPrivateKey() (*rsa.PrivateKey, error) {
	var keyTemplate string = `-----BEGIN PRIVATE KEY-----
%s
-----END PRIVATE KEY-----`

	var key string = fmt.Sprintf(keyTemplate, c.PrivateKey)
	return utils.LoadPrivateKey(key)

}

func (c *Profitshare) PaserCertificate() (*x509.Certificate, error) {
	var keyTemplate string = `
-----BEGIN CERTIFICATE-----
%s
-----END CERTIFICATE-----`
	var key string = fmt.Sprintf(keyTemplate, c.Certificate)
	return utils.LoadCertificate(key)
}

//初始化请求客户端
func (c *Profitshare) NewClient() (*client.Client, error) {

	var (
		privateKey  *rsa.PrivateKey
		certificate *x509.Certificate
		err         error
	)

	//解析加密私钥
	if privateKey, err = c.PaserPrivateKey(); err != nil {
		return nil, err
	}
	//解析证书
	if certificate, err = c.PaserCertificate(); err != nil {
		return nil, err
	}

	opts := []option.ClientOption{
		option.WithMerchant(c.MchID, c.SerialNo, privateKey),   // 设置商户信息，用于生成签名信息
		option.WithWechatPay([]*x509.Certificate{certificate}), // 设置微信支付平台证书信息，对回包进行校验
		option.WithHTTPClient(&http.Client{}),                  // 可以不设置
		option.WithTimeout(2 * time.Second),                    // 自行进行超时时间配置
	}

	return client.NewClient(context.Background(), opts...)
}

func NewJsOrderDefault(appID string, mchID string, private_key string,SerialNo string, certificate string) *Profitshare {
	return &Profitshare{
			AppID:       appID,
			MchID:       mchID,
			SerialNo:    SerialNo,
			PrivateKey:  private_key,
			Certificate: certificate,
	}
}