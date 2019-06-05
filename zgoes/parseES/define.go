package parseES

import (
	"fmt"
)

const smallsString = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

const digits = "0123456789"

type ESStruct struct {
	Filter map[string]interface{} `json:"filter"`
	From   interface{}            `json:"from"`
	Size   interface{}            `json:"size"`
	Sort   []Sort                 `json:"sort"`
	Data   interface{}            `json:"data"`
	Page   interface{}            `json:"page"`
	Source []interface{}          `json:"_source"`
	Aggs   []AggCell              `json:"aggs"`
	Aggs2  [][]AggCell            `json:"aggs2"`
}

type Sort struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

var actions = [...]string{"filter", "terms", "top_hits", "avg", "min", "max", "sum", "ranges", "filter-terms", "filter-ranges", "filter-top_hits", "filter-tophits", "filter-avg", "filter-min", "filter-max", "filter-sum"}

//var actionAll = [...]string{"filter", "terms", "top_hits", "avg", "min", "max", "sum", "ranges", "filter-terms", "filter-ranges", "filter-top_hits", "filter-tophits", "filter-avg", "filter-min", "filter-max", "filter-sum"}

func smallItoa(n int) string {
	if n < 10 {
		return digits[n : n+1]
	}
	return smallsString[n*2 : n*2+2]
}

type AggCell struct {
	Name   string      `json:"name"`
	Path   string      `json:"path"`
	Action string      `json:"action"` // range filter terms top_hits metric
	Field  string      `json:"field"`
	Sort   string      `json:"sort"`
	Size   int         `json:"size"`
	Ranges [][2]int    `json:"ranges"`
	Query  interface{} `json:"query"`
	IsSon  bool        `json:"is_son"` // 与其前一个cell关系，并列或者嵌套
}

type EsParams struct {
	Aggs []AggCell `json:"aggs"`
}

type AggRes struct {
	FieldName  string // field Name
	FilterName string
	NestedName string
	Path       string
	AggTerm    map[string]interface{}
	IsSon      bool
	// IsBucket   bool
	//IsNested bool
}

func (ac *AggCell) AggName() string {
	if ac.Name != "" {
		return ac.Name
	}
	var field string
	if ac.Field == "" {
		field = "field"
	} else {
		field = ac.Field
	}
	return ac.Action + "_" + field
}

func (ac *AggCell) IsNested() bool {
	if ac.Path == "" {
		return false
	}
	return true
}

func (ac *AggCell) NestedField() string {
	return ac.Path + "." + ac.Field
}

func (ac *AggCell) NestedName() string {
	return "nested_" + ac.Field
}

func (ac *AggCell) AggSize() int {
	if ac.Size != 0 {
		return ac.Size
	}
	return 10
}

func (ac *AggCell) AggAction() (string, error) {
	var res string

	action := func(s string) string {
		for i := 0; i < len(actions); i++ {
			if s == actions[i] {
				return s
			}
		}
		return ""
	}
	res = action(ac.Action)
	if res != "" {
		return res, nil
	}

	return "", fmt.Errorf("action Error, you gave %v, want one of %v", ac.Action, actions)
}

func (ac *AggCell) AggRanges() []map[string]interface{} {
	//fmt.Println(ac.Ranges)
	N := len(ac.Ranges)
	res := make([]map[string]interface{}, 0)

	for i := 0; i < N-1; i++ {
		res = append(res, map[string]interface{}{
			"key":  "range_" + smallItoa(i),
			"from": ac.Ranges[i][0],
			"to":   ac.Ranges[i][1],
		})
	}

	if ac.Ranges[N-1][1] == 0 {
		res = append(res, map[string]interface{}{
			"key":  "range_" + smallItoa(N-1),
			"from": ac.Ranges[N-1][0],
		})
	} else {
		res = append(res, map[string]interface{}{
			"key":  "range_" + smallItoa(N-1),
			"from": ac.Ranges[N-1][0],
			"to":   ac.Ranges[N-1][1],
		})
	}
	return res
}
