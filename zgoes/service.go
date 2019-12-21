/*
  Elasticsearch 客户端 基于http实现 可以执行原生DSL语句
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
func InitEs(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgoes {
	muLabel.Lock()
	defer muLabel.Unlock()

	var hsm map[string][]*config.ConnDetail

	if len(label) > 0 && len(currentLabels) > 0 { //此时是destory操作,传入的hsm是nil
		//fmt.Println("--destory--前",currentLabels)
		for _, v := range label {
			delete(currentLabels, v)
		}
		hsm = currentLabels
		//fmt.Println("--destory--后",currentLabels)

	} else { //这是第一次创建操作或etcd中变更时init again操作
		hsm = hsmIn
		//currentLabels = hsm	//this operation is error
		for k, v := range hsm { //so big bug can't set hsm to currentLabels，must be for, may be have old label
			currentLabels[k] = v
		}
	}

	if len(hsm) == 0 {
		return nil
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
	NewDsl() *DSL
	AddOneData(ctx context.Context, index, table, id, dataJson string) (interface{}, error)
	UpOneData(ctx context.Context, index, table, id, dataJson string) (interface{}, error)
	DeleteDsl(ctx context.Context, index, table, dsl string) (interface{}, error)
	UpDateByQuery(ctx context.Context, index, table, dsl string) (interface{}, error)
	ExistsIndices(ctx context.Context, index, table string) (bool, error)
	CreateIndices(ctx context.Context, index, table string) (bool, error)
	AddOneDataAutoId(ctx context.Context, index, table, dataJson string) ([]byte, error)
	QueryDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) ([]byte, error)
	UpdateByQuery(ctx context.Context, index, table, dsl string) ([]byte, error)
	CreateOneData(ctx context.Context, index, table, id, dataJson string) ([]byte, error)
	ScrollSearch(ctx context.Context, index, table, dsl string, scrollId string) ([]byte, error)
}

func (e *zgoes) New(label ...string) (*zgoes, error) {
	return GetEs(label...)
}

func (e *zgoes) SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error) {
	return e.res.SearchDsl(ctx, index, table, dsl, args)
}

func (e *zgoes) AddOneData(ctx context.Context, index, table, id, dsl string) (interface{}, error) {
	return e.res.AddOneData(ctx, index, table, id, dsl)
}
func (e *zgoes) UpOneData(ctx context.Context, index, table, id, dataJson string) (interface{}, error) {
	return e.res.UpOneData(ctx, index, table, id, dataJson)
}

func (e *zgoes) DeleteDsl(ctx context.Context, index, table, dsl string) (interface{}, error) {
	return e.res.DeleteDsl(ctx, index, table, dsl)
}

func (e *zgoes) UpDateByQuery(ctx context.Context, index, table, dsl string) (interface{}, error) {
	return e.res.UpDateByQuery(ctx, index, table, dsl)
}

func (e *zgoes) ExistsIndices(ctx context.Context, index, table string) (bool, error) {
	return e.res.ExistsIndices(ctx, index, table)
}
func (e *zgoes) CreateIndices(ctx context.Context, index, table string) (bool, error) {
	return e.res.CreateIndices(ctx, index, table)
}

func (e *zgoes) AddOneDataAutoId(ctx context.Context, index, table, dataJson string) ([]byte, error) {
	return e.res.AddOneDataAutoId(ctx, index, table, dataJson)
}

func (e *zgoes) QueryDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) ([]byte, error) {
	return e.res.QueryDsl(ctx, index, table, dsl, args)
}

func (e *zgoes) UpdateByQuery(ctx context.Context, index, table, dsl string) ([]byte, error) {
	return e.res.UpdateByQuery(ctx, index, table, dsl)
}

func (e *zgoes) CreateOneData(ctx context.Context, index, table, id, dsl string) ([]byte, error) {
	return e.res.CreateOneData(ctx, index, table, id, dsl)
}

func (e *zgoes) ScrollSearch(ctx context.Context, index, table, dsl string, scrollId string) ([]byte, error) {
	return e.res.ScrollSearch(ctx, index, table, dsl, scrollId)
}

//func (e *zgoes)NewDsl() *mode.DSL {
//	return &mode.DSL{
//		querys: make(map[string]interface{}),
//		args:   make(map[string]interface{}),
//	}
//}

func (e *zgoes) NewDsl() *DSL {
	return NewDSL()
}
