package notice

type Notice struct {
	AppID       string `json:"app_id"`      //用户的appid
	MchID       string `json:"mch_id"`      //商户号
	SerialNo    string `json:"serial_no"`   //商户证书序列号
	PrivateKey  string `json:"private_key"` //加密私钥
}

//解析通知内容
func (n *Notice) ParseNotify()  {
	
}

//验证密钥是否正确
func (n *Notice) VerifySign(signInfo SignInfo) (err error) {

	return nil
}
