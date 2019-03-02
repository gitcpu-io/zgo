/*
@Time : 2019-02-26 12:23
@Author : zhangjianguo
@File : resource
@Software: GoLand
*/
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
	SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error)
	QueryTmp(ctx context.Context, index, table, tmp string, args map[string]interface{}) (interface{}, error)
}

var mu sync.RWMutex

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//方法初始化从uris中获取uri
func NewEsResourcer(label string) EsResourcer {

	//get hosts by label
	mu.RLock()
	defer mu.RUnlock()
	//var hosts []*config.ConnDetail
	//if al, ok := currentLabels[label]; ok {
	//	for _, v := range al {
	//		hosts = append(hosts, v)
	//	}
	//}
	var uri = ""
	if al, ok := currentLabels[label]; ok {
		lf := al[0]
		uri = lf.Uri
	}
	return &esResource{
		label: label,
		//hosts: la.Uri,
		uri: uri,
	}
}

type esResource struct {
	label string
	mu    sync.RWMutex
	//hosts []*config.ConnDetail
	uri string
}

func (e *esResource) GetConChan() *http.Client {
	return &http.Client{}
}

/*
@Time : 2019-02-26 12:23
@Author : zhangjianguo
@File : service
@Software: GoLand
@parms: index:索引名称
@parms: doc:文档类型
@parms: dsl:es原生语句
*/
func (e *esResource) SearchDsl(ctx context.Context, index, table, dsl string, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{} //定义es返回结构提
	uri := e.uri + "/" + index + "/" + table + "/" + "_search?pretty"
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

func (e *esResource) QueryTmp(ctx context.Context, index, table, tmp string, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{} //定义es返回结构提
	uri := e.uri + "/" + index + "/" + table + "/" + "_search/template?pretty"
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(tmp))
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
