package zgoes

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"io/ioutil"
	"strings"

	//jsoniter "github.com/json-iterator/go"
	"net/http"
	"sync"
)

//对外接口
type EsResourcer interface {
	SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error)
	AddOneData(ctx context.Context, index, table, id, dataJson string) (interface{}, error)
	UpOneData(ctx context.Context, index, table, id, dataJson string) (interface{}, error)
	DeleteDsl(ctx context.Context, index, table, dsl string) (interface{}, error)
	UpDateByQuery(ctx context.Context, index, table, dsl string) (interface{}, error)
}

var mu sync.RWMutex

//接口实现
type esResource struct {
	label string       //配置标签
	mu    sync.RWMutex //读写锁
	uri   string       //绑定地址
}

//获取原生http
func (e *esResource) GetConChan() *http.Client {
	return &http.Client{}
}

//方法初始化从uris中获取uri
func NewEsResourcer(label string) EsResourcer {
	mu.RLock()
	defer mu.RUnlock()
	var uri = ""
	if al, ok := currentLabels[label]; ok {
		lf := al[0]
		// uri = lf.Uri
		uri = fmt.Sprintf("http://%s:%s@%s:%v", lf.Username, lf.Password, lf.Host, lf.Port)
	}
	return &esResource{
		label: label,
		uri:   uri,
	}
}

//根据dsl语句执行查询
func (e *esResource) SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{}
	//定义es结果集返回结构体
	uri := e.uri + "/" + index + "/" + table + "/" + "_search?pretty"
	if ignoreUnavailable, ok := args["ignoreUnavailable"]; ok {
		uri = uri + "&ignore_unavailable=" + ignoreUnavailable.(string)
	}
	//拼接es请求uti[索引+文档+_search]
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(dsl)) //post请求
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json") //设置json协议解析头
	resp, err := e.GetConChan().Do(req)                //获取绑定的地址执行请求
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer resp.Body.Close()

	be, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	if err := zgoutils.Utils.Unmarshal(be, &maps); err != nil {
		fmt.Print(err)
		return nil, err
	}

	return maps, err
}

//根据dsl语句执行删除
func (e *esResource) DeleteDsl(ctx context.Context, index, table, dsl string) (interface{}, error) {
	maps := map[string]interface{}{}
	//定义es结果集返回结构体
	uri := e.uri + "/" + index + "/" + table + "/" + "_delete_by_query"
	//拼接es请求uti[索引+文档+_delete]
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(dsl)) //post请求
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json") //设置json协议解析头
	resp, err := e.GetConChan().Do(req)                //获取绑定的地址执行请求
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer resp.Body.Close()

	be, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	if err := zgoutils.Utils.Unmarshal(be, &maps); err != nil {
		fmt.Print(err)
		return nil, err
	}

	return maps, err
}

func (e *esResource) AddOneData(ctx context.Context, index, table, id, dataJson string) (interface{}, error) {
	uri := e.uri + "/" + index + "/" + table + "/" + id
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(dataJson)) //post请求
	if err != nil {
		return nil, fmt.Errorf("es add data create request error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := e.GetConChan().Do(req)
	if err != nil {
		return nil, fmt.Errorf("es add data post error: %v", err)
	}
	defer resp.Body.Close()

	be, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("es add data read body error: %v", err)
	}
	var result interface{}

	if err := zgoutils.Utils.Unmarshal(be, &result); err != nil {
		return nil, fmt.Errorf("es add data umarshal error: %v", err)
	}

	return result, err
}

func (e *esResource) UpOneData(ctx context.Context, index, table, id, dataJson string) (interface{}, error) {
	uri := e.uri + "/" + index + "/" + table + "/" + id + "/" + "_update"
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(dataJson)) //post请求
	if err != nil {
		return nil, fmt.Errorf("es Up data create request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.GetConChan().Do(req)
	if err != nil {
		return nil, fmt.Errorf("es Up data post error: %v", err)
	}
	defer resp.Body.Close()

	be, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("es Up data read body error: %v", err)
	}
	var result interface{}

	if err := zgoutils.Utils.Unmarshal(be, &result); err != nil {
		return nil, fmt.Errorf("es Up data umarshal error: %v", err)
	}
	return result, err
}

func (e *esResource) UpDateByQuery(ctx context.Context, index, table, dsl string) (interface{}, error) {
	uri := e.uri + "/" + index + "/" + table + "/" + "_update_by_query"
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(dsl)) //post请求
	if err != nil {
		return nil, fmt.Errorf("es Up data create request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.GetConChan().Do(req)
	if err != nil {
		return nil, fmt.Errorf("es Up data post error: %v", err)
	}
	defer resp.Body.Close()

	be, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("es Up data read body error: %v", err)
	}
	var result interface{}

	if err := zgoutils.Utils.Unmarshal(be, &result); err != nil {
		return nil, fmt.Errorf("es Up data umarshal error: %v", err)
	}
	return result, err
}