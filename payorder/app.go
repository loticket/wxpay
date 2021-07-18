package payorder

import (
	"context"
	"github.com/loticket/wxpay/auth/signers"
	"github.com/loticket/wxpay/utils"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	"encoding/json"
)

const payOrderAppGateway = "https://api.mch.weixin.qq.com/v3/pay/transactions/app"

type PayApp struct {
	pay PayOrder
}

func (p *PayApp) ClientOption(pay PayOrder)  {
	p.pay = pay
}

func (p *PayApp) Pay(out_trade_no string,amount int64,openid string,description string,attach string)(map[string]interface{},error) {
	var payApp AppPay = AppPay{
		Appid :p.pay.AppID,
		Mchid :p.pay.MchID,
		Description:description,
		OutTradeNo:out_trade_no,
		TimeExpire:time.Now().Local().Add(2*time.Hour).Format("2006-01-02T15:04:05Z07:00"),
		NotifyUrl:p.pay.NotifyURL,
		Attach:attach,
		GoodsTag:"wx",
		Amount: Amount{
			Total:amount,
			Currency:"CNY",
		},
	}

	client,err := p.pay.NewClient()
	if err != nil {
		return nil,err
	}

	httpResp,errRes:= client.Post(context.Background(), payOrderAppGateway,payApp)
	if errRes != nil {
		return nil,errRes
	}

	//获取返回参数
	resByte,errByte := ioutil.ReadAll(httpResp.Body)
	if errByte != nil {
		return nil,errByte
	}

	defer httpResp.Body.Close()
	var appPrepay PayPrepay
	if err = json.Unmarshal(resByte,&appPrepay);err != nil {
		return nil,err
	}


	return p.paserPrePay(appPrepay)

}


func (p *PayApp) paserPrePay(payPrepay PayPrepay)(map[string]interface{},error){

	var (
		buffer    strings.Builder
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
		nonceStr = utils.GenerateNonceStr(32)
	)

	buffer.WriteString(p.pay.AppID)
	buffer.WriteString("\n")
	buffer.WriteString(timestamp)
	buffer.WriteString("\n")
	buffer.WriteString(nonceStr)//随机字符串
	buffer.WriteString("\n")
	buffer.WriteString(payPrepay.PrepayId)
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
		"appid":p.pay.AppID,
		"partnerid":p.pay.MchID,
		"prepayid":payPrepay.PrepayId,
		"package":"Sign=WXPay",
		"timestamp":timestamp,
		"noncestr":nonceStr,
		"sign":signature,
	},nil
}

func NewAppPay() Wxpay {
	return new(PayApp)
}