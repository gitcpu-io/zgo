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
	Search(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error)
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary


type esResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}


func (e *esResource) Add(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	//es := NewEsResource()

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
	//es := NewEsResource()

	maps := map[string]interface{}{} //定义es返回结构提

	return maps, nil
}

func (e *esResource) Search(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	es := NewEsResource()            //获取zgo封装的的client
	maps := map[string]interface{}{} //定义es返回结构提
	url := "http://101.201.28.195:9200"
	uri := url + "/" + index + "/" + table + "/" + "_search?pretty"
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(dsl))
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := es.GetEsClient().Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&maps); err != nil {
		fmt.Print(err)
		return nil, err
	}
	return maps, err
}
