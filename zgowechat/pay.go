package zgowechat

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/rubinus/zgo/zgoutils"
	"hash"
	"net/http"
	"strings"
	"sync"
)

/*
@Time : 2019-10-11 15:21
@Author : rubinus.chu
@File : pay
@project: zgo
*/

type Payer interface {
	//统一下单
	Order(body zgoutils.BodyMap) (wxRes *UnifiedOrderResponse, err error)

	//查询订单
	OrderQuery(body zgoutils.BodyMap) (wxRes *QueryOrderResponse, err error)

	//关闭订单
	OrderClose(body zgoutils.BodyMap) (wxRes *CloseOrderResponse, err error)

	//撤销订单
	OrderCancel(body zgoutils.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRes *ReverseResponse, err error)

	//申请退款
	OrderRefund(body zgoutils.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRes *RefundResponse, err error)

	//查询退款
	QueryRefund(body zgoutils.BodyMap) (wxRes *QueryRefundResponse, err error)

	//提交付款码支付
	MicroPay(body zgoutils.BodyMap) (wxRes *MicropayResponse, err error)

	//企业向微信用户个人付款
	Transfer(body zgoutils.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRes *TransfersResponse, err error)

	//下载对账单
	DownloadBill(body zgoutils.BodyMap) (wxRes string, err error)

	//下载资金账单
	DownloadFundFlow(body zgoutils.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRes string, err error)

	//拉取订单评价数据
	BatchQueryComment(body zgoutils.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRes string, err error)

	//设置支付国家
	SetCountry(country int)

	//设置微信服务器主动通知指定的页面http/https路径。
	SetNotifyUrl(url string)

	//获取微信支付所需参数里的Sign值
	GetParamSign(appId, mchId, apiKey string, bm zgoutils.BodyMap) (sign string)

	////获取微信支付沙箱环境所需参数里的Sign值
	GetSanBoxParamSign(appId, mchId, apiKey string, bm zgoutils.BodyMap) (sign string, err error)

	//解析微信支付异步通知的结果到zgoutils.BodyMap
	ParseNotifyResultToBodyMap(req *http.Request) (bm zgoutils.BodyMap, err error)

	//解析微信支付异步通知的参数
	ParseNotifyResult(req *http.Request) (notifyReq *NotifyRequest, err error)

	//解析微信退款异步通知的参数
	ParseRefundNotifyResult(req *http.Request) (notifyReq *RefundNotifyRequest, err error)

	//解密微信退款异步通知的加密数据
	DecryptRefundNotifyReqInfo(reqInfo, apiKey string) (refundNotify *RefundNotify, err error)

	//微信同步返回参数验签或异步通知参数验签
	VerifySign(apiKey, signType string, bean interface{}) (ok bool, err error)

	//APP支付，统一下单获取支付参数后，再次计算APP支付所需要的的sign
	GetAppPaySign(appid, partnerid, noncestr, prepayid, signType, timestamp, apiKey string) (paySign string)

	//微信内H5支付，统一下单获取支付参数后，再次计算出微信内H5支付需要用的paySign
	GetH5PaySign(appId, nonceStr, packages, signType, timeStamp, apiKey string) (paySign string)

	//JSAPI支付，统一下单获取支付参数后，再次计算出小程序用的paySign
	GetMiniPaySign(appId, nonceStr, prepayId, signType, timeStamp, apiKey string) (paySign string)

	//授权码查询openid(AccessToken:157字符)
	GetOpenIdByAuthCode(appId, mchId, apiKey, authCode, nonceStr string) (openIdRsp *OpenIdByAuthCodeRsp, err error)

	//添加微信证书 Byte 数组
	AddCertFileByte(certFile, keyFile, pkcs12File []byte)

	//添加微信证书 Path 路径
	AddCertFilePath(certFilePath, keyFilePath, pkcs12FilePath string) (err error)
}

type Country int

type PayClient struct {
	AppId      string
	MchId      string
	ApiKey     string
	BaseURL    string
	IsProd     bool
	NotifyUrl  string
	CertFile   []byte
	KeyFile    []byte
	Pkcs12File []byte
	mu         sync.RWMutex
}

//初始化微信客户端 ok
//    appId：应用ID
//    mchId：商户ID
//    ApiKey：API秘钥值
//    IsProd：是否是正式环境
func NewPayClient(appId, mchId, apiKey string, isProd bool) (client *PayClient) {
	return &PayClient{
		AppId:  appId,
		MchId:  mchId,
		ApiKey: apiKey,
		IsProd: isProd,
	}
}

