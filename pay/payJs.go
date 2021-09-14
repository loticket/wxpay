package pay

import (
	"context"
	"encoding/json"
	"github.com/loticket/wxpay/auth/signers"
	"github.com/loticket/wxpay/utils"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var payOrderJsGateway = "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi"

type JsOrder struct {
	*Config
}


//实例化
//appID 微信唯一标识, mchID 商户好, private_key 私钥,SerialNo 证书编号, certificate 证书, notifyURL 通知地址
func NewJsOrderDefault(appID string, mchID string, private_key string,SerialNo string, certificate string, notifyURL string) *JsOrder {
	order := JsOrder{
		&Config{
			AppID:       appID,
			MchID:       mchID,
			SerialNo:    SerialNo,
			PrivateKey:  private_key,
			Certificate: certificate,
			NotifyURL:   notifyURL,
		},
	}
	return &order
}

// NewOrder return an instance of order package
func NewJsOrder(cfg *Config) *JsOrder {
	order := JsOrder{cfg}
	return &order
}

//支付
func (o *JsOrder) Pay(out_trade_no string,amount int64,openid string,description string,attach string)(BridgeJs,error) {
  var param Params = Params{
	  TotalFee:amount,
	  Description:description,
	  OutTradeNo:out_trade_no,
	  OpenID:openid,
	  Attach:attach,
  }

	return o.BridgeConfig(param)
}

//小程序支付下单
func (o *JsOrder) prePayOrder(p Params) (PayPrepay,error){
	var payJs JsapiPay = JsapiPay{
		Appid :o.AppID,
		Mchid :o.MchID,
		Description:p.Description,
		OutTradeNo:p.OutTradeNo,
		TimeExpire:time.Now().Local().Add(2*time.Hour).Format("2006-01-02T15:04:05Z07:00"),
		NotifyUrl:o.NotifyURL,
		Attach:p.Attach,
		GoodsTag:"",
		Amount: Amount{
			Total:p.TotalFee,
			Currency:"CNY",
		},
		Payer: Payer{
			Openid: p.OpenID,
		},
		SettleInfo: SettleInfo{
			ProfitSharing: true,
		},
	}

	client,err := o.NewClient()
	if err != nil {
		return PayPrepay{},err
	}

	httpResp,errRes:= client.Post(context.Background(), payOrderJsGateway,payJs)

	//fmt.Println(httpResp.Header.Values("Wechatpay-Serial"))

	if errRes != nil {
		return PayPrepay{},errRes
	}

	//获取返回参数
	resByte,errByte := ioutil.ReadAll(httpResp.Body)
	if errByte != nil {
		return PayPrepay{},errByte
	}

	defer httpResp.Body.Close()
	var apiPrepay PayPrepay
	if err = json.Unmarshal(resByte,&apiPrepay);err != nil {
		return PayPrepay{},err
	}


	return apiPrepay,nil
}

//app返回信息等待支付
func (o *JsOrder) BridgeConfig(p Params) (cfg BridgeJs, err error) {
  apiPrepay,err := o.prePayOrder(p)
  if err != nil {
  	 return BridgeJs{},err
  }

  var (
  	 buffer    strings.Builder
  	 timestamp = strconv.FormatInt(time.Now().Unix(), 10)
  	 nonceStr = utils.GenerateNonceStr(32)
  )
  apiPrepay.PrepayId = "prepay_id=" + apiPrepay.PrepayId
  buffer.WriteString(o.AppID)
  buffer.WriteString("\n")
  buffer.WriteString(timestamp)
  buffer.WriteString("\n")
  buffer.WriteString(nonceStr)//随机字符串
  buffer.WriteString("\n")
  buffer.WriteString(apiPrepay.PrepayId)
  buffer.WriteString("\n")

  var rsaStr string = buffer.String() //加密的串


	//解析私钥
	rasKey,err := o.PaserPrivateKey()
	if err != nil {
		return BridgeJs{},err
	}

	signature , err := signers.Sha256WithRsa(rsaStr,rasKey)
	if err != nil {
		return BridgeJs{}, err
	}
	
  return BridgeJs{
	  AppId:o.AppID,
	  TimeStamp:timestamp,
	  NonceStr:nonceStr,
	  SignType:"RSA",
          Package: apiPrepay.PrepayId,
	  PaySign:signature,
  },nil
}
