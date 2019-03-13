package mode

import (
	"fmt"
	"reflect"
	"testing"
)

func show(T interface{}) {
	fmt.Println(T, reflect.TypeOf(T))
}

func TestRang(t *testing.T) {
	//ps := map[string]interface{}{"type_name": "小区"}
	rg := map[string]interface{}{"gte": 10, "lte": 20}

	a := fmt.Sprintf(DSL["bool"], unmarshal(RangeMap("city", rg)), "[]", "[]")

	fmt.Println(a)
}

func TestMatch(t *testing.T) {
	mch := map[string]interface{}{"cityarea2": "朝阳"}
	a := fmt.Sprintf(DSL["bool"], unmarshal(MatchMap(mch)), "[]", "[]")

	fmt.Println(a, MatchMap(mch))
}

func TestSimpleAggs(t *testing.T) {
	res := SimpleAggs("cityarea_id", 20)
	show(res)
}

func TestAggs2Field(t *testing.T) {
	res := Aggs2Field("cityarea_id", 5, "borough_id", 10)
	show(fmt.Sprintf(`{"aggs": %s}`, res))
}

func TestQueryDsl(t *testing.T) {
	fmt.Println("#########################QueryDsl")
	must := make([]interface{}, 0)
	should := make([]interface{}, 0)
	filter := make([]interface{}, 0)
	must_not := make([]interface{}, 0)

	var res interface{}
	res = TermField("cityarea_id", 5)
	show(res)
	must = append(must, res)
	res = MatchPhraseField("cityarea_name", "朝阳")
	should = append(should, res)

	res = TermField("borough_id", 10)
	must_not = append(must_not, res)

	res = RangeField("cityarea_id", "lt", 200, "gte", 10)
	show(res)
	filter = append(filter, res)

	aggs := SimpleAggs("cityarea_id", 5)
	sort := SimpleSort("borough_id", true)

	resMap := make(map[string]interface{})
	resMap["must"] = must
	resMap["should"] = should
	resMap["must_not"] = must_not
	resMap["filter"] = filter

	resMap["sort"] = sort
	resMap["aggs"] = aggs
	resMap["from"] = 10
	resMap["size"] = 100
	resMap["_source"] = []string{"_id", "cityarea_id"}

	dsl := QueryDsl(resMap)
	show(dsl)
}

func showQuery(T interface{}, flag string) {
	fmt.Println(flag)
	show(unmarshal(T))
}

func TestQueryBoolEmpty(t *testing.T) {
	println("############################################")
	args := make(map[string]interface{})
	aggs := SimpleAggs("cityarea_id", 5)
	args["aggs"] = aggs
	var res interface{}

	// res = MatchPhraseField("cityarea_name", "朝阳")

	xx := []int{1, 2, 3, 4, 5}
	res = TermsField("cityarea_id", xx)
	must := make([]interface{}, 0)
	must = append(must, res)

	should := make([]interface{}, 0)
	should = append(should, res)

	filter := make([]interface{}, 0)
	filter = append(filter, res)

	must_not := make([]interface{}, 0)
	must_not = append(must_not, res)

	showQuery(MustQuery(must), "must")
	showQuery(ShouldQuery(should), "should")
	showQuery(MustNotQuery(must_not), "must_not")
	showQuery(FilterQuery(filter), "filter")

	args["must"] = must
	res = QueryDsl(args)
	show(res)

}
