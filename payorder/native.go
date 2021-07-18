package payorder

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"
)

type PayNaviate struct {
	pay PayOrder
}

func (p *PayNaviate) ClientOption(pay PayOrder) {
	p.pay = pay
}

func (p *PayNaviate) Pay(out_trade_no string,amount int64,openid string,description string,attach string)(map[string]interface{},error) {

	var nativePay NativePay = NativePay{
		Appid:       p.pay.AppID,
		Mchid:       p.pay.MchID,
		Description: description,
		OutTradeNo:  out_trade_no,
		TimeExpire:  time.Now().Local().Add(2 * time.Hour).Format("2006-01-02T15:04:05Z07:00"),
		Attach:attach,
		NotifyUrl:p.pay.NotifyURL,
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

	httpResp,errRes:= client.Post(context.Background(), payOrderNativeGateway,nativePay)
	if errRes != nil {
		return nil,errRes
	}

	resByte,errByte := ioutil.ReadAll(httpResp.Body)
	if errByte != nil {
		return nil,errByte
	}

	defer httpResp.Body.Close()

	var bridgeNative BridgeNative

	if err := json.Unmarshal(resByte,&bridgeNative);err != nil {
		return nil, err
	}

	return map[string]interface{}{"code_url":bridgeNative.CodeUrl}, nil
}

func NewNativePay() Wxpay {
	return new(PayNaviate)
}