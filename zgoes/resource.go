package zgoes

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strings"
	"sync"
)

type EsResourcer interface {
	Add(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Del(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Set(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Get(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Search(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var uris = map[string]string{} //定义es返回结构提

//项目初始化加载配置文件
func InitEsResource(hsm map[string][]string) {
	uris["sell_write"] = "http://101.201.28.195:9200"
}

//方法初始化从uris中获取uri
func NewEsResourcer(label string) EsResourcer {
	return &esResource{
		label: label,
		url:   uris[label],
	}
}

type esResource struct {
	label string
	mu    sync.RWMutex
	url   string
}

func (e *esResource) GetConChan() *http.Client {
	return &http.Client{}
}

func (e *esResource) Add(ctx context.Context, args map[string]interface{}) (interface{}, error) {

	maps := map[string]interface{}{} //定义es返回结构提

	return maps, nil
}

func (e *esResource) Del(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	//es := NewEsResource()

	maps := map[string]interface{}{} //定义es返回结构提

	return maps, nil
}

func (e *esResource) Set(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	//es := NewEsResource()

	maps := map[string]interface{}{} //定义es返回结构提

	return maps, nil
}

func (e *esResource) Get(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{} //定义es返回结构提
	return maps, nil
}

func (e *esResource) Search(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{} //定义es返回结构提
	index := args["index"].(string)
	table := args["table"].(string)
	dsl := args["dsl"].(string)

	uri := e.url + "/" + index + "/" + table + "/" + "_search?pretty"
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(dsl))
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.GetConChan().Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&maps); err != nil {
		fmt.Print(err)
		return nil, err
	}
	return maps, nil
}