//提交付款码支付 ok
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_10&index=1
func (w *PayClient) MicroPay(body zgoutils.BodyMap) (wxRes *MicropayResponse, err error) {
	var bs []byte
	if w.IsProd {
		bs, err = w.do(body, wxMicropay)
	} else {
		bs, err = w.do(body, wxSandboxMicropay)
	}
	if err != nil {
		return
	}
	wxRes = new(MicropayResponse)
	if err = xml.Unmarshal(bs, wxRes); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal：%v", err.Error())
	}
	return
}

//统一下单 ok
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func (w *PayClient) Order(body zgoutils.BodyMap) (wxRes *UnifiedOrderResponse, err error) {
	var bs []byte
	if w.IsProd {
		bs, err = w.do(body, wxUnifiedorder)
	} else {
		body.Set("total_fee", 101)
		bs, err = w.do(body, wxSandboxUnifiedorder)
	}
	if err != nil {
		return
	}
	wxRes = new(UnifiedOrderResponse)
	if err = xml.Unmarshal(bs, wxRes); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal：%v", err.Error())
	}
	return
}

//查询订单 ok
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
func (w *PayClient) OrderQuery(body zgoutils.BodyMap) (wxRes *QueryOrderResponse, err error) {
	var bs []byte
	if w.IsProd {
		bs, err = w.do(body, wxOrderquery)
	} else {
		bs, err = w.do(body, wxSandboxOrderquery)
	}
	if err != nil {
		return
	}
	wxRes = new(QueryOrderResponse)
	if err = xml.Unmarshal(bs, wxRes); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal：%v", err.Error())
	}
	return
}

//关闭订单 ok
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
func (w *PayClient) OrderClose(body zgoutils.BodyMap) (wxRes *CloseOrderResponse, err error) {
	var bs []byte
	if w.IsProd {
		bs, err = w.do(body, wxCloseorder)
	} else {
		bs, err = w.do(body, wxSandboxCloseorder)
	}
	if err != nil {
		return
	}
	wxRes = new(CloseOrderResponse)
	if err = xml.Unmarshal(bs, wxRes); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal：%v", err.Error())
	}
	return
}

//撤销订单 ok
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_11&index=3
func (w *PayClient) OrderCancel(body zgoutils.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRes *ReverseResponse, err error) {
	var (
		bs        []byte
		tlsConfig *tls.Config
	)
	if w.IsProd {
		if tlsConfig, err = w.addCertConfig(certFilePath, keyFilePath, pkcs12FilePath); err != nil {
			return nil, err
		}
		bs, err = w.do(body, wxReverse, tlsConfig)
	} else {
		bs, err = w.do(body, wxSandboxReverse)
	}
	if err != nil {
		return
	}
	wxRes = new(ReverseResponse)
	if err = xml.Unmarshal(bs, wxRes); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal：%v", err.Error())
	}
	return
}

//申请退款 ok
//    注意：如已使用client.AddCertFilePath()或client.AddCertFileByte()添加过证书，参数certFilePath、keyFilePath、pkcs12FilePath全传空字符串 ""，如方法需单独使用证书，则传证书Path
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
func (w *PayClient) OrderRefund(body zgoutils.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRes *RefundResponse, err error) {
	var (
		bs        []byte
		tlsConfig *tls.Config
	)
	if w.IsProd {
		if tlsConfig, err = w.addCertConfig(certFilePath, keyFilePath, pkcs12FilePath); err != nil {
			return nil, err
		}
		bs, err = w.do(body, wxRefund, tlsConfig)
	} else {
		bs, err = w.do(body, wxSandboxRefund)
	}
	if err != nil {
		return
	}
	wxRes = new(RefundResponse)
	if err = xml.Unmarshal(bs, wxRes); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal：%v", err.Error())
	}
	return
}

//查询退款 ok
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_5
func (w *PayClient) QueryRefund(body zgoutils.BodyMap) (wxRes *QueryRefundResponse, err error) {
	var bs []byte
	if w.IsProd {
		bs, err = w.do(body, wxRefundquery)
	} else {
		bs, err = w.do(body, wxSandboxRefundquery)
	}
	if err != nil {
		return
	}
	wxRes = new(QueryRefundResponse)
	if err = xml.Unmarshal(bs, wxRes); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal：%v", err.Error())
	}
	return
}

