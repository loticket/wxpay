package payorder

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

type PayJsapi struct {
	pay PayOrder
}

func (p *PayJsapi) ClientOption(pay PayOrder)  {
	p.pay = pay
}

func (p *PayJsapi) Pay(out_trade_no string,amount int64,openid string,description string,attach string)(map[string]interface{},error) {

	var payJs JsapiPay = JsapiPay{
		Appid :p.pay.AppID,
		Mchid :p.pay.MchID,
		Description:description,
		OutTradeNo:out_trade_no,
		TimeExpire:time.Now().Local().Add(2*time.Hour).Format("2006-01-02T15:04:05Z07:00"),
		NotifyUrl:p.pay.NotifyURL,
		Attach:attach,
		GoodsTag:"",
		Amount: Amount{
			Total:amount,
			Currency:"CNY",
		},
		Payer: Payer{
			Openid: openid,
		},
	}

	client,err := p.pay.NewClient()
	if err != nil {
		return nil,err
	}

	httpResp,errRes:= client.Post(context.Background(), payOrderJsGateway,payJs)
	if errRes != nil {
		return nil,errRes
	}

	resByte,errByte := ioutil.ReadAll(httpResp.Body)
	if errByte != nil {
		return nil,errByte
	}

	defer httpResp.Body.Close()

	var payPrepay PayPrepay

	if err := json.Unmarshal(resByte,&payPrepay);err != nil {
		return nil, err
	}

	return p.paserPrePay(payPrepay)
}

func (p *PayJsapi) paserPrePay(payPrepay PayPrepay)(map[string]interface{},error){

	var (
		buffer    strings.Builder
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
		nonceStr = utils.GenerateNonceStr(32)
	)
	var prepayId string = "prepay_id=" + payPrepay.PrepayId
	buffer.WriteString(p.pay.AppID)
	buffer.WriteString("\n")
	buffer.WriteString(timestamp)
	buffer.WriteString("\n")
	buffer.WriteString(nonceStr)//随机字符串
	buffer.WriteString("\n")
	buffer.WriteString(prepayId)
	buffer.WriteString("\n")

	var rsaStr string = buffer.String() //加密的串


	//解析私钥
	rasKey,err := p.pay.PaserPrivateKey()
	if err != nil {
		return nil,err
	}

	signature , err := signers.Sha256WithRsa(rsaStr,rasKey)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
     "appId":p.pay.AppID,
     "timeStamp":timestamp,
     "nonceStr":nonceStr,
     "package":prepayId,
     "signType":"RSA",
     "paySign":signature,
	},nil
}

func NewPayJsapiPay() Wxpay {
  return new(PayJsapi)
}