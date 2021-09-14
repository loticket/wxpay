package payorder

const (
	payOrderNativeGateway = "https://api.mch.weixin.qq.com/v3/pay/transactions/native"
	payOrderJsGateway     = "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi"
)

type Params struct {
	TotalFee    int64  //金额
	Description string //描述
	OutTradeNo  string //订单号码
	OpenID      string //用户openid
	Attach      string //附加数据
}

// 小程序调起支付参数
type BridgeJs struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	SignType  string `json:"signType"`
	Package   string `json:"package"`
	PaySign   string `json:"paySign"`
}

// app调起支付参数
type BridgeApp struct {
	AppId     string `json:"appId"`
	Partnerid string `json:"partnerid"`
	Prepayid  string `json:"prepayid"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	Sign      string `json:"sign"`
}

//订单金额
type Amount struct {
	Total    int64  `json:"total"`    //订单总金额，单位为分。 示例值：100
	Currency string `json:"currency"` //CNY：人民币，境内商户号仅支持人民币。
}

//支付者
type Payer struct {
	Openid string `json:"openid"` //用户在直连商户appid下的唯一标识。
}

//结算信息
type SettleInfo struct {
	ProfitSharing bool `json:"profit_sharing"`
}

//jsapi支付请求
type JsapiPay struct {
	Appid       string     `json:"appid"`                 //由微信生成的应用ID，全局唯一
	Mchid       string     `json:"mchid"`                 //直连商户号
	Description string     `json:"description"`           //商品描述
	OutTradeNo  string     `json:"out_trade_no"`          //商户订单号
	TimeExpire  string     `json:"time_expire,omitempty"` //交易结束时间
	Attach      string     `json:"attach"`                //附加数据
	NotifyUrl   string     `json:"notify_url"`            //通知地址
	GoodsTag    string     `json:"goods_tag,omitempty"`   //订单优惠标记
	Amount      Amount     `json:"amount"`                //订单金额
	Payer       Payer      `json:"payer"`                 //支付者
	SettleInfo  SettleInfo `json:"settle_info"`           //结算信息
}

//预支付下单返回
type PayPrepay struct {
	PrepayId string `json:"prepay_id"`
}

//app支付
type AppPay struct {
	Appid       string `json:"appid"`                 //由微信生成的应用ID，全局唯一
	Mchid       string `json:"mchid"`                 //直连商户号
	Description string `json:"description"`           //商品描述
	OutTradeNo  string `json:"out_trade_no"`          //商户订单号
	TimeExpire  string `json:"time_expire,omitempty"` //交易结束时间
	Attach      string `json:"attach"`                //附加数据
	NotifyUrl   string `json:"notify_url"`            //通知地址
	GoodsTag    string `json:"goods_tag,omitempty"`   //订单优惠标记
	Amount      Amount `json:"amount"`                //订单金额
}

//native支付
type NativePay struct {
	Appid       string `json:"appid"`                 //由微信生成的应用ID，全局唯一
	Mchid       string `json:"mchid"`                 //直连商户号
	Description string `json:"description"`           //商品描述
	OutTradeNo  string `json:"out_trade_no"`          //商户订单号
	TimeExpire  string `json:"time_expire,omitempty"` //交易结束时间
	Attach      string `json:"attach"`                //附加数据
	NotifyUrl   string `json:"notify_url"`            //通知地址
	GoodsTag    string `json:"goods_tag,omitempty"`   //订单优惠标记
	Amount      Amount `json:"amount"`                //订单金额
}

//h5支付
type SceneInfo struct {
	PayerClientIp string `json:"payer_client_ip"`
	H5Info        H5Info `json:"h5_info"`
}

type H5Info struct {
	Types string `json:type`
}

type H5Pay struct {
	Appid       string    `json:"appid"`                 //由微信生成的应用ID，全局唯一
	Mchid       string    `json:"mchid"`                 //直连商户号
	Description string    `json:"description"`           //商品描述
	OutTradeNo  string    `json:"out_trade_no"`          //商户订单号
	TimeExpire  string    `json:"time_expire,omitempty"` //交易结束时间
	Attach      string    `json:"attach"`                //附加数据
	NotifyUrl   string    `json:"notify_url"`            //通知地址
	GoodsTag    string    `json:"goods_tag,omitempty"`   //订单优惠标记
	Amount      Amount    `json:"amount"`                //订单金额
	SceneInfo   SceneInfo `json:"scene_info"`
}

//预支付下单返回
type BridgeNative struct {
	CodeUrl string `json:"code_url"`
}

//预支付下单返回
type BridgeH5 struct {
	H5Url string `json:"h5_url"`
}
