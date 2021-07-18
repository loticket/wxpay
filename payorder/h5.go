package payorder

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"
)

const (
	payOrderH5Gateway = "https://api.mch.weixin.qq.com/v3/pay/transactions/h5"
)


type PayH5 struct {
	pay PayOrder
}

func (p *PayH5) ClientOption(pay PayOrder) {
	p.pay = pay
}

func (p *PayH5) Pay(out_trade_no string,amount int64,openid string,description string,attach string)(map[string]interface{},error) {

	var nativePay H5Pay = H5Pay{
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
		SceneInfo:SceneInfo{
		  PayerClientIp: "127.0.0.1",
		  H5Info: H5Info{
		  	Types: "Wap",
		  },
		},
	}

	client,err := p.pay.NewClient()
	if err != nil {
		return nil,err
	}

	httpResp,errRes:= client.Post(context.Background(), payOrderH5Gateway,nativePay)
	if errRes != nil {
		return nil,errRes
	}

	resByte,errByte := ioutil.ReadAll(httpResp.Body)
	if errByte != nil {
		return nil,errByte
	}

	defer httpResp.Body.Close()

	var bridgeH5 BridgeH5

	if err := json.Unmarshal(resByte,&bridgeH5);err != nil {
		return nil, err
	}

	return map[string]interface{}{"h5_url":bridgeH5.H5Url}, nil
}

func NewH5Pay() Wxpay {
	return new(PayH5)
}