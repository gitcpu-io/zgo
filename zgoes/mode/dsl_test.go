package mode

import (
	"fmt"
	"testing"
)

func TestRang(t *testing.T) {
	//ps := map[string]interface{}{"type_name": "小区"}
	rg := map[string]interface{}{"gte": 10, "lte": 20}

	a := fmt.Sprintf(DSL["bool"], Range("city", rg), "[]", "[]")

	fmt.Println(a)
}

func TestMatch(t *testing.T) {
	mch := map[string]interface{}{"cityarea2": "朝阳"}
	a := fmt.Sprintf(DSL["bool"], Match(mch), "[]", "[]")

	fmt.Println(a, string(Match(mch).([]byte)))
}

func TestMulti(t *testing.T) {
	rg := map[string]interface{}{"gte": 10, "lte": 20}
	mch := map[string]interface{}{"cityarea2": "朝阳"}

}
