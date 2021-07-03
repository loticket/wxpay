package downloader

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/loticket/wxpay/client"
	"github.com/loticket/wxpay/consts"
	"github.com/loticket/wxpay/utils"
	"io/ioutil"
	"net/http"
	"sync"
)

// CertificateDownloader 平台证书下载器，下载完成后可直接获取 x509.Certificate 对象或导出证书内容
type CertificateDownloader struct {
	certContents map[string]string   // 证书文本内容，用于导出
	client       *client.Client      // 微信支付 API v3 Go SDK HTTPClient
	mchAPIv3Key  string              // 商户APIv3密钥
	lock         sync.RWMutex
}

// DownloadCertificates 立即下载平台证书列表
func (d *CertificateDownloader) DownloadCertificates() error {
	resp, err := d.performDownloading()
	if err != nil {
		return err
	}

	rawCertContentMap := make(map[string]string)
	certificateMap := make(map[string]*x509.Certificate)
	for _, rawCertificate := range resp.Data {
		certContent, err := d.decryptCertificate(rawCertificate.EncryptCertificate)
		if err != nil {
			return err
		}

		certificate, err := utils.LoadCertificate(certContent)
		if err != nil {
			return fmt.Errorf("parse downlaoded certificate failed: %v, certcontent:%v", err, certContent)
		}

		serialNo := utils.GetCertificateSerialNumber(*certificate)

		rawCertContentMap[serialNo] = certContent
		certificateMap[serialNo] = certificate
	}

	if len(certificateMap) == 0 {
		return fmt.Errorf("no certificate downloaded")
	}
	return nil
}

func (d *CertificateDownloader) decryptCertificate(encryptCertificate *encryptCertificate) (string, error) {
	plaintext, err := utils.DecryptAES256GCM(d.mchAPIv3Key, *encryptCertificate.AssociatedData, *encryptCertificate.Nonce, *encryptCertificate.Ciphertext)
	if err != nil {
		return "", fmt.Errorf("decrypt downloaded certificate failed: %v", err)
	}

	return plaintext, nil
}

func (d *CertificateDownloader) performDownloading() (*downloadCertificatesResponse, error) {
	result, err := d.client.Get(context.Background(), consts.WechatPayAPIServer+"/v3/certificates")
	if err != nil {
		return nil, err
	}

	resp := new(downloadCertificatesResponse)
	if err = d.UnMarshalResponse(result, resp); err != nil {
		return nil, err
	}
	return resp, nil
}


// UnMarshalResponse 将回包组织成结构化数据
func (d *CertificateDownloader) UnMarshalResponse(httpResp *http.Response, resp interface{}) error {
	body, err := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	if err != nil {
		return err
	}

	httpResp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err = json.Unmarshal(body, resp)
	if err != nil {
		return err
	}
	return nil
}