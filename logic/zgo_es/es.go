package zgo_es

import (
	"context"
	"git.zhugefang.com/gocore/zgo.git/dbs/zgo_db_es"
)

type Es struct {
}

func NewES() *Es {
	return &Es{}
}

func (e *Es) Add(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return zgo_db_es.Add(ctx, index, table, dsl, args)
}

func (e *Es) Del(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return zgo_db_es.Del(ctx, index, table, dsl, args)
}

func (e *Es) Set(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return zgo_db_es.Set(ctx, index, table, dsl, args)
}

func (e *Es) Get(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return zgo_db_es.Get(ctx, index, table, dsl, args)
}

func (e *Es) List(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return zgo_db_es.List(ctx, index, table, dsl, args)
}

func (e *Es) Search(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return zgo_db_es.Search(ctx, index, table, dsl, args)
}
