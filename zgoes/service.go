package zgoes

import (
	"context"
	"sync"
)

var (
	currentLabels = make(map[string][]string)
	muLabel       sync.RWMutex
)

type Eser interface {
	Add(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error)
	Del(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error)
	Set(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error)
	Get(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) ([]interface{}, error)
	Search(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error)
}

type zgoes struct {
	res EsResourcer //使用resource另外的一个接口
}

func (e *zgoes) Add(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return e.res.Add(ctx, index, table, dsl, args)
}

func (e *zgoes) Del(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return e.res.Del(ctx, index, table, dsl, args)
}

func (e *zgoes) Set(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return e.res.Set(ctx, index, table, dsl, args)
}

func (e *zgoes) Get(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return e.res.Get(ctx, index, table, dsl, args)
}

func (e *zgoes) Search(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return e.res.Search(ctx, index, table, dsl, args)
}
