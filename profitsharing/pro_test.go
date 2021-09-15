package profitsharing

import (
	"fmt"
	"testing"
)

func TestProfitshare(t *testing.T) {
    var AppId string =  "-"
    var MchId string = "-"
	var MchSerialNumber string = "-"
	var Certificate string = "-"
	var PrivateKey string = "-"

	pro := NewProfitshareDefault(AppId,MchId,PrivateKey,MchSerialNumber,Certificate)

    var req CreateOrderRequest = CreateOrderRequest{
		Appid:AppId,
		OutOrderNo:"-",
		TransactionId:"-",
		UnfreezeUnsplit:true,
		Receivers:[]CreateOrderReceiver{
			CreateOrderReceiver{
				Account:"oih3251--",
				Amount:200,
				Description:"分账提成",
				Type:"PERSONAL_OPENID",
			},
		},
	}


    res,err := pro.CreateOrder(req)
    if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}
