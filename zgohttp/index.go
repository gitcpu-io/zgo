package zgohttp

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoresponse"
	"github.com/json-iterator/go"
	"net/http"
	"net/url"
	"strings"
)

type Httper interface {
	ParseReq2Map(r *http.Request) (map[string]interface{}, error)
	Json(w http.ResponseWriter, response zgoresponse.Response)
	JsonP(w http.ResponseWriter, response zgoresponse.Response)
}

type zgohttp struct {
}

func NewHttp() Httper {
	return &zgohttp{}
}

func (zh *zgohttp) ParseReq2Map(r *http.Request) (map[string]interface{}, error) {
	//info := fmt.Sprintln(r.Header.Get("Content-Type"))
	//len := r.ContentLength
	paramMap := make(map[string]interface{})
	if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "json") {
		err := jsoniter.NewDecoder(r.Body).Decode(&paramMap)
		if err != nil {
			fmt.Println(err)
		}
		//body := make([]byte, len)
		//r.Body.Read(body)
		queryMap, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			fmt.Println(err)
		}
		for key, value := range queryMap {
			paramMap[key] = value[0]
		}
		return paramMap, nil
	} else if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "form") {
		r.ParseForm()

		for key, value := range r.Form {
			paramMap[key] = value[0]
		}
		return paramMap, nil
	} else {
		queryMap, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			fmt.Println(err)
		}
		for key, value := range queryMap {
			paramMap[key] = value[0]
		}
		return paramMap, nil
	}

}

func (zh *zgohttp) Json(w http.ResponseWriter, response zgoresponse.Response) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	ret, err := jsoniter.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintln(w, string(ret))
}

func (zh *zgohttp) JsonP(w http.ResponseWriter, response zgoresponse.Response) {
	//if response.GetCallBack() == "" {
	//	zh.Json(w, response)
	//	return
	//}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	ret, err := jsoniter.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintln(w, string(ret))
}
