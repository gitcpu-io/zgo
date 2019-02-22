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

func (e *Es) Search(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	return zgo_db_es.Search(ctx, index, table, dsl, args)
}
