package zgoutils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var u = New()

var st = &struct {
	A string
}{
	A: "niubi",
}

func TestUtils_IPs(t *testing.T) {
	fmt.Println(u.IPs())

	fmt.Println(u.IsPrivateIP("192.168.100.162"))
	fmt.Println(u.IsPrivateIP("121.69.135.49"))

	e, err := u.Extract("ba")
	if err != nil {
		panic(err)
	}
	fmt.Println(e)
	fmt.Println(u.GetIntranetIP())

	fmt.Println(u.NowUnix(), u.GetTimestamp(13))
}

func TestUtils_ServiceMethod(t *testing.T) {
	a, b, c := u.GrpcServiceMethodConverts("Foo.Bar")
	fmt.Println(a, b, c)
}

func TestUtils_StructToMap(t *testing.T) {
	s := &struct {
		A string
	}{
		A: "niubi",
	}
	m := u.StructToMap(s)
	fmt.Println(m["A"])
}

func TestUtils_Marshal(t *testing.T) {
	str, err := u.Marshal(st)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(str), err)
}
func BenchmarkUtils_Marshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := u.Marshal(st)
		if err != nil {
			panic(err)
		}
	}
}

func TestUtils_NewDecoder(t *testing.T) {
	f, _ := os.Open("../config/local.json")
	defer f.Close()

	str := u.NewDecoder(f)

	fmt.Print(str)
}

func TestUtils_Unmarshal(t *testing.T) {
	st := &struct {
		B string
	}{}
	var str = `{"B":"456"}`
	u.Unmarshal([]byte(str), st)
	fmt.Println(st)
	fmt.Println(str)
}

func TestInitStructWithDefaultTag(t *testing.T) {
	config := &struct {
		B bool          `default:"true"`
		S string        `default:"default string"`
		I int           `default:"10"`
		T time.Duration `default:"100"`
		E int           `default:""`
	}{}
	u.InitStructWithDefaultTag(config)
	assert.NotNil(t, config)
	assert.Equal(t, true, config.B)
	assert.Equal(t, "default string", config.S)
	assert.Equal(t, 10, config.I)
	assert.Equal(t, time.Duration(100), config.T)
	assert.Equal(t, 0, config.E)
}

func TestGetUUIDV4(t *testing.T) {
	uuid := u.GetUUIDV4()
	assert.Equal(t, 32, len(uuid))
	assert.NotEqual(t, u.GetUUIDV4(), u.GetUUIDV4())
}

func TestGetMD5Base64(t *testing.T) {
	assert.Equal(t, "ERIHLmRX2uZmssDdxQnnxQ==",
		u.GetMD5Base64([]byte("That's all folks!!")))
	assert.Equal(t, "GsJRdI3kAbAnHo/0+3wWJw==",
		u.GetMD5Base64([]byte("中文也没啥问题")))
}

func TestGetTimeInFormatRFC2616(t *testing.T) {
	s := u.GetTimeInFormatRFC2616()
	assert.Equal(t, 29, len(s))
	re := regexp.MustCompile(`^[A-Z][a-z]{2}, [0-9]{2} [A-Z][a-z]{2} [0-9]{4} [0-9]{2}:[0-9]{2}:[0-9]{2} GMT$`)
	assert.True(t, re.MatchString(s))
}

func TestGetTimeInFormatISO8601(t *testing.T) {
	s := u.GetTimeInFormatISO8601()
	assert.Equal(t, 20, len(s))
	// 2006-01-02T15:04:05Z
	re := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$`)
	assert.True(t, re.MatchString(s))
}

func TestGetUrlFormedMap(t *testing.T) {
	m := make(map[string]string)
	m["key"] = "value"
	s := u.GetUrlFormedMap(m)
	assert.Equal(t, "key=value", s)
	m["key2"] = "http://domain/?key=value&key2=value2"
	s2 := u.GetUrlFormedMap(m)
	assert.Equal(t, "key=value&key2=http%3A%2F%2Fdomain%2F%3Fkey%3Dvalue%26key2%3Dvalue2", s2)
}

func TestGetTimeInFormatISO8601WithTZData(t *testing.T) {
	tZData = []byte(`"GMT"`)
	loadLocationFromTZData = func(name string, data []byte) (location *time.Location, e error) {
		if strings.Contains(string(data), name) {
			location, _ = time.LoadLocation(name)
		}
		e = fmt.Errorf("There is a error in test.")
		return location, e
	}
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	s := u.GetTimeInFormatISO8601()
	assert.Equal(t, 20, len(s))
	// 2006-01-02T15:04:05Z
	re := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$`)
	assert.True(t, re.MatchString(s))
}

func TestGetTimeInFormatRFC2616WithTZData(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	s := u.GetTimeInFormatRFC2616()
	assert.Equal(t, 29, len(s))
	re := regexp.MustCompile(`^[A-Z][a-z]{2}, [0-9]{2} [A-Z][a-z]{2} [0-9]{4} [0-9]{2}:[0-9]{2}:[0-9]{2} GMT$`)
	assert.True(t, re.MatchString(s))
}
func TestJosn(t *testing.T) {
	son := "{}";
	sss, _ := u.Marshal(son)
	fmt.Println(string(sss))
}
