/*
@Time : 2019-02-26 12:23
@Author : zhangjianguo
@File : service
@Software: GoLand
*/
package zgoes

import (
	"context"
	"git.zhugefang.com/gocore/zgo/comm"
	"git.zhugefang.com/gocore/zgo/config"
	"sync"
)

var (
	currentLabels = make(map[string][]*config.ConnDetail)
	muLabel       sync.RWMutex
)

//项目初始化  根据用户选择label 初始化Es实例
func InitEs(hsm map[string][]*config.ConnDetail) chan *zgoes {
	muLabel.Lock()
	defer muLabel.Unlock()
	currentLabels = hsm

	//自动为变量初始化对象
	initLabel := ""
	for k, _ := range hsm {
		if k != "" {
			initLabel = k
			break
		}
	}
	out := make(chan *zgoes)
	go func() {
		in, err := GetEs(initLabel)
		if err != nil {
			out <- nil
		}
		out <- in
		close(out)
	}()
	return out

}

type zgoes struct {
	res EsResourcer //使用resource另外的一个接口
}

//GetMongo zgo内部获取一个连接mongo
func GetEs(label ...string) (*zgoes, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgoes{
		res: NewEsResourcer(l), //interface
	}, nil
}

func Es(l string) Eser {
	return &zgoes{
		res: NewEsResourcer(l),
	}
}

//Es 对外
type Eser interface {
	NewEs(label ...string) (*zgoes, error) //初始化方法
	SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error)
	QueryTmp(ctx context.Context, index, table, tmp string, args map[string]interface{}) (interface{}, error)
}

func (e *zgoes) NewEs(label ...string) (*zgoes, error) {
	return GetEs(label...)
}
func (e *zgoes) SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error) {
	return e.res.SearchDsl(ctx, index, table, dsl, args)
}

func (e *zgoes) QueryTmp(ctx context.Context, index, table, tmp string, args map[string]interface{}) (interface{}, error) {
	return e.res.QueryTmp(ctx, index, table, tmp, args)

}
