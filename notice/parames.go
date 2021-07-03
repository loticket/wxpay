package notice

//签名--主要验证签名的来源
type SignInfo struct {
	HeaderTimestamp string `json:"Wechatpay-Timestamp"`
	HeaderNonce     string `json:"Wechatpay-Nonce"`
	HeaderSignature string `json:"Wechatpay-Signature"`
	HeaderSerial    string `json:"Wechatpay-Serial"`
	SignBody        string `json:"sign_body"`
}

type NotifyResponse struct {
	Id            string   `json:"id"`     //通知的唯一ID
	EventType     string   `json:"event_type"` //通知的类型，支付成功通知的类型为TRANSACTION.SUCCESS 示例值：TRANSACTION.SUCCESS
	ResourceType  string   `json:"resource_type"` //通知的资源数据类型，支付成功通知为encrypt-resource 示例值：encrypt-resource
	CreateTime    string   `json:"create_time"`
	Resource     *NotifyResource `json:"resource"`
	Summary      string    `json:"summary"` //回调摘要 示例值：支付成功
	SignInfo     *SignInfo `json:"-"`
}
type NotifyResource struct {
	Algorithm       string `json:"algorithm"` //对开启结果数据进行加密的加密算法，目前只支持AEAD_AES_256_GCM
	Ciphertext      string `json:"ciphertext"` //Base64编码后的开启/停用结果数据密文
	AssociatedData  string `json:"associated_data"` //附加数据
	OriginalType    string `json:"original_type"` //原始回调类型，为transaction
	Nonce           string `json:"nonce"` //加密使用的随机串
}


//解密后
type V3DecryptResult struct {
	Appid           string             `json:"appid"`
	Mchid           string             `json:"mchid"`
	OutTradeNo      string             `json:"out_trade_no"`
	TransactionId   string             `json:"transaction_id"`
	TradeType       string             `json:"trade_type"`
	TradeState      string             `json:"trade_state"`
	TradeStateDesc  string             `json:"trade_state_desc"`
	BankType        string             `json:"bank_type"`
	Attach          string             `json:"attach"`
	SuccessTime     string             `json:"success_time"`
	Payer           *Payer             `json:"payer"`
	Amount          *Amount            `json:"amount"`
	SceneInfo       *SceneInfo         `json:"scene_info"`
	PromotionDetail []*PromotionDetail `json:"promotion_detail"`
}

type Payer struct {
	Openid string `json:"openid"` // 用户在直连商户appid下的唯一标识
}

type Amount struct {
	Total         int    `json:"total,omitempty"`          // 订单总金额，单位为分
	PayerTotal    int    `json:"payer_total,omitempty"`    // 用户支付金额，单位为分
	Currency      string `json:"currency,omitempty"`       // CNY：人民币，境内商户号仅支持人民币
	PayerCurrency string `json:"payer_currency,omitempty"` // 用户支付币种
}

type SceneInfo struct {
	DeviceId string `json:"device_id,omitempty"` // 商户端设备号（发起扣款请求的商户服务器设备号）
}

type PromotionDetail struct {
	Amount              int            `json:"amount"`                         // 优惠券面额
	CouponId            string         `json:"coupon_id"`                      // 券ID
	Name                string         `json:"name,omitempty"`                 // 优惠名称
	Scope               string         `json:"scope,omitempty"`                // 优惠范围：GLOBAL：全场代金券, SINGLE：单品优惠
	Type                string         `json:"type,omitempty"`                 // 优惠类型：CASH：充值, NOCASH：预充值
	StockId             string         `json:"stock_id,omitempty"`             // 活动ID
	WechatpayContribute int            `json:"wechatpay_contribute,omitempty"` // 微信出资，单位为分
	MerchantContribute  int            `json:"merchant_contribute,omitempty"`  // 商户出资，单位为分
	OtherContribute     int            `json:"other_contribute,omitempty"`     // 其他出资，单位为分
	Currency            string         `json:"currency,omitempty"`             // CNY：人民币，境内商户号仅支持人民币
	GoodsDetail         []*GoodsDetail `json:"goods_detail,omitempty"`         // 单品列表信息
}

type GoodsDetail struct {
	GoodsId         string `json:"goods_id"`                    // 商品编码
	Quantity        int    `json:"quantity"`                    // 用户购买的数量
	UnitPrice       int    `json:"unit_price"`                  // 商品单价，单位为分
	DiscountAmount  int    `json:"discount_amount"`             // 商品优惠金额
	GoodsRemark     string `json:"goods_remark,omitempty"`      // 商品备注信息
	MerchantGoodsID string `json:"merchant_goods_id,omitempty"` // 商户侧商品编码，服务商模式下无此字段
}
