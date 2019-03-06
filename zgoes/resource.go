package zgoes

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"io/ioutil"

	//jsoniter "github.com/json-iterator/go"
	"net/http"
	"strings"
	"sync"
)

type EsResourcer interface {
	SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error)
}

var mu sync.RWMutex

//var json = jsoniter.ConfigCompatibleWithStandardLibrary

//方法初始化从uris中获取uri
func NewEsResourcer(label string) EsResourcer {
	mu.RLock()
	defer mu.RUnlock()
	var uri = ""
	if al, ok := currentLabels[label]; ok {
		lf := al[0]
		uri = lf.Uri
	}
	return &esResource{
		label: label,
		uri:   uri,
	}
}

type esResource struct {
	label string
	mu    sync.RWMutex
	uri   string
}

func (e *esResource) GetConChan() *http.Client {
	return &http.Client{}
}

func (e *esResource) SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{}
	//定义es返回结构体=
	uri := e.uri + "/" + index + "/" + table + "/" + "_search?pretty"         //拼接es请求uti[索引+文档+_search]
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

