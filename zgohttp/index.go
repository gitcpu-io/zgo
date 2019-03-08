package zgohttp

import (
	"github.com/kataras/iris"
)

type Httper interface {
	JsonpOK(ctx iris.Context, r interface{}) (int, error)
	JsonpErr(ctx iris.Context, msg string) (int, error)
}

type zgohttp struct {
}

func NewHttp() Httper {
	return &zgohttp{}
}

type ErrResponse struct {
	Status    int    `json:"status"`
	Msg       string `json:"msg"`
	ErrorCode string `json:"errorCode"`
}

var (
	ErrorRequestBodyParseFailed = ErrResponse{Status: 400, Msg: "Request body is not correct", ErrorCode: "001"}
	ErrorDBError                = ErrResponse{Status: 500, Msg: "DB ops failed", ErrorCode: "003"}
	ErrorInternalFaults         = ErrResponse{Status: 500, Msg: "Internal service error", ErrorCode: "004"}
)

func (zh *zgohttp) JsonpOK(ctx iris.Context, r interface{}) (int, error) {
	return ctx.JSONP(iris.Map{"status": 200, "data": r})
}

func (zh *zgohttp) JsonpErr(ctx iris.Context, msg string) (int, error) {
	return ctx.JSONP(iris.Map{"status": 201, "msg": msg})
}
