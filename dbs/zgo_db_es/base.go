package zgo_db_es

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Search(ctx context.Context, index string, table string, dsl string, args map[string]interface{}) (interface{}, error) {
	es := NewEsResource()
	url := "http://101.201.28.195:9200"
	uri := url + "/" + index + "/" + table + "_search?pretty"
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(dsl))
	req.Header.Set("Content-Type", "application/json")
	resp, err := es.GetEsClient().Do(req)
	defer resp.Body.Close()
	s, _ := ioutil.ReadAll(resp.Body)
	inte := make(map[string]interface{})
	_ = json.Unmarshal(s, &inte)
	return inte, err
}
