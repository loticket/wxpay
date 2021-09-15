package profitsharing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/loticket/wxpay/base"
	"io/ioutil"
)

//分账接口
type Profitshare struct {
    base.BasePay
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



func NewProfitshareDefault(appID string, mchID string, private_key string,SerialNo string, certificate string) *Profitshare {
	return &Profitshare{
		base.BasePay{
			AppID:       appID,
			MchID:       mchID,
			SerialNo:    SerialNo,
			PrivateKey:  private_key,
			Certificate: certificate,
		},
	}
}