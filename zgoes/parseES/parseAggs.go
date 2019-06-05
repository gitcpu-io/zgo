package parseES

import (
	"errors"
	"fmt"
	"log"
)

func DoBucket(cell AggCell) AggRes {
	var field, nestedName, name, filterName string
	name = cell.AggName()

	if cell.Path == "" {
		field = cell.Field
	} else {
		field = cell.NestedField()
	}

	res := map[string]interface{}{
		name: map[string]interface{}{
			"terms": map[string]interface{}{
				"field": field,
				"size":  cell.AggSize(),
			},
		},
	}

	if cell.Path == "" {
		return AggRes{
			FieldName:  name,
			FilterName: filterName,
			NestedName: nestedName,
			Path:       cell.Path,
			AggTerm:    res,
			IsSon:      cell.IsSon,
		}
	}

	nestedName = cell.NestedName()

	res = map[string]interface{}{
		nestedName: map[string]interface{}{
			"nested": map[string]string{
				"path": cell.Path,
			},
			"aggs": res,
		},
	}
	return AggRes{
		FieldName:  name,
		NestedName: nestedName,
		AggTerm:    res,
		Path:       cell.Path,
		FilterName: filterName,
		IsSon:      cell.IsSon,
	}
}

func DoMetric(cell AggCell) AggRes {
	var field, action, name, nestedName, filterName string
	name = cell.AggName()

	if cell.Path == "" {
		field = cell.Field
	} else {
		field = cell.NestedField()
	}

	if len(cell.Action) > 7 {
		action = cell.Action[7:]
	} else {
		action = cell.Action
	}

	res := map[string]interface{}{
		cell.AggName(): map[string]interface{}{
			action: map[string]string{
				"field": field,
			},
		},
	}

	if cell.Path == "" {
		return AggRes{
			FieldName:  name,
			AggTerm:    res,
			Path:       cell.Path,
			NestedName: nestedName,
			FilterName: filterName,
			IsSon:      cell.IsSon,
		}
	}

	nestedName = cell.NestedName()

	res = map[string]interface{}{
		nestedName: map[string]interface{}{
			"nested": map[string]string{
				"path": cell.Path,
			},
			"aggs": res,
		},
	}
	return AggRes{
		FieldName:  name,
		AggTerm:    res,
		Path:       cell.Path,
		NestedName: nestedName,
		FilterName: filterName,
		IsSon:      cell.IsSon,
	}
}

func DoRanges(cell AggCell) AggRes {
	var field, name, nestedName, filterName string
	name = cell.AggName()

	if cell.Path == "" {
		field = cell.Field
	} else {
		field = cell.NestedField()
	}

	res := map[string]interface{}{
		cell.AggName(): map[string]interface{}{
			"range": map[string]interface{}{
				"field":  field,
				"keyed":  true,
				"ranges": cell.AggRanges(),
			},
		},
	}
	if cell.Path == "" {
		return AggRes{
			FieldName:  name,
			AggTerm:    res,
			Path:       cell.Path,
			NestedName: nestedName,
			FilterName: filterName,
			IsSon:      cell.IsSon,
		}
	}

	nestedName = cell.NestedName()

	res = map[string]interface{}{
		nestedName: map[string]interface{}{
			"nested": map[string]string{
				"path": cell.Path,
			},
			"aggs": res,
		},
	}
	return AggRes{
		FieldName:  name,
		AggTerm:    res,
		Path:       cell.Path,
		NestedName: nestedName,
		FilterName: filterName,
		IsSon:      cell.IsSon,
	}
}

func DoTopHits(cell AggCell) AggRes {
	var n int
	var name, nestedName, filterName string

	res := make(map[string]interface{})
	if parseMap, ok := cell.Query.(map[string]interface{}); ok {
		if sort, ok := parseMap["sort"]; ok {
			if sortMap, ok := sort.(map[string]interface{}); ok {
				res["sort"] = map[string]interface{}{
					sortMap["field"].(string): map[string]interface{}{
						"order": sortMap["type"],
					},
				}
			} else {
				sortMaps := sort.([]map[string]interface{})
				n = len(sortMaps)
				if n == 1 {
					res["sort"] = map[string]interface{}{
						sortMaps[0]["field"].(string): map[string]interface{}{
							"order": sortMaps[0]["type"],
						},
					}
				} else {
					sortRes := make([]map[string]interface{}, 0)
					for i := 0; i < n; i++ {
						sortRes = append(sortRes, map[string]interface{}{
							sortMaps[i]["field"].(string): map[string]interface{}{
								"order": sortMaps[i]["type"],
							},
						})
					}
					res["sort"] = sortRes
				}
			}
		}
		if _source, ok := parseMap["_source"]; ok {
			_sourceSlice := _source.([]string)
			res["_source"] = _sourceSlice
		}
		if size, ok := parseMap["size"]; ok {
			sizeInt := size.(int)
			res["size"] = sizeInt
		}
		if from, ok := parseMap["from"]; ok {
			fromInt := from.(int)
			res["from"] = fromInt
		}
		if to, ok := parseMap["to"]; ok {
			toInt := to.(int)
			res["to"] = toInt
		}
	}
	name = cell.AggName()
	res = map[string]interface{}{
		name: map[string]interface{}{
			"top_hits": res,
		},
	}
	return AggRes{
		FieldName:  name,
		AggTerm:    res,
		Path:       cell.Path,
		NestedName: nestedName,
		FilterName: filterName,
		IsSon:      cell.IsSon,
	}
}

