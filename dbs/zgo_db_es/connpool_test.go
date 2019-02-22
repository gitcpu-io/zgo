package zgo_db_es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	uri := "http://101.201.28.195:9200/active_bj_house_sell/spider/_search?pretty"
	dsl := `{"query": {"match_all": {}}}`
	er := NewEsResource()
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(dsl))
	req.Header.Set("Content-Type", "application/json")

	resp, err := er.GetEsClient().Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	r, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(r))
	mp := make(map[string]interface{})
	_ = json.Unmarshal(r, &mp)
	fmt.Println(mp)
}
