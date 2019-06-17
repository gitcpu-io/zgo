package zgohttp

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"runtime"
	"strconv"
	"time"
)

func getRequestLogs(ctx context.Context) string {
	var status, ip, method, path string
	status = strconv.Itoa(ctx.GetStatusCode())
	path = ctx.Path()
	method = ctx.Method()
	ip = ctx.RemoteAddr()
	// the date should be logged by iris' Logger, so we skip them
	return fmt.Sprintf("%v %s %s %s", status, path, method, ip)
}

type Httper interface {
	JsonpOK(ctx iris.Context, r interface{}) (int, error)
	JsonpErr(ctx iris.Context, msg string) (int, error)
	JsonOK(ctx iris.Context, r interface{}) (int, error)
	JsonServiceErr(ctx iris.Context) (int, error)
	JsonParamErr(ctx iris.Context) (int, error)
	JsonErr(ctx iris.Context, status int, code string, msg string) (int, error)
	JsonExpectErr(ctx iris.Context, msg string) (int, error)
	UseBefore(ctx iris.Context) // 捕获异常，开始计时时间
	AsyncMid(ctx iris.Context)  // 使用go程异步
}

type zgohttp struct {
}

func New() Httper {
	return &zgohttp{}
}

type ErrResponse struct {
	Status    int                    `json:"code"`
	Msg       string                 `json:"message"`
	ErrorCode string                 `json:"errorCode"`
	Data      map[string]interface{} `json:"data"`
	Time      int64                  `json:"time"`
}

var (
	ErrorRequestBodyParseFailed = ErrResponse{Status: 400, Msg: "Request body is not correct", ErrorCode: "001"}
	ErrorDBError                = ErrResponse{Status: 500, Msg: "DB ops failed", ErrorCode: "003"}
	ErrorInternalFaults         = ErrResponse{Status: 500, Msg: "Internal service error", ErrorCode: "004"}
)

func (zh *zgohttp) JsonpOK(ctx iris.Context, r interface{}) (int, error) {
	startTime := ctx.Values().GetInt64Default("startTime", time.Now().UnixNano())
	takeTime := (time.Now().UnixNano() - startTime) / 1e6
	return ctx.JSONP(iris.Map{"code": 200, "data": r, "message": "", "time": takeTime})
}

func (zh *zgohttp) JsonpErr(ctx iris.Context, msg string) (int, error) {
	startTime := ctx.Values().GetInt64Default("startTime", time.Now().UnixNano())
	takeTime := (time.Now().UnixNano() - startTime) / 1e6
	return ctx.JSONP(iris.Map{"code": 400, "message": msg, "data": make(map[string]interface{}), "time": takeTime})
}

// JsonOK 正常的返回方法
func (zh *zgohttp) JsonOK(ctx iris.Context, r interface{}) (int, error) {
	startTime := ctx.Values().GetInt64Default("startTime", time.Now().UnixNano())
	takeTime := (time.Now().UnixNano() - startTime) / 1e6
	return ctx.JSON(iris.Map{"code": 200, "data": r, "message": "操作成功", "time": takeTime})
}

// JsonExpectErr 预期内的错误，适用于调用func后 return出来的errors!=nil时的返回值
func (zh *zgohttp) JsonExpectErr(ctx iris.Context, msg string) (int, error) {
	startTime := ctx.Values().GetInt64Default("startTime", time.Now().UnixNano())
	takeTime := (time.Now().UnixNano() - startTime) / 1e6
	return ctx.JSON(ErrResponse{Status: 500, Msg: msg, ErrorCode: "500", Data: make(map[string]interface{}), Time: takeTime})
}

// JsonOtherErr 其他自定义返回方法 （业务本身的异常)
func (zh *zgohttp) JsonErr(ctx iris.Context, status int, code string, msg string) (int, error) {
	startTime := ctx.Values().GetInt64Default("startTime", time.Now().UnixNano())
	takeTime := (time.Now().UnixNano() - startTime) / 1e6
	return ctx.JSON(ErrResponse{Status: status, Msg: msg, ErrorCode: code, Data: make(map[string]interface{}), Time: takeTime})
}

// JsonServiceErr defer recover到panic的时候用的异常方法
func (zh *zgohttp) JsonServiceErr(ctx iris.Context) (int, error) {
	msg := "服务器开小差了，稍后再试吧"
	startTime := ctx.Values().GetInt64Default("startTime", time.Now().UnixNano())
	takeTime := (time.Now().UnixNano() - startTime) / 1e6
	return ctx.JSON(ErrResponse{Status: 500, Msg: msg, ErrorCode: "500", Data: make(map[string]interface{}), Time: takeTime})
}

// JsonParamErr 参数验证不通过时调用
func (zh *zgohttp) JsonParamErr(ctx iris.Context) (int, error) {
	msg := "参数错误"
	startTime := ctx.Values().GetInt64Default("startTime", time.Now().UnixNano())
	takeTime := (time.Now().UnixNano() - startTime) / 1e6
	return ctx.JSON(ErrResponse{Status: 400, Msg: msg, ErrorCode: "400", Data: make(map[string]interface{}), Time: takeTime})
}

func (zh *zgohttp) UseBefore(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			//fmt.Println(err)
			zh.JsonServiceErr(ctx)
			if ctx.IsStopped() {
				return
			}
			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break
				}
				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}
			// when stack finishes
			logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
			logMessage += fmt.Sprintf("At Request: %s\n", getRequestLogs(ctx))
			logMessage += fmt.Sprintf("Trace: %s\n", err)
			logMessage += fmt.Sprintf("\n%s", stacktrace)
			ctx.Application().Logger().Warn(logMessage)
			//ctx.StatusCode(500)
			ctx.StopExecution()
		}
	}()
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,x-token")
	ctx.Header("content-type", "application/json")
	start := time.Now().UnixNano()
	ctx.Values().Set("startTime", start)
	ctx.Next()
}

func (zh *zgohttp) AsyncMid(ctx iris.Context) {
	ch := make(chan int)
	go func() {
		ctx.Next()
		ch <- 1
	}()
	select {
	case <-ch:
		return
	}
}
