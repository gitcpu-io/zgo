package mode

import (
	"fmt"

	"git.zhugefang.com/gocore/zgo/zgoutils"
)

var DSL = map[string]string{
	"bool": `{"source":"{\"query\":{\"bool\":{\"must\":{{#toJson}}must{{/toJson}},\"should\":{{#toJson}}should{{/toJson}},\"must_not\":{{#toJson}}must_not{{/toJson}} }}}}","params":{"must":%s,"should":%s,"must_not":%s}}`,
}

var boolTemplate = `{"source":"{\"query\":{\"bool\":{\"must\":{{#toJson}}must{{/toJson}},\"should\":{{#toJson}}should{{/toJson}},\"must_not\":{{#toJson}}must_not{{/toJson}},\"filter\":{{#toJson}}filter{{/toJson}}}},\"sort\":{{#toJson}}sort{{/toJson}},\"_source\":{{#toJson}}_source{{/toJson}},\"from\": {{from}},\"size\": {{size}},\"aggs\":{{#toJson}}aggs{{/toJson}}}",
"params":{"must":%s,"should":%s,"must_not":%s, "filter": %s,"sort": %s,"_source":%s,"aggs":%s,"from": %d,"size":%d}}`

var emptyMap = map[string]interface{}{}

//var emptySlice = []string{}

var template = `{
	"source":"{\"query\":{{#toJson}}query{{/toJson}},\"sort\":{{#toJson}}sort{{/toJson}},\"_source\":{{#toJson}}_source{{/toJson}},\"from\": {{from}},\"size\": {{size}},\"aggs\":{{#toJson}}aggs{{/toJson}}}",
	"params":{
		"query": %s,
		"sort": %s,
		"_source":%s,
		"aggs":%s,
		"from": %d,
		"size":%d
	}
}`

func ifelse(input, defvalue string) string {
	if input == "" {
		return defvalue
	}
	return input
}

func GetSimpleDsl(query, sort, _source, aggs string, from, size int) string {
	var t_sort, t_source, t_aggs string

	t_sort = ifelse(sort, "{}")
	t_source = ifelse(_source, "[]")
	t_aggs = ifelse(aggs, "{}")

	return fmt.Sprintf(template, query, t_sort, t_source, t_aggs, from, size)
}

func GetBoolDsl(query, sort, _source, aggs string, from, size int) string {
	var t_sort, t_source, t_aggs string

	t_sort = ifelse(sort, "{}")
	t_source = ifelse(_source, "[]")
	t_aggs = ifelse(aggs, "{}")

	return fmt.Sprintf(boolTemplate, query, t_sort, t_source, t_aggs, from, size)
}

func parseArgs(args map[string]interface{}, key string) (interface{}, bool) {
	values, ok := args[key]
	if !ok {
		return emptyMap, false
	}
	return values, true
}

func QueryDsl(args map[string]interface{}) interface{} {
	boolMap := make(map[string]interface{})

	must, ok := parseArgs(args, "must")
	if ok {
		boolMap["must"] = must
	}

	must_not, ok := parseArgs(args, "must_not")
	if ok {
		boolMap["must_not"] = must_not
	}
	filter, ok := parseArgs(args, "filter")
	if ok {
		boolMap["filter"] = filter
	}

	should, ok := parseArgs(args, "should")
	if ok {
		boolMap["should"] = should
		boolMap["minimum_should_match"] = 1
	}

	queryMap := make(map[string]interface{})
	sort, ok := parseArgs(args, "sort")
	if ok {
		queryMap["sort"] = sort
	}
	aggs, ok := parseArgs(args, "aggs")
	if ok {
		queryMap["aggs"] = aggs
	}
	_source, ok := parseArgs(args, "_source")
	if ok {
		queryMap["_source"] = _source
	}
	from, ok := args["from"]
	if ok {
		queryMap["from"] = from
	}
	size, ok := args["size"]
	if ok {
		queryMap["size"] = size
	}
	if len(boolMap) > 0 {
		queryMap["query"] = map[string]interface{}{
			"bool": boolMap,
		}
	}

	return unmarshal(queryMap)
}

var NestedElem = `{
	"nested" :{
		"path": %q,
		"query": %s
	}
}`

var NestedSimpleAggs = `{
	"nested_field":{
		"nested":{
			"path": %q
		},
		"aggs": {
			"aggs_field" : {
				"terms":{
					"field": %q
				}
			}
		}
	}
}`

