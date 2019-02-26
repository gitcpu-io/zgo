/*
@Time : 2019-02-26 12:23
@Author : zhangjianguo
@File : service
@Software: GoLand
*/
package zgoes

import (
	"context"
	"git.zhugefang.com/gocore/zgo.git/comm"
	"git.zhugefang.com/gocore/zgo.git/config"
	"sync"
)

var (
	currentLabels = make(map[string][]config.ConnDetail)
	muLabel       sync.RWMutex
)

//项目初始化  根据用户选择label 初始化Es实例
func InitEs(hsm map[string][]config.ConnDetail) {
	muLabel.Lock()
	defer muLabel.Unlock()
	currentLabels = hsm
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
	Add(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Del(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Set(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Get(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Search(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

func (e *zgoes) NewEs(label ...string) (*zgoes, error) {
	return GetEs(label...)
}

func (e *zgoes) Add(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return e.res.Add(ctx, args)
}

func (e *zgoes) Del(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return e.res.Del(ctx, args)
}

func (e *zgoes) Set(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return e.res.Set(ctx, args)
}

func (e *zgoes) Get(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return e.res.Get(ctx, args)
}

func (e *zgoes) Search(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return e.res.Search(ctx, args)
}
