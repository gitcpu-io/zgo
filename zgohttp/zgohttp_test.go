package zgohttp

import (
	"git.zhugefang.com/gocore/zgo.git/zgoresponse"
	"net/http"
	"testing"
)

func TestRouter_GET(t *testing.T) {
	router := New()
	//router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	//router.GET("/mu/:name", Hello1)

	//router.GET("/hello/mu/:name", Hello2)

	http.ListenAndServe(":8080", router)
}

//func Index(w http.ResponseWriter, r *http.Request, _ Params) {
//	fmt.Fprint(w, "Welcome!\n")
//}

//func Hello(w http.ResponseWriter, r *http.Request, ps Params) {
//	info := fmt.Sprintln(r.Header.Get("Content-Type"))
//	len := r.ContentLength
//	body := make([]byte, len)
//	r.Body.Read(body)
//	fmt.Println(w, 	r.URL.Query())
//	fmt.Fprintln(w, info, string(body))
//	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
//}

func Hello(w http.ResponseWriter, r *http.Request, ps Params) {
	//m, _ := ParseReq2Map(r)
	//ret := zgoresponse.NewServerResponse(200, m)
	ret := zgoresponse.NewServerError(404, 404, "not found")
	Json(w, ret)
	//fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

//func Hello1(w http.ResponseWriter, r *http.Request, ps Params) {
//	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
//}
//
//func Hello2(w http.ResponseWriter, r *http.Request, ps Params) {
//	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
//}
