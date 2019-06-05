package parseES

import (
	"fmt"
	"reflect"
	"strings"
)

type Bool struct {
	Must    []map[string]map[string]interface{} `json:"must"`
	MustNot []map[string]map[string]interface{} `json:"must_not"`
	Should  []map[string]map[string]interface{} `json:"should"`
	Minimum int                                 `json:"minimum_should_match"`
}

func NewBool(minimum int) *Bool {
	return &Bool{Minimum: minimum}
}

var statusTemp map[string]map[string]interface{}

var termMap = map[string]string{}

var matchMap = map[string]string{}

var rangeMap = map[string]string{}

var wildCard = map[string]string{}

var nestSort = map[string]string{
	"public_time":        "house_info",
	"gov_id":             "house_info",
	"created":            "house_info",
	"house_price":        "house_info",
	"verify_status":      "house_info",
	"sid":                "house_subway",
	"source_type":        "house_info",
	"house_floor":        "house_info",
	"source":             "house_info",
	"line_name":          "house_subway",
	"verify_reason":      "house_info",
	"house_built_year":   "house_info",
	"verify_time":        "house_info",
	"house_fitment":      "house_info",
	"house_toward":       "house_info",
	"source_url":         "house_info",
	"owner_phone":        "house_info",
	"house_desc":         "house_info",
	"bid":                "house_subway",
	"rang":               "house_subway",
	"user_id":            "house_info",
	"broker_num":         "house_info",
	"owner_name":         "house_info",
	"station_name":       "house_subway",
	"price_change":       "house_info",
	"cityarea2_id":       "cityarea2",
	"price_change_value": "house_info",
	"lid":                "house_subway",
	"status":             "house_info",
	"cityarea2_py":       "cityarea2",
	"service_phone":      "house_info",
	"updated":            "house_info",
	"price_updated":      "house_info",
	"dial_time":          "house_info",
	"tag":                "house_info",
	"not_tag":            "house_info",
	"not_source_name":    "house_info",
	"source_name":        "house_info",
	"cityarea2_name":     "cityarea2",
	"not_cityarea2_name": "cityarea2",
}

func Init(match, term, wild, ranges map[string]string) {
	matchMap, termMap, wildCard, rangeMap = match, term, wild, ranges
}

func (ret *Bool) Parse(parseMap map[string]interface{}) *Bool {
	ret.JustSpecial(parseMap)
	ret.CheckExistInMap(parseMap, termMap, "term")
	ret.CheckExistInMap(parseMap, rangeMap, "range")
	ret.CheckExistInMap(parseMap, matchMap, "match")
	ret.CheckExistInMap(parseMap, wildCard, "wildcard")
	ret.CheckShouldInMap(parseMap)
	statusTemp = nil
	return ret
}

func (ret *Bool) JustSpecial(parseMap map[string]interface{}) *Bool {
	if parseMap["status"] != nil {
		statusTemp = map[string]map[string]interface{}{"term": {"house_info.status": parseMap["status"]}}
		tempBool := NewBool(0)
		tempBool.Must = append(tempBool.Must, statusTemp)
		ret.Must = append(ret.Must, map[string]map[string]interface{}{"nested": {"path": "house_info", "query": map[string]interface{}{"bool": tempBool}}})
	}

	if parseMap["price_change"] != nil && parseMap["price_updated"] != nil {
		tempBool := NewBool(0)
		tempBool.Must = append(tempBool.Must, map[string]map[string]interface{}{"term": {"house_info.price_change": parseMap["price_change"]}})
		array := strings.Split(parseMap["price_updated"].(string), "-")
		temp := map[string]string{"gte": array[0], "lte": array[1]}
		tempBool.Must = append(tempBool.Must, map[string]map[string]interface{}{"range": {"house_info.price_updated": temp}})
		if statusTemp != nil {
			tempBool.Must = append(tempBool.Must, statusTemp)
		}
		ret.Must = append(ret.Must, map[string]map[string]interface{}{"nested": {"path": "house_info", "query": map[string]interface{}{"bool": tempBool}}})
	}

	if parseMap["tag"] != nil {
		ret.JudgeTag(parseMap["tag"], true)
	}

	if parseMap["not_tag"] != nil {
		ret.JudgeTag(parseMap["not_tag"], false)
	}

	return ret
}