func makeFilter(query map[string]interface{}) map[string]interface{} {
	ret := new(Bool)
	return map[string]interface{}{"bool": ret.Parse(query)}
}

func DoFilter(cell AggCell) AggRes {
	var name, nestedName, filterName string
	filterName = cell.AggName()

	if parseMap, ok := cell.Query.(map[string]interface{}); ok {
		filter := makeFilter(parseMap)
		res := map[string]interface{}{
			filterName: map[string]interface{}{
				"filter": filter,
			},
		}

		return AggRes{
			FieldName:  name,
			AggTerm:    res,
			Path:       cell.Path,
			NestedName: nestedName,
			FilterName: filterName,
			IsSon:      cell.IsSon,
		}
	} else {
		log.Fatalf("Filter Params Error: %v", cell)
		return AggRes{}
	}
}

func FilterBucket(cell AggCell) AggRes {
	var filterName string
	filterName = cell.AggName()

	if parseMap, ok := cell.Query.(map[string]interface{}); ok {
		filter := makeFilter(parseMap)
		cell.Action = "terms"
		bucket := DoBucket(cell)
		res := map[string]interface{}{
			filterName: map[string]interface{}{
				"filter": filter,
				"aggs":   bucket.AggTerm,
			},
		}

		bucket.FilterName = filterName
		bucket.AggTerm = res
		return bucket
	} else {
		log.Fatalf("Filter Params Error: %v", cell)
		return AggRes{}
	}
}

func FilterMetric(cell AggCell) AggRes {
	var filterName string
	filterName = cell.AggName()

	if parseMap, ok := cell.Query.(map[string]interface{}); ok {
		filter := makeFilter(parseMap)
		metric := DoMetric(cell)
		res := map[string]interface{}{
			filterName: map[string]interface{}{
				"filter": filter,
				"aggs":   metric.AggTerm,
			},
		}

		metric.FilterName = filterName
		metric.AggTerm = res
		return metric
	} else {
		log.Fatalf("Filter Params Error: %v", cell)
		return AggRes{}
	}
}

func FilterRanges(cell AggCell) AggRes {
	var filterName string
	filterName = cell.AggName()

	if parseMap, ok := cell.Query.(map[string]interface{}); ok {
		filter := makeFilter(parseMap)
		ranges := DoRanges(cell)
		res := map[string]interface{}{
			filterName: map[string]interface{}{
				"filter": filter,
				"aggs":   ranges.AggTerm,
			},
		}
		ranges.FilterName = filterName
		ranges.AggTerm = res
		return ranges
	} else {
		log.Fatalf("Filter Params Error: %v", cell)
		return AggRes{}
	}
}

func FilterTopHits(cell AggCell) AggRes {
	var filterName string
	filterName = cell.AggName()

	if parseMap, ok := cell.Query.(map[string]interface{}); ok {
		filterMap := parseMap["filter"].(map[string]interface{})
		filter := makeFilter(filterMap)
		delete(parseMap, "filter")
		cell.Query = parseMap
		tophits := DoTopHits(cell)
		res := map[string]interface{}{
			filterName: map[string]interface{}{
				"filter": filter,
				"aggs":   tophits.AggTerm,
			},
		}
		tophits.FilterName = filterName
		tophits.AggTerm = res
		return tophits
	} else {
		log.Fatalf("Filter Params Error: %v", cell)
		return AggRes{}
	}
	return AggRes{}
}

func ProductAggs(aggParams []AggCell) (interface{}, error) {
	N := len(aggParams)
	if N == 0 {
		return nil, errors.New("agg params lens 0")
	}

	aggTerms := make([]AggRes, 0)
	var term AggRes

	for i := 0; i < N; i++ {
		action, err := aggParams[i].AggAction()
		if err != nil {
			return nil, fmt.Errorf("Aggs Param %d, %v", i+1, err)
		}

		switch action {
		case "terms":
			term = DoBucket(aggParams[i])
		case "ranges":
			term = DoRanges(aggParams[i])
		case "filter":
			term = DoFilter(aggParams[i])
		case "top_hits", "tophits":
			term = DoTopHits(aggParams[i])
		case "avg", "min", "max", "sum":
			term = DoMetric(aggParams[i])
		case "filter-terms":
			term = FilterBucket(aggParams[i])
		case "filter-ranges":
			term = FilterRanges(aggParams[i])
		case "filter-top_hits", "filter-tophits":
			term = FilterTopHits(aggParams[i])
		case "filter-avg", "filter-min", "filter-max", "filter-sum":
			term = FilterMetric(aggParams[i])
		}

		aggTerms = append(aggTerms, term)
	}

	if N == 1 {
		// return aggTerms[0].AggTerm, nil
		return aggTerms[0], nil
	} else {
		return CombineTerms(aggTerms), nil
	}
}