// m = {field: value}
var TermMap = func(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["term"] = m
	return term
}

//
var TermField = func(field string, value interface{}) interface{} {
	term := map[string]interface{}{}
	term["term"] = map[string]interface{}{
		field: value,
	}
	return term
}

// m = {op1: val1, op2: val2}  op = lt gt lte gte
var RangeMap = func(field string, m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["range"] = map[string]interface{}{
		field: m,
	}
	return term
}

//
var RangeField = func(field string, op1 string, val1 interface{},
	op2 string, val2 interface{}) interface{} {
	valMap := make(map[string]interface{})
	if op1 != "" && val1 != nil {
		valMap[op1] = val1
	}

	if op2 != "" && val2 != nil {
		valMap[op2] = val2
	}

	if len(valMap) == 0 {
		return valMap
	}

	term := map[string]interface{}{}
	term["range"] = map[string]interface{}{
		field: valMap,
	}
	return term
}

// m = {field: value}
var MatchMap = func(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["match"] = m
	return term
}

//
var MatchField = func(field string, value interface{}) interface{} {
	term := map[string]interface{}{}
	term["match"] = map[string]interface{}{
		field: value,
	}
	return term
}

// m = {field: value} or {field: {"query": value, "slop": slop}}
var MatchPhraseMap = func(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["match_phrase"] = m
	return term
}

var MatchPhraseField = func(field string, value interface{}) interface{} {
	term := map[string]interface{}{}
	term["match_phrase"] = map[string]interface{}{
		field: value,
	}
	return term
}

var MatchPhraseSlop = func(field string, value interface{}, slop int) interface{} {
	term := map[string]interface{}{}
	term["match_phrase"] = map[string]interface{}{
		field: map[string]interface{}{
			"query": value,
			"slop":  slop,
		},
	}
	return term
}

// m = { field: value}
var WildcardMap = func(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["wildcard"] = m
	return term
}

//
var WildcardField = func(field string, value interface{}) interface{} {
	term := map[string]interface{}{}
	term["wildcard"] = map[string]interface{}{
		field: value,
	}
	return term
}

//
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
}

//
var GeoBoxMap = func(location map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["geo_bounding_box"] = location
	return term
}

//
var GeoBoxField = func(field string, left_lat, left_lon,
	right_lat, right_lon float64) interface{} {

	term := map[string]interface{}{}
	term["geo_bounding_box"] = map[string]interface{}{
		field: map[string]interface{}{
			"top_left": map[string]float64{
				"lat": left_lat,
				"lon": left_lon,
			},
			"bottom_right": map[string]float64{
				"lat": right_lat,
				"lon": right_lon,
			},
		},
	}
	return term
}

func SimpleAggs(field string, size int) interface{} {
	//return fmt.Sprintf(`{"aggs_field": {"terms":{"field": %q,"size": %d}}}`, field, size)
	return map[string]interface{}{
		"aggs_field": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": field,
				"size":  size,
			},
		},
	}
}

func Aggs2Field(field1 string, size1 int, field2 string, size2 int) interface{} {
	return map[string]interface{}{
		"aggs_field": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": field1,
				"size":  size1,
			}, "aggs": map[string]interface{}{
				"aggs_field": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": field2,
						"size":  size2,
					},
				},
			},
		},
	}
}

func SimpleSort(field string, isAsc bool) interface{} {
	var order = "asc"
	if !isAsc {
		order = "desc"
	}
	return map[string]interface{}{
		field: map[string]interface{}{
			"order": order,
		},
	}
}

func BoolQuery(op string, mapWhere interface{}) interface{} {
	opMap := make(map[string]interface{})
	opMap[op] = mapWhere

	if op == "should" {
		opMap["minimum_should_match"] = 1
	}

	return map[string]interface{}{
		"bool": opMap,
	}
}

func MustQuery(mapWhere []interface{}) interface{} {
	return BoolQuery("must", mapWhere)
}

func ShouldQuery(mapWhere []interface{}) interface{} {
	return BoolQuery("should", mapWhere)
}

func MustNotQuery(mapWhere []interface{}) interface{} {
	return BoolQuery("must_not", mapWhere)
}

func FilterQuery(mapWhere []interface{}) interface{} {
	return BoolQuery("filter", mapWhere)
}

func unmarshal(T interface{}) interface{} {
	t, _ := zgoutils.Utils.Marshal(T)
	return string(t)
}
