package zgoalipay

import (
  "errors"
  "fmt"
  "github.com/gitcpu-io/zgo/zgoutils"
)

/*
@Time : 2019-10-18 17:46
@Author : rubinus.chu
@File : fund
@project: zgo
*/

//alipay.fund.trans.uni.transfer(统一转账到支付宝账户接口)
//    文档地址：https://docs.open.alipay.com/api_28/alipay.fund.trans.uni.transfer/
func (a *PayClient) FundTransUniTransfer(body zgoutils.BodyMap) (tradeRes *FundTransUniTransferResponse, err error) {
  var bs []byte
  var p1, p2 string
  p1 = body.Get("out_biz_no")
  p2 = body.Get("trans_amount")
  if p1 == null && p2 == null {
    return nil, errors.New("out_biz_no and trans_amount are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "alipay.fund.trans.uni.transfer"); err != nil {
    return
  }
  tradeRes = new(FundTransUniTransferResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayFundTransUniTransferResponse.Code != "10000" {
    info := tradeRes.AlipayFundTransUniTransferResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.fund.trans.common.query(查询转账订单接口)
//	文档地址：https://docs.open.alipay.com/api_28/alipay.fund.trans.common.query
func (a *PayClient) FundTransCommonQuery(body zgoutils.BodyMap) (tradeRes *FundTransCommonQueryResponse, err error) {
  var bs []byte
  var p1 = body.Get("out_biz_no")
  if p1 == null {
    return nil, errors.New("out_biz_no are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "alipay.fund.trans.common.query"); err != nil {
    return
  }
  tradeRes = &FundTransCommonQueryResponse{}
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayFundTransCommonQueryResponse.Code != "10000" {
    info := tradeRes.AlipayFundTransCommonQueryResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.fund.trans.toaccount.transfer(单笔转账到支付宝账户接口)
//    文档地址：https://docs.open.alipay.com/api_28/alipay.fund.trans.toaccount.transfer
func (a *PayClient) FundTransToaccountTransfer(body zgoutils.BodyMap) (tradeRes *FundTransToaccountTransferResponse, err error) {
  var bs []byte
  trade := body.Get("out_biz_no")
  if trade == null {
    return nil, errors.New("out_biz_no is not allowed to be null")
  }
  if bs, err = a.do(body, "alipay.fund.trans.toaccount.transfer"); err != nil {
    return
  }
  tradeRes = new(FundTransToaccountTransferResponse)
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayFundTransToaccountTransferResponse.Code != "10000" {
    info := tradeRes.AlipayFundTransToaccountTransferResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.fund.trans.order.query(查询转账订单接口)
//	文档地址：https://docs.open.alipay.com/api_28/alipay.fund.trans.order.query
func (a *PayClient) FundTransOrderQuery(body zgoutils.BodyMap) (tradeRes *FundTransOrderQueryResponse, err error) {
  var bs []byte
  var p1, p2 string
  p1 = body.Get("out_biz_no")
  p2 = body.Get("order_id")
  if p1 == null && p2 == null {
    return nil, errors.New("out_biz_no and order are not allowed to be null at the same time")
  }
  if bs, err = a.do(body, "alipay.fund.trans.order.query"); err != nil {
    return
  }
  tradeRes = &FundTransOrderQueryResponse{}
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayFundTransOrderQueryResponse.Code != "10000" {
    info := tradeRes.AlipayFundTransOrderQueryResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.fund.account.query(支付宝资金账户资产查询接口)
//	文档地址：https://docs.open.alipay.com/api_28/alipay.fund.account.query
func (a *PayClient) FundAccountQuery(body zgoutils.BodyMap) (tradeRes *FundAccountQueryResponse, err error) {
  var bs []byte
  var (
    userId             string
    accountProductCode string
    accountSceneCode   string
  )
  userId = body.Get("alipay_user_id")
  if userId == null {
    return nil, errors.New("alipay_user_id not allowed to be null")
  }

  accountProductCode = body.Get("account_product_code")
  accountSceneCode = body.Get("account_scene_code")
  if accountProductCode == null && accountSceneCode == null {
    return nil, errors.New("account_product_code and account_scene_code are not allowed to be null at the same time")
  }

  if bs, err = a.do(body, "alipay.fund.account.query"); err != nil {
    return
  }
  tradeRes = &FundAccountQueryResponse{}
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayFundAccountQueryResponse.Code != "10000" {
    info := tradeRes.AlipayFundAccountQueryResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}

//alipay.fund.trans.refund(资金退回接口)
//	文档地址：https://docs.open.alipay.com/api_28/alipay.fund.trans.refund/
func (a *PayClient) FundTransRefund(body zgoutils.BodyMap) (tradeRes *FundTansRefundResponse, err error) {
  var bs []byte
  var (
    p1 string
    p2 string
    p3 string
  )
  p1 = body.Get("order_id")
  if p1 == null {
    return nil, errors.New("order_id not allowed to be null")
  }

  p2 = body.Get("out_request_no")
  p3 = body.Get("refund_amount")
  if p2 == null || p3 == null {
    return nil, errors.New("out_request_no or refund_amount are not allowed to be null at the same time")
  }

  if bs, err = a.do(body, "alipay.fund.trans.refund"); err != nil {
    return
  }
  tradeRes = &FundTansRefundResponse{}
  if err = zgoutils.Utils.Unmarshal(bs, tradeRes); err != nil {
    return nil, err
  }
  if tradeRes.AlipayFundTansRefundResponse.Code != "10000" {
    info := tradeRes.AlipayFundTansRefundResponse
    return nil, fmt.Errorf(`{"code":"%v","msg":"%v","sub_code":"%v","sub_msg":"%v"}`, info.Code, info.Msg, info.SubCode, info.SubMsg)
  }
  tradeRes.SignData = getSignData(bs)
  return
}
