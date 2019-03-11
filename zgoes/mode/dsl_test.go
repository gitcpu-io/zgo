package mode

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	//ps := map[string]interface{}{"type_name": "小区"}
	rg := map[string]interface{}{"gte": 10, "lte": 20}

	a := fmt.Sprintf(DSL["bool"], Rang("city", rg), "[]", "[]")

	fmt.Println(a)
}
