/*
  Elasticsearch 客户端 基于http实现 可以执行原生DSL语句
*/
package zgoes

import (
	"context"
	"git.zhugefang.com/gocore/zgo/comm"
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/zgoes/mode"
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
	//currentLabels = hsm

	for k, v := range hsm { //so big bug can't set hsm to currentLabels，must be for, may be have old label
		currentLabels[k] = v
	}

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
	res EsResourcer
}

func GetEs(label ...string) (*zgoes, error) {
	//根据配置获取具体配置信息
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgoes{
		res: NewEsResourcer(l),
	}, nil
}

func Es(l string) Eser {
	return &zgoes{
		res: NewEsResourcer(l),
	}
}

/*
 ElasticSearch 对外使用接口
*/
type Eser interface {
	// 根据配置名称获取Elastic实例 如果所在项目中只使用一个Elastic实例时，则无需初始化（调用NewEs）,可以直接使用接口
	New(label ...string) (*zgoes, error)
	// param ctx:上线文
	// param index:索引文明
	// param table:文档名称
	// param dsl: 原生elastic语句
	// 根据elastic dsl 语句查询数据 该接口只能执行查询操作
	SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error)
	NewDsl() *mode.DSL
}

func (e *zgoes) New(label ...string) (*zgoes, error) {
	return GetEs(label...)
}

func (e *zgoes) SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error) {
	return e.res.SearchDsl(ctx, index, table, dsl, args)
}

//func (e *zgoes)NewDsl() *mode.DSL {
//	return &mode.DSL{
//		querys: make(map[string]interface{}),
//		args:   make(map[string]interface{}),
//	}
//}

func (e *zgoes) NewDsl() *mode.DSL {
	return mode.NewDSL()
}