//下载对账单 ok
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_6
func (w *PayClient) DownloadBill(body zgoutils.BodyMap) (wxRes string, err error) {
	var bs []byte
	if w.IsProd {
		bs, err = w.do(body, wxDownloadbill)
	} else {
		bs, err = w.do(body, wxSandboxDownloadbill)
	}
	if err != nil {
		return
	}
	wxRes = string(bs)
	return
}

//下载资金账单 ok
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_18&index=7
//    好像不支持沙箱环境，因为沙箱环境默认需要用MD5签名，但是此接口仅支持HMAC-SHA256签名
func (w *PayClient) DownloadFundFlow(body zgoutils.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRes string, err error) {
	var (
		bs        []byte
		tlsConfig *tls.Config
	)
	if w.IsProd {
		if tlsConfig, err = w.addCertConfig(certFilePath, keyFilePath, pkcs12FilePath); err != nil {
			return null, err
		}
		bs, err = w.do(body, wxDownloadfundflow, tlsConfig)
	} else {
		bs, err = w.do(body, wxSandboxDownloadfundflow)
	}
	if err != nil {
		return
	}
	wxRes = string(bs)
	return
}

//拉取订单评价数据
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_17&index=11
//    好像不支持沙箱环境，因为沙箱环境默认需要用MD5签名，但是此接口仅支持HMAC-SHA256签名
func (w *PayClient) BatchQueryComment(body zgoutils.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRes string, err error) {
	var (
		bs        []byte
		tlsConfig *tls.Config
	)
	if w.IsProd {
		body.Set("sign_type", SignType_HMAC_SHA256)
		if tlsConfig, err = w.addCertConfig(certFilePath, keyFilePath, pkcs12FilePath); err != nil {
			return null, err
		}
		bs, err = w.do(body, wxBatchquerycomment, tlsConfig)
	} else {
		bs, err = w.do(body, wxSandboxBatchquerycomment)
	}
	if err != nil {
		return
	}
	wxRes = string(bs)
	return
}

//企业向微信用户个人付款
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_1
//    注意：此方法未支持沙箱环境，默认正式环境，转账请慎重
func (w *PayClient) Transfer(body zgoutils.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRes *TransfersResponse, err error) {
	body.Set("mch_appid", w.AppId)
	body.Set("mchid", w.MchId)
	var (
		bs        []byte
		tlsConfig *tls.Config
		agent     *gorequest.SuperAgent
		errs      []error
		res       gorequest.Response
	)
	if tlsConfig, err = w.addCertConfig(certFilePath, keyFilePath, pkcs12FilePath); err != nil {
		return nil, err
	}
	body.Set("sign", getReleaseSign(w.ApiKey, SignType_MD5, body))
	agent = zgoutils.HttpAgent().TLSClientConfig(tlsConfig)
	if w.BaseURL != null {
		agent.Post(w.BaseURL + wxTransfers)
	} else {
		agent.Post(wxBaseUrlCh + wxTransfers)
	}
	if res, bs, errs = agent.Type("xml").SendString(generateXml(body)).EndBytes(); len(errs) > 0 {
		return nil, errs[0]
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %v", res.StatusCode)
	}
	wxRes = new(TransfersResponse)
	if err = xml.Unmarshal(bs, wxRes); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal：%v", err.Error())
	}
	return
}

//向微信发送请求 ok
func (w *PayClient) do(body zgoutils.BodyMap, path string, tlsConfig ...*tls.Config) (bytes []byte, err error) {
	body.Set("appid", w.AppId)
	body.Set("mch_id", w.MchId)
	var (
		sign string
		errs []error
		res  gorequest.Response
	)
	if w.NotifyUrl != null {
		if path == wxUnifiedorder || path == wxRefund { //只在统一下单 和 退款时才用到notify_url
			body.Set("notify_url", w.NotifyUrl)
		}
	}
	if body.Get("sign") != null {
		goto GoRequest
	}
	if !w.IsProd {
		body.Set("sign_type", SignType_MD5)
		if sign, err = getSignBoxSign(w.MchId, w.ApiKey, body); err != nil {
			return
		}
	} else {
		sign = getReleaseSign(w.ApiKey, body.Get("sign_type"), body)
	}
	body.Set("sign", sign)
GoRequest:
	agent := zgoutils.HttpAgent()
	if w.IsProd && tlsConfig != nil {
		agent.TLSClientConfig(tlsConfig[0])
	}
	if w.BaseURL != null {
		agent.Post(w.BaseURL + path)
	} else {
		agent.Post(wxBaseUrlCh + path)
	}
	if res, bytes, errs = agent.Type("xml").SendString(generateXml(body)).EndBytes(); len(errs) > 0 {
		return nil, errs[0]
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %v", res.StatusCode)
	}
	if strings.Contains(string(bytes), "HTML") {
		return nil, errors.New(string(bytes))
	}
	return
}

