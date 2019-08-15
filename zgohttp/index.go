package zgohttp

import (
	"bufio"
	"bytes"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"io/ioutil"
	"net/http"
	"net/url"
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
	JsonFree(ctx iris.Context, content interface{}) (int, error) // 自定义返回结构体
	UseBefore(ctx iris.Context)                                  // 捕获异常，开始计时时间
	AsyncMid(ctx iris.Context)                                   // 使用go程异步
	Get(url string) ([]byte, error)
	Post(url string, play url.Values) ([]byte, error)
	PostJson(url string, jsonData []byte, handler ...map[string]string) ([]byte, error)
	PostForm(url string, formData []byte) ([]byte, error)
	//钉钉机器人
	Ding(token string, msg string)
	GetByProxy(httpUrl string, proxyAddr string, params map[string]interface{}) ([]byte, error)
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

func (u *zgohttp) Ding(token string, msg string) {
	url := "https://oapi.dingtalk.com/robot/send?access_token=" + token
	mps := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": msg},
	}
	bytes, _ := zgoutils.Utils.Marshal(mps)
	u.PostJson(url, bytes)
	return
}

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

// JsonOK 正常的返回方法
func (zh *zgohttp) JsonFree(ctx iris.Context, content interface{}) (int, error) {
	//startTime := ctx.Values().GetInt64Default("startTime", time.Now().UnixNano())
	//takeTime := (time.Now().UnixNano() - startTime) / 1e6
	return ctx.JSON(content)
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

func (zh *zgohttp) Get(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(bufio.NewReader(resp.Body))
	if err != nil {
		return nil, err
	}
	return body, err
}

func (zh *zgohttp) Post(url string, play url.Values) ([]byte, error) {
	resp, err := http.PostForm(url, play)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func (zh *zgohttp) PostJson(url string, jsonData []byte, handler ...map[string]string) ([]byte, error) {

	reader := bytes.NewReader([]byte(jsonData))
	request, err := http.NewRequest("POST", url, reader)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	if len(handler) > 0 {
		for k, v := range handler[0] {
			request.Header.Set(k, v)
		}
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer request.Body.Close()
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return respBytes, nil
}

func (zh *zgohttp) PostForm(url string, jsonData []byte) ([]byte, error) {

	var tmp interface{}

	err := zgoutils.Utils.Unmarshal(jsonData, &tmp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	tmpMap, ok := tmp.(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("Got (%v) is not JSON Map", string(jsonData))
	}

	var formData = make(map[string][]string)
	for k, v := range tmpMap {
		formData[k] = []string{fmt.Sprintf("%v", v)}
	}

	client := http.Client{}
	resp, err := client.PostForm(url, formData)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return respBytes, nil
}

func (zh *zgohttp) GetByProxy(httpUrl string, proxyAddr string, params map[string]interface{}) ([]byte, error) {
	proxy, err := url.Parse("http://" + proxyAddr)
	if err != nil {
		return nil, err
	}
	netTransport := &http.Transport{
		Proxy:               http.ProxyURL(proxy),
		MaxIdleConnsPerHost: 10,
	}
	httpClient := &http.Client{
		Transport: netTransport,
	}
	resp, err := httpClient.Get(httpUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(bufio.NewReader(resp.Body))
	if err != nil {
		return nil, err
	}
	return body, err
}
