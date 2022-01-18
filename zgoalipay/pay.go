package zgoalipay

import (
  "crypto"
  "crypto/rand"
  "crypto/rsa"
  "crypto/sha1"
  "crypto/sha256"
  "crypto/x509"
  "encoding/base64"
  "encoding/pem"
  "errors"
  "fmt"
  "github.com/gitcpu-io/zgo/zgoutils"
  "github.com/parnurzeal/gorequest"
  "hash"
  "net/http"
  "net/url"
  "strings"
  "time"
)

/*
@Time : 2019-10-11 15:21
@Author : rubinus.chu
@File : pay
@project: zgo
*/

type Payer interface {
  //统一收单交易支付接口
  OrderPay(body zgoutils.BodyMap) (tradeRes *TradePayResponse, err error)

  //统一收单线下交易查询
  OrderQuery(body zgoutils.BodyMap) (tradeRes *TradeQueryResponse, err error)

  //统一收单交易结算接口
  OrderSettle(body zgoutils.BodyMap) (tradeRes *TradeOrderSettleResponse, err error)

  //统一收单线下交易预创建
  Order(body zgoutils.BodyMap) (tradeRes *TradePrecreateResponse, err error)

  //统一收单交易创建接口
  OrderCreate(body zgoutils.BodyMap) (tradeRes *TradeCreateResponse, err error)

  //统一收单交易关闭接口
  OrderClose(body zgoutils.BodyMap) (tradeRes *TradeCloseResponse, err error)

  //统一收单交易撤销接口
  OrderCancel(body zgoutils.BodyMap) (tradeRes *TradeCancelResponse, err error)

  //统一收单交易退款接口
  OrderRefund(body zgoutils.BodyMap) (tradeRes *TradeRefundResponse, err error)

  //统一收单退款页面接口
  OrderPageRefund(body zgoutils.BodyMap) (tradeRes *TradePageRefundResponse, err error)

  //统一收单交易退款查询
  OrderFastPayRefundQuery(body zgoutils.BodyMap) (tradeRes *TradeFastpayRefundQueryResponse, err error)

  //统一收单下单并支付页面接口 -- pc出现二维码，手机支付宝扫一扫
  OrderPagePay(body zgoutils.BodyMap) (payUrl string, err error)

  //app支付接口2.0
  OrderAppPay(body zgoutils.BodyMap) (payUrl string, err error)

  //手机网站支付接口2.0 -- h5页面
  OrderWapPay(body zgoutils.BodyMap) (payUrl string, err error)

  //统一转账到支付宝账户接口
  FundTransUniTransfer(body zgoutils.BodyMap) (tradeRes *FundTransUniTransferResponse, err error)

  //统一转账查询接口
  FundTransCommonQuery(body zgoutils.BodyMap) (tradeRes *FundTransCommonQueryResponse, err error)

  //单笔转账到支付宝账户接口
  FundTransToaccountTransfer(body zgoutils.BodyMap) (tradeRes *FundTransToaccountTransferResponse, err error)

  //查询转账订单接口
  FundTransOrderQuery(body zgoutils.BodyMap) (tradeRes *FundTransOrderQueryResponse, err error)

  //支付宝资金账户资产查询接口
  FundAccountQuery(body zgoutils.BodyMap) (tradeRes *FundAccountQueryResponse, err error)

  //资金退回接口
  FundTransRefund(body zgoutils.BodyMap) (tradeRes *FundTansRefundResponse, err error)

  //支付宝订单信息同步接口
  //waiting

  //换取授权访问令牌
  SystemOauthToken(body zgoutils.BodyMap) (tradeRes *SystemOauthTokenResponse, err error)

  //换取应用授权令牌
  OpenAuthTokenApp(body zgoutils.BodyMap) (tradeRes *OpenAuthTokenAppResponse, err error)

  //支付宝会员授权信息查询接口
  UserInfoShare() (tradeRes *UserInfoShareResponse, err error)

  //芝麻分
  ZhimaCreditScoreGet(body zgoutils.BodyMap) (tradeRes *ZhimaCreditScoreGetResponse, err error)

  //************
  //设置 应用公钥证书SN
  SetAppCertSN(appCertSN string)

  //设置 支付宝根证书SN
  SetRootCertSN(rootCertSN string)

  //设置支付后的ReturnUrl
  SetReturnUrl(url string)

  //设置支付宝服务器主动通知商户服务器里指定的页面http/https路径。
  SetNotifyUrl(url string)

  //设置编码格式，如utf-8,gbk,gb2312等，默认推荐使用 utf-8
  SetCharset(charset string)

  //设置签名算法类型，目前支持RSA2和RSA，默认推荐使用 RSA2
  SetSignType(signType string)

  //设置应用授权
  SetAppAuthToken(appAuthToken string)

  //设置用户信息授权
  SetAuthToken(authToken string)

  //************************

  //解析支付宝支付完成后的Notify信息
  ParseNotifyResult(req *http.Request) (notifyReq *NotifyRequest, err error)

  //支付宝同步返回验签或异步通知验签
  VerifySign(aliPayPublicKey string, bean interface{}, syncSign ...string) (ok bool, err error)
}

