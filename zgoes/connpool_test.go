package zgoes

import "testing"

const (
	label_sell = "label_sell"
	label_rent = "label_rent"
)

func TestEsSearch(t *testing.T) {
	InitEs(map[string][]string{
		label_sell: []string{"localhost:27017"},
		label_rent: []string{"localhost:27017"},
	}) //测试时表示使用nsq，在zgo_start中使用一次

}
