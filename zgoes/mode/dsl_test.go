package mode

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoes"
	"reflect"
	"testing"
)

func show(T interface{}) {
	fmt.Println(T, reflect.TypeOf(T))
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

	showQuery(BoolMust(must), "must")
	showQuery(BoolShould(should), "should")
	showQuery(BoolMustNot(must_not), "must_not")
	showQuery(BoolFilter(filter), "filter")

	args["must"] = must
	res = QueryDsl(args)
	show(res)
}

func TestDslStruct_1(t *testing.T) {
	fmt.Println("#########################QueryDsl")

	dsl := zgoes.NewDSL()
	var res interface{}
	res = TermField("cityarea_id", 5)
	show(res)
	dsl.Must(res)

	res = MatchPhraseField("cityarea_name", "朝阳")
	dsl.Should(res)

	res = TermField("borough_id", 10)
	dsl.MustNot(res)

	res = RangeField("cityarea_id", "lt", 200, "gte", 10)
	show(res)
	dsl.Filter(res)

	aggs := SimpleAggs("cityarea_id", 5)
	sort := SimpleSort("borough_id", true)

	dsl.SetAggs(aggs)
	dsl.SetSort(sort)
	dsl.SetFrom(10)
	dsl.SetSize(100)

	dsl.Set_Source([]string{"_id", "cityarea_id"})

	dslstr := dsl.QueryDsl()
	show(dsl)
	show(dslstr)
}

func TestDslStruct_2(t *testing.T) {
	fmt.Println("#########################QueryDsl")

	dsl := zgoes.NewDSL()

	m1 := dsl.TermField("cityarea_id", 5)
	show(m1)
	//dsl.Must(m1)

	s1 := dsl.MatchPhraseField("cityarea_name", "朝阳")
	//dsl.Should(s1)

	mn1 := dsl.TermField("borough_id", 10)
	//dsl.MustNot(mn1)

	f1 := dsl.RangeField("cityarea_id", "lt", 200, "gte", 10)
	//dsl.Filter(f1)

	aggs := dsl.SimpleAggs("cityarea_id", 5)
	sort := dsl.SimpleSort("borough_id", true)

	//dsl.SetAggs(aggs)
	//dsl.SetSort(sort)
	//dsl.SetFrom(10)
	//dsl.SetSize(100)

	//dsl.Set_Source([]string{"_id", "cityarea_id"})
	//dsl.Set_SourceField("_id", "cityarea_id")

	//dsl.Must(m1).Should(s1).MustNot(mn1).Filter(f1).SetAggs(aggs).SetSort(sort).SetFrom(10).SetSize(100).Set_SourceField("_id", "cityarea_id")
	dsl.Must(m1).Should(s1).MustNot(mn1).Filter(f1)
	dsl.SetAggs(aggs).SetSort(sort).SetFrom(10).SetSize(100).Set_SourceField("_id", "cityarea_id")

	dslstr := dsl.QueryDsl()
	show(dsl)
	show(dslstr)

	binterface := dsl.BoolDslTerm()
	show(binterface)
	bstring := dsl.BoolDslString()
	show(bstring)
}

func TestUpEs(t *testing.T) {

}