type PayClient struct {
  AppId            string
  PrivateKey       string
  AlipayRootCertSN string
  AppCertSN        string
  ReturnUrl        string
  NotifyUrl        string
  Charset          string
  SignType         string
  AppAuthToken     string
  AuthToken        string
  IsProd           bool
}

//初始化支付宝客户端
//    注意：如果使用支付宝公钥证书验签，请设置 支付宝根证书SN（client.SetAlipayRootCertSN()）、应用公钥证书SN（client.SetAppCertSN()）
//    appId：应用ID
//    PrivateKey：应用私钥
//    IsProd：是否是正式环境
func NewPayClient(appId, privateKey string, isProd bool) (client *PayClient) {
  return &PayClient{
    AppId:      appId,
    PrivateKey: privateKey,
    IsProd:     isProd,
  }
}

//alipay.trade.fastpay.refund.query(统一收单交易退款查询)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.fastpay.refund.query
func (a *PayClient) OrderFastPayRefundQuery(body zgoutils.BodyMap) (tradeRes *TradeFastpayRefundQueryResponse, err error) {
  var (
    p1, p2 string
    bs     []byte
  )
  p1 = body.Get("out_trade_no")
  p2 = body.Get("trade_no")
  if p1 == null && p2 == null {
    return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "alipay.trade.fastpay.refund.query"); err != nil {
    return
  }
  tradeRes = new(TradeFastpayRefundQueryResponse)

  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayTradeFastpayRefundQueryResponse.Code != "10000" {
    info := tradeRes.AlipayTradeFastpayRefundQueryResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.trade.order.settle(统一收单交易结算接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.order.settle
func (a *PayClient) OrderSettle(body zgoutils.BodyMap) (tradeRes *TradeOrderSettleResponse, err error) {
  var (
    p1, p2 string
    bs     []byte
  )
  p1 = body.Get("out_request_no")
  p2 = body.Get("trade_no")
  if p1 == null || p2 == null {
    return nil, errors.New("out_request_no or trade_no are not allowed to be null")
  }
  if bs, err = a.do(body, "alipay.trade.order.settle"); err != nil {
    return
  }
  tradeRes = new(TradeOrderSettleResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayTradeOrderSettleResponse.Code != "10000" {
    info := tradeRes.AlipayTradeOrderSettleResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.trade.create(统一收单交易创建接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.create
func (a *PayClient) OrderCreate(body zgoutils.BodyMap) (tradeRes *TradeCreateResponse, err error) {
  var (
    p1, p2 string
    bs     []byte
  )
  p1 = body.Get("out_trade_no")
  p2 = body.Get("buyer_id")
  if p1 == null && p2 == null {
    return nil, errors.New("out_trade_no and buyer_id are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "alipay.trade.create"); err != nil {
    return
  }
  tradeRes = new(TradeCreateResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayTradeCreateResponse.Code != "10000" {
    info := tradeRes.AlipayTradeCreateResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.trade.close(统一收单交易关闭接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.close
func (a *PayClient) OrderClose(body zgoutils.BodyMap) (tradeRes *TradeCloseResponse, err error) {
  var (
    p1, p2 string
    bs     []byte
  )
  p1 = body.Get("out_trade_no")
  p2 = body.Get("trade_no")
  if p1 == null && p2 == null {
    return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "alipay.trade.close"); err != nil {
    return
  }
  tradeRes = new(TradeCloseResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayTradeCloseResponse.Code != "10000" {
    info := tradeRes.AlipayTradeCloseResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.trade.cancel(统一收单交易撤销接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.cancel
func (a *PayClient) OrderCancel(body zgoutils.BodyMap) (tradeRes *TradeCancelResponse, err error) {
  var (
    p1, p2 string
    bs     []byte
  )
  p1 = body.Get("out_trade_no")
  p2 = body.Get("trade_no")
  if p1 == null && p2 == null {
    return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "alipay.trade.cancel"); err != nil {
    return
  }
  tradeRes = new(TradeCancelResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayTradeCancelResponse.Code != "10000" {
    info := tradeRes.AlipayTradeCancelResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.trade.refund(统一收单交易退款接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.refund
func (a *PayClient) OrderRefund(body zgoutils.BodyMap) (tradeRes *TradeRefundResponse, err error) {
  var (
    p1, p2 string
    bs     []byte
  )
  p1 = body.Get("out_trade_no")
  p2 = body.Get("trade_no")
  if p1 == null && p2 == null {
    return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "alipay.trade.refund"); err != nil {
    return nil, err
  }
  tradeRes = new(TradeRefundResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayTradeRefundResponse.Code != "10000" {
    info := tradeRes.AlipayTradeRefundResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.trade.refund(统一收单退款页面接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.page.refund
func (a *PayClient) OrderPageRefund(body zgoutils.BodyMap) (tradeRes *TradePageRefundResponse, err error) {
  var (
    p1, p2 string
    bs     []byte
  )
  p1 = body.Get("out_trade_no")
  p2 = body.Get("trade_no")
  if p1 == null && p2 == null {
    return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "	alipay.trade.page.refund"); err != nil {
    return
  }
  tradeRes = new(TradePageRefundResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayTradePageRefundResponse.Code != "10000" {
    info := tradeRes.AlipayTradePageRefundResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.trade.precreate(统一收单线下交易预创建)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.precreate
func (a *PayClient) Order(body zgoutils.BodyMap) (tradeRes *TradePrecreateResponse, err error) {
  var bs []byte
  p := body.Get("out_trade_no")
  if p == null {
    return nil, errors.New("out_trade_no is not allowed to be null")
  }
  if bs, err = a.do(body, "alipay.trade.precreate"); err != nil {
    return
  }
  tradeRes = new(TradePrecreateResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayTradePrecreateResponse.Code != "10000" {
    info := tradeRes.AlipayTradePrecreateResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.trade.pay(统一收单交易支付接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.pay
func (a *PayClient) OrderPay(body zgoutils.BodyMap) (tradeRes *TradePayResponse, err error) {
  var bs []byte
  p := body.Get("out_trade_no")
  if p == null {
    return nil, errors.New("out_trade_no is not allowed to be null")
  }
  if bs, err = a.do(body, "alipay.trade.pay"); err != nil {
    return
  }
  tradeRes = new(TradePayResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayTradePayResponse.Code != "10000" {
    info := tradeRes.AlipayTradePayResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.trade.query(统一收单线下交易查询)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.query
func (a *PayClient) OrderQuery(body zgoutils.BodyMap) (tradeRes *TradeQueryResponse, err error) {
  var (
    p1, p2 string
    bs     []byte
  )
  p1 = body.Get("out_trade_no")
  p2 = body.Get("trade_no")
  if p1 == null && p2 == null {
    return nil, errors.New("out_trade_no and trade_no are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "alipay.trade.query"); err != nil {
    return
  }
  tradeRes = new(TradeQueryResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayTradeQueryResponse.Code != "10000" {
    info := tradeRes.AlipayTradeQueryResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.trade.app.pay(app支付接口2.0)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.app.pay
func (a *PayClient) OrderAppPay(body zgoutils.BodyMap) (payUrl string, err error) {
  var bs []byte
  trade := body.Get("out_trade_no")
  if trade == null {
    return null, errors.New("out_trade_no is not allowed to be null")
  }
  if bs, err = a.do(body, "alipay.trade.app.pay"); err != nil {
    return null, err
  }
  payUrl = string(bs)
  return
}

//alipay.trade.wap.pay(手机网站支付接口2.0)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.wap.pay
func (a *PayClient) OrderWapPay(body zgoutils.BodyMap) (payUrl string, err error) {
  var bs []byte
  trade := body.Get("out_trade_no")
  if trade == null {
    return null, errors.New("out_trade_no is not allowed to be null")
  }
  body.Set("product_code", "QUICK_WAP_WAY")
  if bs, err = a.do(body, "alipay.trade.wap.pay"); err != nil {
    return null, err
  }
  payUrl = string(bs)
  return
}

//alipay.trade.page.pay(统一收单下单并支付页面接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.page.pay
func (a *PayClient) OrderPagePay(body zgoutils.BodyMap) (payUrl string, err error) {
  var bs []byte
  trade := body.Get("out_trade_no")
  if trade == null {
    return null, errors.New("out_trade_no is not allowed to be null")
  }
  body.Set("product_code", "FAST_INSTANT_TRADE_PAY")
  if bs, err = a.do(body, "alipay.trade.page.pay"); err != nil {
    return null, err
  }
  payUrl = string(bs)
  return
}

//alipay.trade.orderinfo.sync(支付宝订单信息同步接口)
//    文档地址：https://docs.open.alipay.com/api_1/alipay.trade.orderinfo.sync
func (a *PayClient) OrderOrderinfoSync(body zgoutils.BodyMap) {

}

//alipay.system.oauth.token(换取授权访问令牌)
//    文档地址：https://docs.open.alipay.com/api_9/alipay.system.oauth.token
func (a *PayClient) SystemOauthToken(body zgoutils.BodyMap) (tradeRes *SystemOauthTokenResponse, err error) {
  var bs []byte
  grantType := body.Get("grant_type")
  if grantType == null {
    return nil, errors.New("grant_type is not allowed to be null")
  }
  code := body.Get("code")
  refreshToken := body.Get("refresh_token")
  if code == null && refreshToken == null {
    return nil, errors.New("code and refresh_token are not allowed to be null at the same time")
  }
  if bs, err = systemOauthToken(a.AppId, a.PrivateKey, body, "alipay.system.oauth.token", a.IsProd); err != nil {
    return
  }
  tradeRes = new(SystemOauthTokenResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipaySystemOauthTokenResponse.AccessToken == null {
    info := tradeRes.ErrorResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//向支付宝发送请求
func systemOauthToken(appId, privateKey string, body zgoutils.BodyMap, method string, isProd bool) (bytes []byte, err error) {
  body.Set("app_id", appId)
  body.Set("method", method)
  body.Set("format", "JSON")
  body.Set("charset", "utf-8")
  body.Set("sign_type", "RSA2")
  body.Set("timestamp", time.Now().Format(TimeLayout))
  body.Set("version", "1.0")
  var (
    sign, address string
    errs          []error
  )
  pKey := FormatPrivateKey(privateKey)
  if sign, err = getRsaSign(body, "RSA2", pKey); err != nil {
    return
  }
  body.Set("sign", sign)
  agent := zgoutils.HttpAgent()
  if !isProd {
    address = zfbSandboxBaseUrlUtf8
  } else {
    address = zfbBaseUrlUtf8
  }
  if _, bytes, errs = agent.Post(address).Type("form-data").SendString(FormatURLParam(body)).EndBytes(); len(errs) > 0 {
    return nil, errs[0]
  }
  return
}

//alipay.user.info.share(支付宝会员授权信息查询接口)
//    body：此接口无需body参数
//    文档地址：https://docs.open.alipay.com/api_2/alipay.user.info.share
func (a *PayClient) UserInfoShare() (tradeRes *UserInfoShareResponse, err error) {
  var bs []byte
  if bs, err = a.do(nil, "alipay.user.info.share"); err != nil {
    return nil, err
  }
  tradeRes = new(UserInfoShareResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayUserInfoShareResponse.Code != "10000" {
    info := tradeRes.AlipayUserInfoShareResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.open.auth.token.app(换取应用授权令牌)
//    文档地址：https://docs.open.alipay.com/api_9/alipay.open.auth.token.app
func (a *PayClient) OpenAuthTokenApp(body zgoutils.BodyMap) (tradeRes *OpenAuthTokenAppResponse, err error) {
  var bs []byte
  grantType := body.Get("grant_type")
  if grantType == null {
    return nil, errors.New("grant_type is not allowed to be null")
  }
  code := body.Get("code")
  refreshToken := body.Get("refresh_token")
  if code == null && refreshToken == null {
    return nil, errors.New("code and refresh_token are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "alipay.open.auth.token.app"); err != nil {
    return
  }
  tradeRes = new(OpenAuthTokenAppResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayOpenAuthTokenAppResponse.Code != "10000" {
    info := tradeRes.AlipayOpenAuthTokenAppResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//zhima.credit.score.get(芝麻分)
//    文档地址：https://docs.open.alipay.com/api_8/zhima.credit.score.get
func (a *PayClient) ZhimaCreditScoreGet(body zgoutils.BodyMap) (tradeRes *ZhimaCreditScoreGetResponse, err error) {
  var (
    p1, p2 string
    bs     []byte
  )
  if p1 = body.Get("product_code"); p1 == null {
    body.Set("product_code", "w1010100100000000001")
  }
  if p2 = body.Get("transaction_id"); p2 == null {
    return nil, errors.New("transaction_id is not allowed to be null")
  }
  if bs, err = a.do(body, "zhima.credit.score.get"); err != nil {
    return
  }
  tradeRes = new(ZhimaCreditScoreGetResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.ZhimaCreditScoreGetResponse.Code != "10000" {
    info := tradeRes.ZhimaCreditScoreGetResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//向支付宝发送请求
func (a *PayClient) do(body zgoutils.BodyMap, method string) (bytes []byte, err error) {
  var (
    bodyStr, sign, address, urlParam string
    bodyBs                           []byte
    res                              gorequest.Response
    errs                             []error
  )
  if body != nil {
    if bodyBs, err = zgoutils.Utils.Marshal(body); err != nil {
      return nil, fmt.Errorf("zgoutils.Utils.Marshal：%v", err.Error())
    }
    bodyStr = string(bodyBs)
  }
  pubBody := make(zgoutils.BodyMap)
  pubBody.Set("app_id", a.AppId)
  pubBody.Set("method", method)
  pubBody.Set("format", "JSON")
  if a.AppCertSN != null {
    pubBody.Set("app_cert_sn", a.AppCertSN)
  }
  if a.AlipayRootCertSN != null {
    pubBody.Set("alipay_root_cert_sn", a.AlipayRootCertSN)
  }
  if a.ReturnUrl != null {
    pubBody.Set("return_url", a.ReturnUrl)
  }
  if a.Charset == null {
    pubBody.Set("charset", "utf-8")
  } else {
    pubBody.Set("charset", a.Charset)
  }
  if a.SignType == null {
    pubBody.Set("sign_type", "RSA2")
  } else {
    pubBody.Set("sign_type", a.SignType)
  }
  pubBody.Set("timestamp", time.Now().Format(TimeLayout))
  pubBody.Set("version", "1.0")
  if a.NotifyUrl != null {
    pubBody.Set("notify_url", a.NotifyUrl)
  }
  if a.AppAuthToken != null {
    pubBody.Set("app_auth_token", a.AppAuthToken)
  }
  if a.AuthToken != null {
    pubBody.Set("auth_token", a.AuthToken)
  }
  if bodyStr != null {
    pubBody.Set("biz_content", bodyStr)
  }
  if sign, err = getRsaSign(pubBody, pubBody.Get("sign_type"), FormatPrivateKey(a.PrivateKey)); err != nil {
    return
  }
  pubBody.Set("sign", sign)
  urlParam = FormatURLParam(pubBody)
  if method == "alipay.trade.app.pay" {
    return []byte(urlParam), nil
  }
  if method == "alipay.trade.page.pay" {
    if !a.IsProd {
      return []byte(zfbSandboxBaseUrl + "?" + urlParam), nil
    } else {
      return []byte(zfbBaseUrl + "?" + urlParam), nil
    }
  }
  agent := zgoutils.HttpAgent()
  if !a.IsProd {
    address = zfbSandboxBaseUrlUtf8
  } else {
    address = zfbBaseUrlUtf8
  }
  if res, bytes, errs = agent.Post(address).Type("form-data").SendString(urlParam).EndBytes(); len(errs) > 0 {
    return nil, errs[0]
  }
  if res.StatusCode != 200 {
    return nil, fmt.Errorf("HTTP Request Error, StatusCode = %v", res.StatusCode)
  }
  if method == "alipay.trade.wap.pay" {
    if res.Request.URL.String() == zfbSandboxBaseUrl || res.Request.URL.String() == zfbBaseUrl {
      return nil, errors.New("alipay.trade.wap.pay error,please check the parameters")
    }
    return []byte(res.Request.URL.String()), nil
  }
  return
}

//	AppId      string `json:"app_id"`      //支付宝分配给开发者的应用ID
//	Method     string `json:"method"`      //接口名称
//	Format     string `json:"format"`      //仅支持 JSON
//	ReturnUrl  string `json:"return_url"`  //HTTP/HTTPS开头字符串
//	Charset    string `json:"charset"`     //请求使用的编码格式，如utf-8,gbk,gb2312等，推荐使用 utf-8
//	SignType   string `json:"sign_type"`   //商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用 RSA2
//	Sign       string `json:"sign"`        //商户请求参数的签名串
//	Timestamp  string `json:"timestamp"`   //发送请求的时间，格式"yyyy-MM-dd HH:mm:ss"
//	Version    string `json:"version"`     //调用的接口版本，固定为：1.0
//	NotifyUrl  string `json:"notify_url"`  //支付宝服务器主动通知商户服务器里指定的页面http/https路径。
//	BizContent string `json:"biz_content"` //业务请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入文档

type OpenApiRoyaltyDetailInfoPojo struct {
  RoyaltyType  string `json:"royalty_type,omitempty"`
  TransOut     string `json:"trans_out,omitempty"`
  TransOutType string `json:"trans_out_type,omitempty"`
  TransInType  string `json:"trans_in_type,omitempty"`
  TransIn      string `json:"trans_in"`
  Amount       string `json:"amount,omitempty"`
  Desc         string `json:"desc,omitempty"`
}

//设置 应用公钥证书SN
//    appCertSN：应用公钥证书SN，通过 gopay.GetCertSN() 获取
func (a *PayClient) SetAppCertSN(appCertSN string) {
  a.AppCertSN = appCertSN
}

//设置 支付宝根证书SN
//    alipayRootCertSN：支付宝根证书SN，通过 gopay.GetCertSN() 获取
func (a *PayClient) SetRootCertSN(rootCertSN string) {
  a.AlipayRootCertSN = rootCertSN
}

//设置支付后的ReturnUrl
func (a *PayClient) SetReturnUrl(url string) {
  a.ReturnUrl = url
}

//设置支付宝服务器主动通知商户服务器里指定的页面http/https路径。
func (a *PayClient) SetNotifyUrl(url string) {
  a.NotifyUrl = url
}

//设置编码格式，如utf-8,gbk,gb2312等，默认推荐使用 utf-8
func (a *PayClient) SetCharset(charset string) {
  if charset == null {
    a.Charset = "utf-8"
  } else {
    a.Charset = charset
  }
}

//设置签名算法类型，目前支持RSA2和RSA，默认推荐使用 RSA2
func (a *PayClient) SetSignType(signType string) {
  if signType == null {
    a.SignType = "RSA2"
  } else {
    a.SignType = signType
  }
}

//设置应用授权
func (a *PayClient) SetAppAuthToken(appAuthToken string) {
  a.AppAuthToken = appAuthToken
}

//设置用户信息授权
func (a *PayClient) SetAuthToken(authToken string) {
  a.AuthToken = authToken
}

//获取参数签名
func getRsaSign(bm zgoutils.BodyMap, signType, privateKey string) (sign string, err error) {
  var (
    block          *pem.Block
    h              hash.Hash
    key            *rsa.PrivateKey
    hashs          crypto.Hash
    encryptedBytes []byte
  )

  if block, _ = pem.Decode([]byte(privateKey)); block == nil {
    return null, errors.New("pem.Decode：privateKey decode error")
  }
  if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
    fmt.Println("bodyStr==", err)

    return
  }
  switch signType {
  case "RSA":
    h = sha1.New()
    hashs = crypto.SHA1
  case "RSA2":
    h = sha256.New()
    hashs = crypto.SHA256
  default:
    h = sha256.New()
    hashs = crypto.SHA256
  }
  if _, err = h.Write([]byte(bm.EncodeAliPaySignParams())); err != nil {
    return
  }
  if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, key, hashs, h.Sum(nil)); err != nil {
    return
  }
  sign = base64.StdEncoding.EncodeToString(encryptedBytes)
  return
}

//格式化请求URL参数
func FormatURLParam(body zgoutils.BodyMap) (urlParam string) {
  v := url.Values{}
  for key, value := range body {
    v.Add(key, value.(string))
  }
  urlParam = v.Encode()
  return
}

func getSignData(bs []byte) (signData string) {
  str := string(bs)
  indexStart := strings.Index(str, `":`)
  indexEnd := strings.Index(str, `,"sign"`)
  signData = str[indexStart+2 : indexEnd]
  return
}
