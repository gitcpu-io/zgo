package mode

import "git.zhugefang.com/gocore/zgo/zgoutils"

var DSL = map[string]string{
	"bool": `{"source":"{\"query\":{\"bool\":{\"must\":{{#toJson}}must{{/toJson}},\"should\":{{#toJson}}should{{/toJson}},\"must_not\":{{#toJson}}must_not{{/toJson}} }}}}","params":{"must":%s,"should":%s,"must_not":%s}}`,
}

var Term = func(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["term"] = m
	t, _ := zgoutils.Utils.Marshal(term)
	return t
}
var Rang = func(f string, m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term[f] = m
	t, _ := zgoutils.Utils.Marshal(term)
	return t
}
var Match = func(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["term"] = m
	t, _ := zgoutils.Utils.Marshal(term)
	return t
}