//设置支付国家（默认：中国国内）
//    根据支付地区情况设置国家
//    country：<China：中国国内，China2：中国国内（冗灾方案），SoutheastAsia：东南亚，Other：其他国家>
func (w *PayClient) SetCountry(country int) {
	switch country {
	case 1:
		w.BaseURL = wxBaseUrlCh
	case 2:
		w.BaseURL = wxBaseUrlCh2
	case 3:
		w.BaseURL = wxBaseUrlHk
	case 4:
		w.BaseURL = wxBaseUrlUs
	default:
		w.BaseURL = wxBaseUrlCh
	}
}

//设置微信服务器主动通知指定的页面http/https路径。
func (w *PayClient) SetNotifyUrl(url string) {
	w.NotifyUrl = url
}

//获取微信支付正式环境Sign值
func getReleaseSign(apiKey string, signType string, bm zgoutils.BodyMap) (sign string) {
	var h hash.Hash
	if signType == SignType_HMAC_SHA256 {
		h = hmac.New(sha256.New, []byte(apiKey))
	} else {
		h = md5.New()
	}
	h.Write([]byte(bm.EncodeWechatSignParams(apiKey)))
	sign = strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	return
}

//获取微信支付沙箱环境Sign值
func getSignBoxSign(mchId, apiKey string, bm zgoutils.BodyMap) (sign string, err error) {
	var (
		sandBoxApiKey string
		h             hash.Hash
	)
	if sandBoxApiKey, err = getSanBoxKey(mchId, zgoutils.GetRandomString(32), apiKey, SignType_MD5); err != nil {
		return
	}
	h = md5.New()
	h.Write([]byte(bm.EncodeWechatSignParams(sandBoxApiKey)))
	sign = strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	return
}

//从微信提供的接口获取：SandboxSignKey
func getSanBoxKey(mchId, nonceStr, apiKey, signType string) (key string, err error) {
	body := make(zgoutils.BodyMap)
	body.Set("mch_id", mchId)
	body.Set("nonce_str", nonceStr)
	//沙箱环境：获取沙箱环境ApiKey
	if key, err = getSanBoxSignKey(mchId, nonceStr, getReleaseSign(apiKey, signType, body)); err != nil {
		return
	}
	return
}

//从微信提供的接口获取：SandboxSignkey
func getSanBoxSignKey(mchId, nonceStr, sign string) (key string, err error) {
	reqs := make(zgoutils.BodyMap)
	reqs.Set("mch_id", mchId)
	reqs.Set("nonce_str", nonceStr)
	reqs.Set("sign", sign)
	var (
		byteList    []byte
		errorList   []error
		keyResponse *getSignKeyResponse
	)
	if _, byteList, errorList = zgoutils.HttpAgent().Post(wxSandboxGetsignkey).Type("xml").SendString(generateXml(reqs)).EndBytes(); len(errorList) > 0 {
		return null, errorList[0]
	}
	keyResponse = new(getSignKeyResponse)
	if err = xml.Unmarshal(byteList, keyResponse); err != nil {
		return
	}
	if keyResponse.ReturnCode == "FAIL" {
		return null, errors.New(keyResponse.ReturnMsg)
	}
	return keyResponse.SandboxSignkey, nil
}

//生成请求XML的Body体
func generateXml(bm zgoutils.BodyMap) (reqXml string) {
	var buffer strings.Builder
	buffer.WriteString("<xml>")
	for key := range bm {
		buffer.WriteByte('<')
		buffer.WriteString(key)
		buffer.WriteString("><![CDATA[")
		buffer.WriteString(bm.Get(key))
		buffer.WriteString("]]></")
		buffer.WriteString(key)
		buffer.WriteByte('>')
	}
	buffer.WriteString("</xml>")
	reqXml = buffer.String()
	return
}
