package zgo

import (
	"git.zhugefang.com/gocore/zgo.git/logic/zgo_es"
)

var Es *zgo_es.Es

func init() {
	Es = zgo_es.NewES()
}
