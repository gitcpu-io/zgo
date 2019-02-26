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
	"git.zhugefang.com/gocore/zgo.git/config"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
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
	Match(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

var mu sync.RWMutex

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//方法初始化从uris中获取uri
func NewEsResourcer(label string) EsResourcer {
	//get hosts by label
	mu.RLock()
	defer mu.RUnlock()

	var hosts []*config.ConnDetail
	if al, ok := currentLabels[label]; ok {
		for _, v := range al {
			hosts = append(hosts, v)
		}
	}

	return &esResource{
		label: label,
		hosts: hosts,
	}
}

type esResource struct {
	label string
	mu    sync.RWMutex
	hosts []*config.ConnDetail
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
func (e *esResource) Add(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{} //定义es返回结构提
	return maps, nil
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
func (e *esResource) Del(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{} //定义es返回结构提
	return maps, nil
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
func (e *esResource) DelByQuery(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{} //定义es返回结构提
	return maps, nil
}

/*
@Time : 2019-02-26 12:23
@Author : zhangjianguo
@File : service
@Software: GoLand
@des :
@parms: index:索引名称
@parms: doc:文档类型
@parms: dsl:es原生语句
*/
func (e *esResource) DelById(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{} //定义es返回结构提
	return maps, nil
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
func (e *esResource) Set(ctx context.Context, args map[string]interface{}) (interface{}, error) {

	maps := map[string]interface{}{} //定义es返回结构提

	return maps, nil
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
func (e *esResource) Get(ctx context.Context, args map[string]interface{}) (interface{}, error) {

	maps := map[string]interface{}{} //定义es返回结构提
	return maps, nil
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
func (e *esResource) Search(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{} //定义es返回结构提
	index := args["index"].(string)
	table := args["table"].(string)
	dsl := args["dsl"].(string)

	var url string
	r := e.hosts[rand.Intn(1)]
	if len(e.hosts) > 0 {
		//从数组中随机选择一个
		url = e.hosts[0].Uri
	}
	fmt.Println(rand.Intn(1), r, url, "000")

	uri := url + "/" + index + "/" + table + "/" + "_search?pretty"
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

/*
@Time : 2019-02-26 12:23
@Author : zhangjianguo
@File : service
@Software: GoLand
@parms: index:索引名称
@parms: doc:文档类型
@parms: dsl:es原生语句
*/
func (e *esResource) bulk(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return nil, nil
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
func (e *esResource) Match(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	maps := map[string]interface{}{} //定义es返回结构提
	return maps, nil

}
