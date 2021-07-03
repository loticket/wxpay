package errors

const  (
	TRADE_ERROR = "交易错误"
	SYSTEMERROR = "系统错误"
	SIGN_ERROR  = "签名错误"
	RULELIMIT   = "业务规则限制"
	PARAM_ERROR = "参数错误"
	OUT_TRADE_NO_USED = "商户订单号重复"
	ORDERNOTEXIST = "订单不存在"
	ORDER_CLOSED = "订单已关闭"
	OPENID_MISMATCH = "openid和appid不匹配"
	NOTENOUGH = "余额不足"
	NOAUTH = "商户无权限"
	MCH_NOT_EXISTS = "商户号不存在"
	INVALID_TRANSACTIONID = "订单号非法"
	INVALID_REQUEST = "无效请求"
	BANKERROR = "银行系统异常"
	APPID_MCHID_NOT_MATCH = "appid和mch_id不匹配"
	ACCOUNTERROR = "账号异常"
)