package notice

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/loticket/wxpay/consts"
	"github.com/loticket/wxpay/utils"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Notice struct {
	AppID       string `json:"app_id"`      //用户的appid
	MchID       string `json:"mch_id"`      //商户号
	MchAPIv3Key string `json:"serial_no"`   //商户证书序列号
	PublicKey   *rsa.PublicKey  //平台的签名公钥
}

//解析通知内容
func (n *Notice) ParseNotify(request *http.Request,v3DecryptResult interface{}) error  {

	if err := n.checkParameters(request.Header);err != nil { //检查通知的header头部
		return err
	}


	body, err := n.getRequestBody(request)
	if err != nil  {
		return err
	}

	if len(body) == 0 {
		return errors.New("request body is empty")
	}

	//if err := n.VerifySign(request.Header,body);err != nil {
	//	return err
	//}


	var ret NotifyResponse
	if err := json.Unmarshal(body,&ret);err != nil {
		return err
	}

	plaintext, err := utils.DecryptAES256GCM(n.MchAPIv3Key, ret.Resource.AssociatedData, ret.Resource.Nonce, ret.Resource.Ciphertext)

	if err != nil {
		return err
	}

	ret.Resource.Plaintext = plaintext

	if err = json.Unmarshal([]byte(plaintext), &v3DecryptResult); err != nil {
		return err
	}

	return nil
}


func (n *Notice) checkParameters(header http.Header) error {

	if strings.TrimSpace(header.Get(consts.WechatPaySerial)) == "" {
		return fmt.Errorf("empty %s, WechatPaySerial=[%s]", consts.WechatPaySerial, requestID)
	}
	
	if strings.TrimSpace(header.Get(consts.WechatPayTimestamp)) == "" {
		return fmt.Errorf("empty %s,WechatPayTimestamp=[%s]", consts.WechatPayTimestamp, requestID)
	}

	if strings.TrimSpace(header.Get(consts.WechatPayNonce)) == "" {
		return fmt.Errorf("empty %s, WechatPayNonce=[%s]", consts.WechatPayNonce, requestID)
	}

	timeStampStr := strings.TrimSpace(header.Get(consts.WechatPayTimestamp))
	timeStamp, err := strconv.Atoi(timeStampStr)
	if err != nil {
		return fmt.Errorf("invalid timestamp:[%s] request-id=[%s] err:[%v]", timeStampStr, requestID, err)
	}

	if math.Abs(float64(timeStamp)-float64(time.Now().Unix())) >= consts.FiveMinute {
		return fmt.Errorf("timestamp=[%d] expires, request-id=[%s]", timeStamp, requestID)
	}
	return nil
}

func (n *Notice) buildMessage(header http.Header, body []byte) (string, error) {
	timeStamp := header.Get(consts.WechatPayTimestamp)
	nonce := header.Get(consts.WechatPayNonce)

	message := fmt.Sprintf("%s\n%s\n%s\n", timeStamp, nonce, string(body))
	return message, nil
}

func (n *Notice) getRequestBody(request *http.Request) ([]byte, error) {
	reqBody, err := request.GetBody()
	if err != nil {
		return nil, fmt.Errorf("get request body err: %v", err)
	}

	body, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return nil, fmt.Errorf("read request body err: %v", err)
	}

	return body, nil
}

//验证密钥是否正确
func (n *Notice) VerifySign(header http.Header,body []byte) (err error) {
	//检查验证签名信息

	message,err := n.buildMessage(header,body)
	if err != nil {
		return err
	}

	serialNumber := header.Get(consts.WechatPaySerial) //平台证书编号
	signature := header.Get(consts.WechatPaySignature) //通知签名

	if strings.TrimSpace(serialNumber) == "" {
		return fmt.Errorf("serialNumber is empty, verifier need input serialNumber")
	}
	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("message is empty, verifier need input message")
	}
	if strings.TrimSpace(signature) == "" {
		return fmt.Errorf("signature is empty, verifier need input signature")
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("verify failed: signature not base64 encoded")
	}

	hashed := sha256.Sum256([]byte(message))
	err = rsa.VerifyPKCS1v15(n.PublicKey, crypto.SHA256, hashed[:], sigBytes)
	if err != nil {
		return fmt.Errorf("verifty signature with public key err:%s", err.Error())
	}


	return nil
}