func ProductMultiAggs(multiParams [][]AggCell) (interface{}, error) {
	N := len(multiParams)
	ress := make([]AggRes, 0)

	for i := 0; i < N; i++ {
		resTemp, err := ProductAggs(multiParams[i])
		if err != nil {
			return nil, err
		}
		res, ok := resTemp.(AggRes)
		if ok {
			ress = append(ress, res)
		} else {
			return nil, errors.New("product Aggs error, param index: " + smallItoa(i+1))
		}
	}

	if N == 1 {
		return ress[0], nil
	} else {
		return CombineTerms(ress), nil
	}
}

func Parallel(term, tail AggRes) AggRes {
	var name string

	if len(tail.FilterName) > 0 {
		name = tail.FilterName
	} else {
		if len(tail.NestedName) > 0 {
			name = tail.NestedName
		} else {
			name = tail.FieldName
		}
	}

	term.AggTerm[name] = tail.AggTerm[name]
	return term
}

func InClusion(term, tail AggRes) AggRes {
	if len(term.Path) == 0 {
		if len(term.FilterName) == 0 {
			term.AggTerm["aggs"] = tail.AggTerm
			return term
		} else {
			term.AggTerm[term.FilterName].(map[string]interface{})["aggs"].(map[string]interface{})[term.FieldName].(map[string]interface{})["aggs"] = tail.AggTerm
			return term
		}

	} else {
		if term.Path != tail.Path { // path 不同，需要先出来，才能再进入
			if len(term.FilterName) == 0 {
				if len(term.NestedName) == 0 {
					panic("nestedName is None, but Path exists")
				}
				term.AggTerm[term.NestedName].(map[string]interface{})["aggs"].(map[string]interface{})[term.FieldName].(map[string]interface{})["aggs"] = map[string]interface{}{
					"out_nested": map[string]interface{}{
						"reverse_nested": map[string]interface{}{},
						"aggs":           tail.AggTerm,
					},
				}
				return term
			} else {
				term.AggTerm[term.FilterName].(map[string]interface{})["aggs"].(map[string]interface{})[term.NestedName].(map[string]interface{})["aggs"].(map[string]interface{})[term.FieldName].(map[string]interface{})["aggs"] = map[string]interface{}{
					"out_nested": map[string]interface{}{
						"reverse_nested": map[string]interface{}{},
						"aggs":           tail.AggTerm,
					},
				}
				return term
			}
		} else { // path 相同，无需重复进入
			if len(tail.FilterName) == 0 {
				var innerMap map[string]interface{}
				innerMap = tail.AggTerm[tail.NestedName].(map[string]interface{})["aggs"].(map[string]interface{})
				if len(term.FilterName) == 0 {
					term.AggTerm[term.NestedName].(map[string]interface{})["aggs"].(map[string]interface{})[term.FieldName].(map[string]interface{})["aggs"] = innerMap
					return term
				} else {
					term.AggTerm[term.FilterName].(map[string]interface{})["aggs"].(map[string]interface{})[term.NestedName].(map[string]interface{})["aggs"].(map[string]interface{})[term.FieldName].(map[string]interface{})["aggs"] = innerMap
					return term
				}
			} else {
				if len(term.FilterName) == 0 {
					term.AggTerm[term.NestedName].(map[string]interface{})["aggs"].(map[string]interface{})[term.FieldName].(map[string]interface{})["aggs"] = map[string]interface{}{
						"out_nested": map[string]interface{}{
							"reverse_nested": map[string]interface{}{},
							"aggs":           tail.AggTerm,
						},
					}
					return term
				} else {
					term.AggTerm[term.FilterName].(map[string]interface{})["aggs"].(map[string]interface{})[term.NestedName].(map[string]interface{})["aggs"].(map[string]interface{})[term.FieldName].(map[string]interface{})["aggs"] = map[string]interface{}{
						"out_nested": map[string]interface{}{
							"reverse_nested": map[string]interface{}{},
							"aggs":           tail.AggTerm,
						},
					}
					return term
				}
			}

		}
	}
}

func CombineTerms(terms []AggRes) AggRes {
	N := len(terms)
	var i, j int
	var result = terms[N-1]

	for i = N - 2; i >= 0; i-- { // i 前面
		j = i + 1 // j 后面
		if terms[j].IsSon {
			result = InClusion(terms[i], result)
		} else {
			result = Parallel(terms[i], result)
		}
	}
	// return result.AggTerm
	return result
}
