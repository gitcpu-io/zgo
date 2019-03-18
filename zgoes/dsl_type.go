package zgoes

import (
	"git.zhugefang.com/gocore/zgo/zgoutils"
)

type DSL struct {
	args   map[string]interface{}
	querys map[string]interface{}
}

func NewDSL() *DSL {
	return &DSL{
		querys: make(map[string]interface{}),
		args:   make(map[string]interface{}),
	}
}

func (dsl *DSL) addOp(op string, term ...interface{}) *DSL {
	val, ok := dsl.querys[op]
	var querys []interface{}

	if !ok {
		querys = make([]interface{}, 0)
	} else {
		querys = val.([]interface{})
	}

	querys = append(querys, term...)
	dsl.querys[op] = querys
	return dsl
}

func (dsl *DSL) Set(key string, value interface{}) *DSL {
	dsl.args[key] = value
	return dsl
}

func (dsl *DSL) SetAggs(value interface{}) *DSL {
	dsl.Set("aggs", value)
	return dsl
}

func (dsl *DSL) SetFrom(value interface{}) *DSL {
	dsl.Set("from", value)
	return dsl
}

func (dsl *DSL) SetSize(value interface{}) *DSL {
	dsl.Set("size", value)
	return dsl
}

// value type []string
func (dsl *DSL) Set_Source(value interface{}) *DSL {
	dsl.Set("_source", value)
	return dsl
}

// val type string // val 可为多值 field1 field2...
func (dsl *DSL) Set_SourceField(val ...interface{}) *DSL {
	dsl.Set("_source", val)
	return dsl
}

func (dsl *DSL) SetSort(value interface{}) *DSL {
	dsl.Set("sort", value)
	return dsl
}

func (dsl *DSL) Must(term ...interface{}) *DSL {
	return dsl.addOp("must", term...)
}

func (dsl *DSL) MustNot(term ...interface{}) *DSL {
	return dsl.addOp("must_not", term...)
}

func (dsl *DSL) Filter(term ...interface{}) *DSL {
	return dsl.addOp("filter", term...)
}

func (dsl *DSL) Should(term ...interface{}) *DSL {
	return dsl.addOp("should", term...)
}

// m = {field: value}
func (dsl *DSL) TremMap(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["term"] = m
	return term
}

//
func (dsl *DSL) TermField(field string, value interface{}) interface{} {
	term := map[string]interface{}{}
	term["term"] = map[string]interface{}{
		field: value,
	}
	return term
}

// m = { field: values} // values = [val0, val1, val2, ...]
func (dsl *DSL) TermsMap(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["terms"] = m
	return term
}

// values = [val0, val1, val2, ...]
func (dsl *DSL) TermsField(field string, values interface{}) interface{} {
	term := map[string]interface{}{}
	term["terms"] = map[string]interface{}{
		field: values,
	}
	return term
}

// m = {op1: val1, op2: val2}  op = lt gt lte gte
func (dsl *DSL) RangeMap(field string, m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["range"] = map[string]interface{}{
		field: m,
	}
	return term
}

//
func (dsl *DSL) RangeField(field string, op1 string, val1 interface{},
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
func (dsl *DSL) MatchMap(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["match"] = m
	return term
}

//
func (dsl *DSL) MatchField(field string, value interface{}) interface{} {
	term := map[string]interface{}{}
	term["match"] = map[string]interface{}{
		field: value,
	}
	return term
}

// m = {field: value} or {field: {"query": value, "slop": slop}}
func (dsl *DSL) MatchPhraseMap(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["match_phrase"] = m
	return term
}

func (dsl *DSL) MatchPhraseField(field string, value interface{}) interface{} {
	term := map[string]interface{}{}
	term["match_phrase"] = map[string]interface{}{
		field: value,
	}
	return term
}

func (dsl *DSL) MatchPhraseSlop(field string, value interface{}, slop int) interface{} {
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
func (dsl *DSL) WildcardMap(m map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["wildcard"] = m
	return term
}

//
func (dsl *DSL) WildcardField(field string, value interface{}) interface{} {
	term := map[string]interface{}{}
	term["wildcard"] = map[string]interface{}{
		field: value,
	}
	return term
}

//
func (dsl *DSL) GeoBox(f string, top_left map[string]interface{},
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
func (dsl *DSL) GeoBoxMap(location map[string]interface{}) interface{} {
	term := map[string]interface{}{}
	term["geo_bounding_box"] = location
	return term
}

//
func (dsl *DSL) GeoBoxField(field string, left_lat, left_lon,
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

// field : field
// return {"aggs_field": {"terms":{"field": field,"size": 10}}}
func (dsl *DSL) SimpleAggs(field string, size int) interface{} {
	return map[string]interface{}{
		"aggs_field": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": field,
				"size":  size,
			},
		},
	}
}

// field : path.field
// return {"nested_field": {"nested":{"path": path}, "aggs":{"aggs_field":{"terms":{"field": path.field, "size": size}}}}}
func (dsl *DSL) NestedAggs(path, field string, size int) interface{} {
	return map[string]interface{}{
		"nested_field": map[string]interface{}{
			"nested": map[string]interface{}{
				"path": path,
			},
			"aggs": map[string]interface{}{
				"aggs_field": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": field,
						"size":  size,
					},
				},
			},
		},
	}
}