func (ret *Bool) JudgeTag(value interface{}, flag bool) {
	var sValue string
	switch value.(type) {
	case string:
		sValue = value.(string)
	default:
		sValue = fmt.Sprint(value)
	}
	if strings.Contains(sValue, ",") {
		for _, tag := range strings.Split(sValue, ",") {
			tempBool := NewBool(0)
			tempBool.Must = append(tempBool.Must, map[string]map[string]interface{}{"term": {"house_info.tag": tag}})
			if flag {
				if statusTemp != nil {
					tempBool.Must = append(tempBool.Must, statusTemp)
				}
				ret.Must = append(ret.Must, map[string]map[string]interface{}{"nested": {"path": "house_info", "query": map[string]interface{}{"bool": tempBool}}})
			} else {
				ret.MustNot = append(ret.MustNot, map[string]map[string]interface{}{"nested": {"path": "house_info", "query": map[string]interface{}{"bool": tempBool}}})
			}
		}
	} else {
		tempBool := NewBool(0)
		tempBool.Must = append(tempBool.Must, map[string]map[string]interface{}{"term": {"house_info.tag": sValue}})
		if flag {
			if statusTemp != nil {
				tempBool.Must = append(tempBool.Must, statusTemp)
			}
			ret.Must = append(ret.Must, map[string]map[string]interface{}{"nested": {"path": "house_info", "query": map[string]interface{}{"bool": tempBool}}})
		} else {
			ret.MustNot = append(ret.MustNot, map[string]map[string]interface{}{"nested": {"path": "house_info", "query": map[string]interface{}{"bool": tempBool}}})
		}
	}
}

func (ret *Bool) CheckExistInMap(parseMap map[string]interface{}, destMap map[string]string, name string) *Bool {
	for key, value := range parseMap {
		var sValue string
		if mapValue, ok := destMap[key]; ok {
			switch value.(type) {
			case string:
				sValue = value.(string)
			default:
				sValue = fmt.Sprint(value)
			}
			ret.CreateBool(mapValue, name, key, sValue)
		}
	}
	return ret
}

func (ret *Bool) CheckShouldInMap(parseMap map[string]interface{}) *Bool {
	if value, ok := parseMap["should"]; ok {
		for _, shouldItem := range value.([]interface{}) {
			//RecursionShould(shouldItem, ret)
			resultTemp := make(map[string]interface{}, 0)
			itemBool := NewBool(0)
			for _, shouldMap := range shouldItem.([]interface{}) {
				tempBool := NewBool(0)
				tempBool.JustSpecial(shouldMap.(map[string]interface{}))
				tempBool.CheckExistInMap(shouldMap.(map[string]interface{}), termMap, "term")
				tempBool.CheckExistInMap(shouldMap.(map[string]interface{}), rangeMap, "range")
				tempBool.CheckExistInMap(shouldMap.(map[string]interface{}), matchMap, "match")
				tempBool.CheckExistInMap(shouldMap.(map[string]interface{}), wildCard, "wildcard")
				if tempBool.Should != nil {
					itemBool.Should = append(itemBool.Should, tempBool.Should...)
				}
				if tempBool.Must != nil {
					itemBool.Should = append(itemBool.Should, map[string]map[string]interface{}{"bool": {"must": tempBool.Must}})
				}
				if tempBool.MustNot != nil {
					itemBool.MustNot = append(itemBool.MustNot, tempBool.MustNot...)
				}
			}
			resultTemp["should"] = itemBool.Should
			resultTemp["must_not"] = itemBool.MustNot
			ret.Must = append(ret.Must, map[string]map[string]interface{}{"bool": resultTemp})
		}

	}
	return ret
}

func RecursionShould(shouldItem interface{}, ret *Bool) {
	resultTemp := make(map[string]interface{}, 0)
	itemBool := NewBool(0)
	if reflect.TypeOf(shouldItem) == reflect.TypeOf(map[string]interface{}{}) {
		tempBool := NewBool(0)
		tempBool.JustSpecial(shouldItem.(map[string]interface{}))
		tempBool.CheckExistInMap(shouldItem.(map[string]interface{}), termMap, "term")
		tempBool.CheckExistInMap(shouldItem.(map[string]interface{}), rangeMap, "range")
		tempBool.CheckExistInMap(shouldItem.(map[string]interface{}), matchMap, "match")
		tempBool.CheckExistInMap(shouldItem.(map[string]interface{}), wildCard, "wildcard")
		if tempBool.Should != nil {
			itemBool.Should = append(itemBool.Should, tempBool.Should...)
		} else if tempBool.Must != nil {
			itemBool.Should = append(itemBool.Should, map[string]map[string]interface{}{"bool": map[string]interface{}{"must": tempBool.Must}})
		} else if tempBool.MustNot != nil {
			itemBool.MustNot = append(itemBool.MustNot, tempBool.MustNot...)
		}
	} else if reflect.TypeOf(shouldItem) == reflect.TypeOf([]interface{}{}) {
		for _, value := range shouldItem.([]interface{}) {
			RecursionShould(value, ret)
		}
	}
	resultTemp["should"] = itemBool.Should
	resultTemp["must_not"] = itemBool.MustNot
	ret.Must = append(ret.Must, map[string]map[string]interface{}{"bool": resultTemp})
}

