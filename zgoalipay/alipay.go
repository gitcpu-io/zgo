package zgoalipay

/*
@Time : 2019-10-11 15:19
@Author : rubinus.chu
@File : alipay
@project: zgo
*/

type AliPay interface {
	Pay(appId, privateKey string, isProd bool) Payer

	//添加其它接口

}

type alipay struct {
}

func (p *alipay) Pay(appId, privateKey string, isProd bool) Payer {
	return NewPayClient(appId, privateKey, isProd)
}

func New() AliPay {
	return &alipay{}
}
