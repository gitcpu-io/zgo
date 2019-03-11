package zgohttp

import (
	"github.com/kataras/iris"
)

type Httper interface {
	JsonpOK(ctx iris.Context, r interface{}) (int, error)
	JsonpErr(ctx iris.Context, msg string) (int, error)
	JsonServiceErr(ctx iris.Context) (int, error)
	JsonParamErr(ctx iris.Context) (int, error)
	JsonOtherErr(ctx iris.Context, code int, msg string) (int, error)
	JsonExpectErr(ctx iris.Context, msg string) (int, error)
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
	return ctx.JSONP(iris.Map{"status": 400, "msg": msg})
}

// JsonOK 正常的返回方法
func (zh *zgohttp) JsonOK(ctx iris.Context, r interface{}) (int, error) {
	return ctx.JSON(iris.Map{"status": 200, "data": r})
}

// JsonExpectErr 预期内的错误，适用于调用func后 return出来的errors!=nil时的返回值
func (zh *zgohttp) JsonExpectErr(ctx iris.Context, msg string) (int, error) {
	return ctx.JSON(ErrResponse{Status: 500, Msg: msg, ErrorCode: "500"})
}

// JsonOtherErr 其他自定义返回方法 （不要轻易使用）
func (zh *zgohttp) JsonOtherErr(ctx iris.Context, code int, msg string) (int, error) {
	return ctx.JSON(ErrResponse{Status: code, Msg: msg, ErrorCode: string(code)})
}

// JsonServiceErr defer recover到panic的时候用的异常方法
func (zh *zgohttp) JsonServiceErr(ctx iris.Context) (int, error) {
	msg := "服务器异常"
	return ctx.JSON(ErrResponse{Status: 500, Msg: msg, ErrorCode: "500"})
}

// JsonParamErr 参数验证不通过时调用
func (zh *zgohttp) JsonParamErr(ctx iris.Context) (int, error) {
	msg := "参数错误"
	return ctx.JSON(ErrResponse{Status: 400, Msg: msg, ErrorCode: "400"})
}