func (ret *Bool) CreateBool(mapValue string, name, key string, value interface{}) *Bool {
	key, value, name, switchFlag := DefineValueByName(value, key, name)
	if mapValue != "" {
		tempBool := NewBool(0)
		switch {
		case switchFlag == "must_not":
			if name == "range" {
				//特殊处理 range的mustnot情况
				ret.MustNot = append(ret.MustNot, map[string]map[string]interface{}{"nested": {"path": mapValue, "query": map[string]interface{}{"bool": value}}})
			} else {
				tempBool.Must = append(tempBool.Must, map[string]map[string]interface{}{name: {mapValue + "." + key: value}})
				ret.MustNot = append(ret.MustNot, map[string]map[string]interface{}{"nested": {"path": mapValue, "query": map[string]interface{}{"bool": tempBool}}})
			}
		case switchFlag == "must":
			tempBool.Must = append(tempBool.Must, map[string]map[string]interface{}{name: {mapValue + "." + key: value}})
			if statusTemp != nil && mapValue == "house_info" {
				tempBool.Must = append(tempBool.Must, statusTemp)
			}
			ret.Must = append(ret.Must, map[string]map[string]interface{}{"nested": {"path": mapValue, "query": map[string]interface{}{"bool": tempBool}}})
		case switchFlag == "should":
			if name == "range" {
				//tempBool.Should = append(tempBool.Should, map[string]map[string]interface{}{name: {mapValue + "." + key: value}})
				//防止range中嵌套should语句
				ret.Must = append(ret.Must, map[string]map[string]interface{}{"nested": {"path": mapValue, "query": map[string]interface{}{"bool": value}}})
			} else {
				tempBool.Should = append(tempBool.Should, map[string]map[string]interface{}{name: {mapValue + "." + key: value}})
				ret.Must = append(ret.Must, map[string]map[string]interface{}{"nested": {"path": mapValue, "query": map[string]interface{}{"bool": tempBool}}})
			}
		}
	} else {
		switch {
		case switchFlag == "must_not":
			if name == "range" {
				ret.MustNot = append(ret.MustNot, value.(*Bool).Should...)
			} else {
				ret.MustNot = append(ret.MustNot, map[string]map[string]interface{}{name: {key: value}})
			}
		case switchFlag == "must":
			ret.Must = append(ret.Must, map[string]map[string]interface{}{name: {key: value}})
		case switchFlag == "should":
			ret.Must = append(ret.Must, map[string]map[string]interface{}{"bool": {"should": value.(*Bool).Should}})
			//ret.Should = append(ret.Should, value.(*Bool).Should...)
		}
	}
	return ret
}

func DefineValueByName(value interface{}, key, name string) (string, interface{}, string, string) {
	var switchFlag = "must"
	if strings.HasPrefix(key, "not_") {
		switchFlag = "must_not"
		key = key[4:]
	}
	if name == "term" {
		if strings.Contains(value.(string), ",") {
			valueList := strings.Split(value.(string), ",")
			return key, valueList, "terms", switchFlag
		} else {
			return key, value, "term", switchFlag
		}
	} else if name == "wildcard" {
		return key, "*" + value.(string) + "*", name, switchFlag
	} else if name == "match" {
		return key, map[string]interface{}{"query": value.(string), "slop": 10}, "match_phrase", switchFlag
	} else if name == "range" {
		array := make([]string, 1)
		should := NewBool(1)
		if switchFlag == "must" {
			switchFlag = "should"
		}
		rangeValue := rangeMap[key]
		if rangeValue != "" {
			key = rangeValue + "." + key
		}

		if strings.Contains(value.(string), ",") {
			array = strings.Split(value.(string), ",")
			for _, rangeValue := range array {
				if strings.Contains(rangeValue, "-") {
					array = strings.Split(rangeValue, "-")
					should.Should = append(should.Should, map[string]map[string]interface{}{"range": {key: map[string]interface{}{"gte": array[0], "lte": array[1]}}})
				} else {
					should.Should = append(should.Should, map[string]map[string]interface{}{"term": {key: rangeValue}})
				}
			}
			if statusTemp != nil && rangeValue == "house_info" {
				should.Must = append(should.Must, statusTemp)
			}
			return key, should, "range", switchFlag
		} else if strings.Contains(value.(string), "-") {
			array = strings.Split(value.(string), "-")
			should.Should = append(should.Should, map[string]map[string]interface{}{"range": {key: map[string]interface{}{"gte": array[0], "lte": array[1]}}})
			if statusTemp != nil && rangeValue == "house_info" {
				should.Must = append(should.Must, statusTemp)
			}
			return key, should, "range", switchFlag
		} else {
			should.Should = append(should.Should, map[string]map[string]interface{}{"term": {key: value}})
			if statusTemp != nil && rangeValue == "house_info" {
				should.Must = append(should.Must, statusTemp)
			}
			return key, should, "range", switchFlag
		}

	}
	return key, value, name, switchFlag
}

func ParseSort(sortMap []Sort) []map[string]interface{} {
	sortRet := make([]map[string]interface{}, 0)
	for _, value := range sortMap {
		term := make(map[string]interface{})
		if mapValue, ok := nestSort[value.Field]; ok {
			term[mapValue+"."+value.Field] = map[string]string{"order": value.Type, "nested_path": value.Field}
			sortRet = append(sortRet, term)
		} else {
			term[value.Field] = map[string]string{"order": value.Type}
			sortRet = append(sortRet, term)
		}
	}
	return sortRet
}
