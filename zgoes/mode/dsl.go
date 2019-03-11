package mode

import "git.zhugefang.com/gocore/zgo/zgoutils"

var DSL = map[string]string{
	"bool": `{"source":"{\"query\":{\"bool\":{\"must\":{{#toJson}}must{{/toJson}},\"should\":{{#toJson}}should{{/toJson}},\"must_not\":{{#toJson}}must_not{{/toJson}} }}}}","params":{"must":%s,"should":%s,"must_not":%s}}`,
}

var SimpleDsl = map[string]interface{}{
	"source": `{\"query\": {{#toJson}}simple{{/toJson}}}`,
}

var Term = func(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["term"] = m
	return term
	//t, _ := zgoutils.Utils.Marshal(term)
	//return t
}
var Range = func(f string, m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["range"] = map[string]interface{}{
		f: m,
	}

	return term
	//t, _ := zgoutils.Utils.Marshal(term)
	//return t
}

var Match = func(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["match"] = m
	return term
	//t, _ := zgoutils.Utils.Marshal(term)
	//return t
}

var MatchPhrase = func(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["match_phrase"] = m
	return term
	//t, _ := zgoutils.Utils.Marshal(term)
	//return t
}

var Wildcard = func(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["wildcard"] = m
	return term
	//t, _ := zgoutils.Utils.Marshal(term)
	//return t
}

var GeoBox = func(f string, top_left map[string]interface{},
	bottom_right map[string]interface{}) interface{} {

	term := map[string]interface{}{}
	term["geo_bounding_box"] = map[string]interface{}{
		f: map[string]interface{}{
			"top_left":     top_left,
			"bottom_right": bottom_right,
		},
	}
	return term
	//t, _ := zgoutils.Utils.Marshal(term)
	//return t
}
