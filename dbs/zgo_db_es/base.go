package zgo_db_es

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Search(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	es := NewEsResource()            //获取zgo封装的的client
	maps := map[string]interface{}{} //定义es返回结构提
	url := "http://101.201.28.195:9200"
	uri := url + "/" + index + "/" + table + "/" + "_search?pretty"
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(dsl))
	if err != nil {
		fmt.Print(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := es.GetEsClient().Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Print(err)
	}
	if err := json.NewDecoder(resp.Body).Decode(&maps); err != nil {
		fmt.Print(err)
	}
	return maps, err
}
