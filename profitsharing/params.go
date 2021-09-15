package profitsharing

import "time"

const (
	SHAREORDER         = "https://api.mch.weixin.qq.com/v3/profitsharing/orders"                      //请求分账
	SHAREORDERSELECT   = "https://api.mch.weixin.qq.com/v3/profitsharing/orders/%s&transaction_id=%s" //查询分账结果
	SHAREORDERUNFREEZE = "https://api.mch.weixin.qq.com/v3/profitsharing/orders/unfreeze"             //解冻剩余资金API
)

type CreateOrderReceiver struct {
	// 1、类型是MERCHANT_ID时，是商户号 2、类型是PERSONAL_OPENID时，是个人openid  3、类型是PERSONAL_SUB_OPENID时，是个人sub_openid
	Account string `json:"account"`
	// 分账金额，单位为分，只能为整数，不能超过原订单支付金额及最大分账比例金额
	Amount int64 `json:"amount"`
	// 分账的原因描述，分账账单中需要体现
	Description string `json:"description"`
	// 可选项，在接收方类型为个人的时可选填，若有值，会检查与 name 是否实名匹配，不匹配会拒绝分账请求 1、分账接收方类型是PERSONAL_OPENID或PERSONAL_SUB_OPENID时，是个人姓名的密文（选传，传则校验） 此字段的加密的方式为：敏感信息加密说明 2、使用微信支付平台证书中的公钥 3、使用RSAES-OAEP算法进行加密 4、将请求中HTTP头部的Wechatpay-Serial设置为证书序列号
	Name string `json:"name,omitempty" encryption:"EM_APIV3"`
	// 1、MERCHANT_ID：商户号 2、PERSONAL_OPENID：个人openid（由父商户APPID转换得到） 3、PERSONAL_SUB_OPENID: 个人sub_openid（由子商户APPID转换得到）
	Type string `json:"type"`
}

type CreateOrderRequest struct {
	// 微信分配的服务商appid
	Appid string `json:"appid"`
	// 服务商系统内部的分账单号，在服务商系统内部唯一，同一分账单号多次请求等同一次。只能是数字、大小写字母_-|*@
	OutOrderNo string `json:"out_order_no"`
	// 分账接收方列表，可以设置出资商户作为分账接受方，最多可有50个分账接收方
	Receivers []CreateOrderReceiver `json:"receivers,omitempty"`
	// 微信分配的子商户公众账号ID，分账接收方类型包含PERSONAL_SUB_OPENID时必填。（直连商户不需要，服务商需要）
	SubAppid string `json:"sub_appid,omitempty"`
	// 微信支付分配的子商户号，即分账的出资商户号。（直连商户不需要，服务商需要）
	SubMchid string `json:"sub_mchid,omitempty"`
	// 微信支付订单号
	TransactionId string `json:"transaction_id"`
	// 1、如果为true，该笔订单剩余未分账的金额会解冻回分账方商户； 2、如果为false，该笔订单剩余未分账的金额不会解冻回分账方商户，可以对该笔订单再次进行分账。
	UnfreezeUnsplit bool `json:"unfreeze_unsplit"`
}

//分账返回信息
type OrdersEntity struct {
	// 微信分账单号，微信系统返回的唯一标识
	OrderId string `json:"order_id"`
	// 商户系统内部的分账单号，在商户系统内部唯一，同一分账单号多次请求等同一次。只能是数字、大小写字母_-|*@
	OutOrderNo string `json:"out_order_no"`
	// 分账接收方列表
	Receivers []OrderReceiverDetail `json:"receivers,omitempty"`
	// 分账单状态（每个接收方的分账结果请查看receivers中的result字段），枚举值： 1、PROCESSING：处理中 2、FINISHED：分账完成  * `PROCESSING` - 处理中，  * `FINISHED` - 分账完成，
	State string `json:"state"`
	// 微信支付分配的子商户号，即分账的出资商户号。（直连商户不需要，服务商需要）
	SubMchid string `json:"sub_mchid,omitempty"`
	// 微信支付订单号
	TransactionId string `json:"transaction_id"`
}

type OrderReceiverDetail struct {
	// 1、类型是MERCHANT_ID时，是商户号 2、类型是PERSONAL_OPENID时，是个人openid 3、类型是PERSONAL_SUB_OPENID时，是个人sub_openid
	Account string `json:"account"`
	// 分账金额，单位为分，只能为整数，不能超过原订单支付金额及最大分账比例金额
	Amount int64 `json:"amount"`
	// 分账创建时间，遵循RFC3339标准格式
	CreateTime time.Time `json:"create_time"`
	// 分账的原因描述，分账账单中需要体现
	Description string `json:"description"`
	// 微信分账明细单号，每笔分账业务执行的明细单号，可与资金账单对账使用
	DetailId string `json:"detail_id"`
	// 分账失败原因。包含以下枚举值： 1、ACCOUNT_ABNORMAL : 分账接收账户异常 2、NO_RELATION : 分账关系已解除 3、RECEIVER_HIGH_RISK : 高风险接收方 4、RECEIVER_REAL_NAME_NOT_VERIFIED : 接收方未实名 5、NO_AUTH : 分账权限已解除  * `ACCOUNT_ABNORMAL` - 分账接收账户异常，  * `NO_RELATION` - 分账关系已解除，  * `RECEIVER_HIGH_RISK` - 高风险接收方，  * `RECEIVER_REAL_NAME_NOT_VERIFIED` - 接收方未实名，  * `NO_AUTH` - 分账权限已解除，
	FailReason string `json:"fail_reason,omitempty"`
	// 分账完成时间，遵循RFC3339标准格式
	FinishTime time.Time `json:"finish_time"`
	// 枚举值： 1、PENDING：待分账 2、SUCCESS：分账成功 3、CLOSED：已关闭  * `PENDING` - 待分账，  * `SUCCESS` - 分账成功，  * `CLOSED` - 已关闭，
	Result string `json:"result"`
	// 1、MERCHANT_ID：商户号 2、PERSONAL_OPENID：个人openid（由父商户APPID转换得到） 3、PERSONAL_SUB_OPENID: 个人sub_openid（由子商户APPID转换得到）  * `MERCHANT_ID` - 商户号，  * `PERSONAL_OPENID` - 个人openid（由父商户APPID转换得到），  * `PERSONAL_SUB_OPENID` - 个人sub_openid（由子商户APPID转换得到）（直连商户不需要，服务商需要），
	Type string `json:"type"`
}

type QueryOrderRequest struct {
	// 微信支付分配的子商户号，即分账的出资商户号。（直连商户不需要，服务商需要）
	SubMchid string `json:"sub_mchid,omitempty"`
	// 微信支付订单号
	TransactionId string `json:"transaction_id"`
	// 商户系统内部的分账单号，在商户系统内部唯一，同一分账单号多次请求等同一次。只能是数字、大小写字母_-|*@ 。 微信分账单号与商户分账单号二选一填写
	OutOrderNo string `json:"out_order_no"`
}

type UnfreezeOrderRequest struct {
	// 分账的原因描述，分账账单中需要体现
	Description string `json:"description"`
	// 商户系统内部的分账单号，在商户系统内部唯一，同一分账单号多次请求等同一次。只能是数字、大小写字母_-|*@
	OutOrderNo string `json:"out_order_no"`
	// 微信支付分配的子商户号，即分账的出资商户号。（直连商户不需要，服务商需要）
	SubMchid string `json:"sub_mchid,omitempty"`
	// 微信支付订单号
	TransactionId string `json:"transaction_id"`
}
