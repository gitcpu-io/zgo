package zgoutils

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

type bytesBox struct {
	sort []string
	box  *bytes.Buffer
}

func newBytesBoxForSlice() *bytesBox {
	return &bytesBox{
		sort: nil,
		box:  new(bytes.Buffer),
	}
}

func newBytesBox(mm map[string]interface{}) *bytesBox {
	keys := make([]string, len(mm))
	var i int

	for k, _ := range mm {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	return &bytesBox{
		sort: keys,
		box:  new(bytes.Buffer),
	}
}

func MarshalMap(mm map[string]interface{}) (string, error) {
	return mapEncoder(mm)
}

func MarshalSlice(ms []interface{}) (string, error) {
	return sliceEncoder(ms)
}

func sliceEncoder(se []interface{}) (string, error) {
	bb := newBytesBoxForSlice()
	//var k string
	//var v interface{}
	var e error
	var js string
	bb.write("[")
	for i, v := range se {
		switch v.(type) {
		case bool:
			boolEncoder(bb, v.(bool))
		case string:
			stringEncoder(bb, v.(string))
		case int8:
			intEncoder(bb, int(v.(int8)))
		case int32:
			intEncoder(bb, int(v.(int32)))
		case int64:
			intEncoder(bb, int(v.(int64)))
		case int:
			intEncoder(bb, v.(int))
		case float32:
			floatEncoder(bb, float64(v.(float32)))
		case float64:
			floatEncoder(bb, v.(float64))
		case []string:
			sliceStringEncoder(bb, v.([]string))
		case []int:
			sliceIntEncoder(bb, v.([]int))
		case []byte:
			stringEncoder(bb, string(v.([]byte)))
		case []uint:
			sliceUintEncoder(bb, v.([]uint))
		case []float64:
			sliceFloat64Encoder(bb, v.([]float64))
		case []float32:
			sliceFloat32Encoder(bb, v.([]float32))

		case map[string]interface{}:
			js, e = mapEncoder(v.(map[string]interface{}))

			if e == nil {
				bb.write(js)
			} else {
				return "", e
			}
		case []interface{}:
			js, e = sliceEncoder(v.([]interface{}))
			if e == nil {
				bb.write(js)
			} else {
				return "", e
			}

		default:
			if v == nil {
				bb.write("null")
			} else {
				return "", fmt.Errorf("Json not support this type %v", reflect.TypeOf(v))
			}
		}
		if i != len(se)-1 {
			bb.write(",")
		}
	}
	bb.write("]")
	return bb.String(), nil
}

func mapEncoder(mm map[string]interface{}) (string, error) {
	bb := newBytesBox(mm)
	var k string
	var v interface{}
	var e error
	var js string

	bb.write("{")
	for i := 0; i < len(bb.sort); i++ {
		k = bb.sort[i]
		v = mm[k]
		if i > 0 {
			bb.write(",")
		}
		bb.writeKey(k)
		switch v.(type) {
		case bool:
			boolEncoder(bb, v.(bool))
		case string:
			stringEncoder(bb, v.(string))
		case int8:
			intEncoder(bb, int(v.(int8)))
		case int32:
			intEncoder(bb, int(v.(int32)))
		case int64:
			intEncoder(bb, int(v.(int64)))
		case int:
			intEncoder(bb, v.(int))
		case float32:
			floatEncoder(bb, float64(v.(float32)))
		case float64:
			floatEncoder(bb, v.(float64))
		case []string:
			sliceStringEncoder(bb, v.([]string))
		case []int:
			sliceIntEncoder(bb, v.([]int))
		case []byte:
			stringEncoder(bb, string(v.([]byte)))
		case []uint:
			sliceUintEncoder(bb, v.([]uint))
		case []float64:
			sliceFloat64Encoder(bb, v.([]float64))
		case []float32:
			sliceFloat32Encoder(bb, v.([]float32))

		case map[string]interface{}:
			js, e = mapEncoder(v.(map[string]interface{}))

			if e == nil {
				bb.write(js)
			} else {
				return "", e
			}

		case []interface{}:
			js, e = sliceEncoder(v.([]interface{}))
			if e == nil {
				bb.write(js)
			} else {
				return "", e
			}

		default:
			if v == nil {
				bb.write("null")
			} else {
				return "", fmt.Errorf("Json not support this type %v", reflect.TypeOf(v))
			}
		}
	}
	bb.write("}")
	return bb.String(), nil
}

func (bb *bytesBox) write(s string) {
	bb.box.WriteString(s)
}

func (bb *bytesBox) writeKey(key string) {
	bb.write("\"" + key + "\":")
}

func (bb *bytesBox) String() string {
	return bb.box.String()
}

func stringEncoder(bb *bytesBox, value string) {
	bb.write("\"" + value + "\"")
}

func intEncoder(bb *bytesBox, value int) {
	vs := strconv.Itoa(value)
	bb.write(vs)
}

func floatEncoder(bb *bytesBox, value float64) {
	vs := strconv.FormatFloat(value, 'f', -1, 64)
	bb.write(vs)
}

func sliceStringEncoder(bb *bytesBox, value []string) {
	bb.write("[")
	var i, n = 0, len(value)
	for i = 0; i < n; i++ {
		if i > 0 {
			bb.write(",")
		}
		stringEncoder(bb, value[i])
	}
	bb.write("]")
}

func sliceIntEncoder(bb *bytesBox, value []int) {
	bb.write("[")
	var i, n = 0, len(value)
	for i = 0; i < n; i++ {
		if i > 0 {
			bb.write(",")
		}
		intEncoder(bb, value[i])
	}
	bb.write("]")
}

func sliceUintEncoder(bb *bytesBox, value []uint) {
	bb.write("[")
	var i, n = 0, len(value)
	for i = 0; i < n; i++ {
		if i > 0 {
			bb.write(",")
		}
		intEncoder(bb, int(value[i]))
	}
	bb.write("]")
}

func sliceFloat64Encoder(bb *bytesBox, value []float64) {
	bb.write("[")
	var i, n = 0, len(value)
	for i = 0; i < n; i++ {
		if i > 0 {
			bb.write(",")
		}
		floatEncoder(bb, value[i])
	}
	bb.write("]")
}

func sliceFloat32Encoder(bb *bytesBox, value []float32) {
	bb.write("[")
	var i, n = 0, len(value)
	for i = 0; i < n; i++ {
		if i > 0 {
			bb.write(",")
		}
		floatEncoder(bb, float64(value[i]))
	}
	bb.write("]")
}

func boolEncoder(bb *bytesBox, value bool) {
	if value {
		bb.box.WriteString("true")
	} else {
		bb.box.WriteString("false")
	}
}
