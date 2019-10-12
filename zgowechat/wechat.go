package zgowechat

/*
@Time : 2019-10-11 15:19
@Author : rubinus.chu
@File : wechat
@project: zgo
*/

type Wechat interface {
	Pay(appId, mchId, apiKey string, isProd bool) Payer

	//添加其它接口

}

type wechat struct {
}

func (w *wechat) Pay(appId, mchId, apiKey string, isProd bool) Payer {
	return NewPayClient(appId, mchId, apiKey, isProd)
}

func New() Wechat {
	return &wechat{}
}