func (dsl *DSL) Aggs2Field(field1 string, size1 int, field2 string, size2 int) interface{} {
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

func (dsl *DSL) SimpleSort(field string, isAsc bool) interface{} {
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

func (dsl *DSL) BoolQuery(op string, mapWhere interface{}) interface{} {
	opMap := make(map[string]interface{})
	opMap[op] = mapWhere

	if op == "should" {
		opMap["minimum_should_match"] = 1
	}

	return map[string]interface{}{
		"bool": opMap,
	}
}

// mapWhere = []interface{} or map[string] interface{}
// return {bool: {must: mapWhere}}
func (dsl *DSL) BoolMust(mapWhere interface{}) interface{} {
	return dsl.BoolQuery("must", mapWhere)
}

// mapWhere = []interface{}
// return {bool: {should: mapWhere}, "minimum_should_match":1}
func (dsl *DSL) BoolShould(mapWhere interface{}) interface{} {
	return dsl.BoolQuery("should", mapWhere)
}

// mapWhere = []interface{} or map[string] interface{}
// return {bool: {must_not: mapWhere}}
func (dsl *DSL) BoolMustNot(mapWhere interface{}) interface{} {
	return dsl.BoolQuery("must_not", mapWhere)
}

// mapWhere = []interface{} or map[string] interface{}
// return {bool: {filter: mapWhere}}
func (dsl *DSL) BoolFilter(mapWhere interface{}) interface{} {
	return dsl.BoolQuery("filter", mapWhere)
}

// mapMixWhere = {"must": mapWhere, "must_not": mapWhere, "filter": mapWhere, "should": mapWhere}
// return {bool: {"filter": mapWhere, "must":...}}
func (dsl *DSL) BoolMix(mapMixWhere map[string]interface{}) interface{} {
	if len(mapMixWhere) > 0 {
		return map[string]interface{}{
			"bool": mapMixWhere,
		}
	}
	return mapMixWhere
}

func marshal(T interface{}) string {
	t, _ := zgoutils.Utils.Marshal(T)
	return string(t)
}

// path = path to nested
// boolTerm = {"bool": {"must" : {"term": {path.field: value}}}}
// return { nested: {"path": path, "query": boolTerm}}
func (dsl *DSL) NestedDslTerm(path string, boolTerm interface{}) interface{} {
	return map[string]interface{}{
		"nested": map[string]interface{}{
			"path":  path,
			"query": boolTerm,
		},
	}
}

func (dsl *DSL) BoolDslTerm() interface{} {
	if len(dsl.querys) > 0 {
		return map[string]interface{}{
			"bool": dsl.querys,
		}
	}
	return dsl.querys
}

func (dsl *DSL) BoolDslString() string {
	return marshal(dsl.BoolDslTerm())
}

func (dsl *DSL) QueryDsl() string {

	_, ok := dsl.querys["should"]
	if ok {
		dsl.querys["minimum_should_match"] = 1
	}

	if len(dsl.querys) > 0 {
		dsl.args["query"] = map[string]interface{}{
			"bool": dsl.querys,
		}
	}

	return marshal(dsl.args)
}
