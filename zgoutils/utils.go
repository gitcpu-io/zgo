package zgoutils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoutils/dateparse"
	"github.com/bwmarrin/snowflake"
	"github.com/fatih/structs"
	"github.com/json-iterator/go"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	ip_cache  = ""
	ip_byName = ""
)

var (
	privateBlocks []*net.IPNet
)

const SysTimeform = "2006-01-02 15:04:05"
const SysTimeformShort = "2006-01-02"
const TimeformYm = "200601"

// 中国时区
var SysTimeLocation *time.Location

var Utils Utilser

func init() {
	Utils = New()

	s, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		SysTimeLocation = time.FixedZone("CST", 8*3600)
	} else {
		SysTimeLocation = s
	}
	//time.Local = s

	for _, b := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10"} {
		if _, block, err := net.ParseCIDR(b); err == nil {
			privateBlocks = append(privateBlocks, block)
		}
	}
}

type Utilser interface {
	GBK2UTF8(s []byte) ([]byte, error)
	UTF82GBK(s []byte) ([]byte, error)
	ToString(data interface{}) (string, error)

	//雪花算法生成器
	Snowflake(i ...int) snowflake.ID
	//sha1算法
	SHA1(str string) string

	//是否是email地址
	IsEmail(email string) (ok bool, err error)
	//是否是中文
	IsChineseWords(words string) (ok bool, err error)
	//是否是身份证号
	IsChineseID(s string) (ok bool, err error)
	//是否是银行卡号
	IsBankCard(n int64) (ok bool, err error)

	//Marshal 序列化为json
	Marshal(in interface{}) ([]byte, error)
	//Unmarshal 反序列化为go 内存对象
	Unmarshal(bytes []byte, in interface{}) error

	NewDecoder(reader io.Reader) *jsoniter.Decoder
	NewEncoder(writer io.Writer) *jsoniter.Encoder

	// mongo bson marshal
	BsonMarshal(in interface{}) ([]byte, error)
	// mongo bson unmarshal
	BsonUnmarshal(bytes []byte, in interface{}) error

	//string转map[string]interface{}
	StringToMap(str string) map[string]interface{}
	//结构体转map[string]interface{}
	StructToMap(interface{}) map[string]interface{}
	//map转结构体
	MapToStruct(input interface{}, out interface{})
	// GrpcServiceMethod converts a gRPC method to a Go method
	GrpcServiceMethodConverts(m string) (string, string, error)

	ParseDns(strDns string) bool
	IPs() []string
	//是否是内网ip
	IsPrivateIP(ipAddr string) bool
	Extract(addr string) (string, error)
	//获取内网ip
	GetIntranetIP() string
	//通过网关名来获取内网或外网ip，name=eth0或eth1
	GetIntranetIPByName(name string) string
	//数字转ip地址
	UInt32ToIP(intIP uint32) net.IP
	//ip地址转数字
	IPToUInt32(ipnr net.IP) uint32

	GetUUIDV4() string
	GetMD5Base64([]byte) string
	ByteToBase64([]byte) string
	Base64ToByte(string) ([]byte, error)
	GetGMTLocation() (*time.Location, error)
	GetTimeInFormatISO8601() string
	GetTimeInFormatRFC2616() string
	//从一个map中返回a=123&b=456
	GetUrlFormedMap(map[string]string) string
	//把a=123&b=456转成map[a:123,b:456]
	ParseUrlEscapeToMap(params string) (map[string]string, error)
	InitStructWithDefaultTag(interface{})

	//获取当前时间时间戳，n= 10/13/19 位时间戳
	GetTimestamp(n int) int64
	NowUnix() int
	FormatFromUnixTime(t int64) string
	FormatFromUnixTimeShort(t int64) string
	FormatUnixTime(year int, month int, day int) string
	FormatUnixTimeShort(year int, month int, day int) string
	FormatUnixTimeYm(year int, month int, day int) string
	//转化14位字符时间为标准时间格式
	FormatStringToStandTimeString(str string) string
	//转化任意格式字符串为标准时间
	ParseTime(str string) (time.Time, error)

	IsYesToday(t int64) bool
	IsYesTodayByTime(t1, t2 int64) bool
	//生成最大数是max的一个随机数
	Random(max int) int
	//生成num个随机数0--9
	RandomInt64(num int) (int64, error)
	//生成num个随机数0--9,A--Z,a-z
	RandomString(num int) string
	RandRangeInt(min, max int) int
	RandRangeInt64(min, max int64) int64
	//获取任意一天0点的unix, 0表示当前，1表示明天，-1表示昨天，以此类推
	GetAnyDayZeroUnixTimestamp(t int) int64

	CreateSign(str string) string
	Addslashes(str string) string
	Stripslashes(str string) string
	Ip4toInt(ip string) int64
	NextDayDuration() time.Duration
	GetInt64(i interface{}, d int64) int64
	GetString(str interface{}, d string) string
	GetInt64FromMap(dm map[string]interface{}, key string, dft int64) int64
	GetInt64FromStringMap(dm map[string]string, key string, dft int64) int64
	GetStringFromMap(dm map[string]interface{}, key string, dft string) string
	GetStringFromStringMap(dm map[string]string, key string, dft string) string
	MarshalMap(mm map[string]interface{}) (string, error)
	RemoteIp(ctx iris.Context) string
}

var loadLocationFromTZData func(name string, data []byte) (*time.Location, error) = nil

var tZData []byte = nil

var jsonIterator = jsoniter.ConfigCompatibleWithStandardLibrary

type utils struct{}

func New() Utilser {
	return &utils{}
}

func (u *utils) Snowflake(i ...int) snowflake.ID {
	// Create a new Node with a Node number of 1
	p := 1
	if len(i) != 0 {
		p = i[0]
	}
	node, err := snowflake.NewNode(int64(p))
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// Generate a snowflake ID.
	return node.Generate()
}
func (u *utils) SHA1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// GBK2UTF8 transform s from GBK to UTF8 format
func (u *utils) GBK2UTF8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF82GBK transform s from UTF8 to GBK format
func (u *utils) UTF82GBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// ToString convert data to string, data must be one of map[string]string, map[string]interface{}, string, []string, struct
func (u *utils) ToString(data interface{}) (string, error) {
	if structs.IsStruct(data) {
		data = structs.Map(data)
	}
	var s string
	switch data.(type) {
	case string:
		s = data.(string)
	case []string:
		for _, v := range data.([]string) {
			s += fmt.Sprintf("%v\n", v)
		}
	case map[string]string:
		for k, v := range data.(map[string]string) {
			s += fmt.Sprintf("%v: %v\n", k, v)
		}
	case map[string]interface{}:
		for k, v := range data.(map[string]interface{}) {
			s += fmt.Sprintf("%v: %v\n", k, v)
		}
	default:
		return "", errors.New("Unsupport data")
	}
	return s, nil
}

// IsEmail determine whether it is email address
func (u *utils) IsEmail(email string) (ok bool, err error) {
	p := `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	ok, err = regexp.MatchString(p, email)
	return
}

// IsBankCard determine whether it is bankcard number
func (u *utils) IsBankCard(n int64) (ok bool, err error) {
	s := strconv.FormatInt(n, 10)
	var sum int
	var i int
	for i = 1; i < len(s); i++ {
		var now int
		now, _ = strconv.Atoi(string(s[len(s)-1-i]))
		if i%2 == 0 {
			sum += now
			continue
		}
		var _i int
		_i = now * 2
		sum += _i / 10
		sum += _i % 10
	}
	var v int
	v, _ = strconv.Atoi(string(s[len(s)-1]))
	if (sum+v)%10 == 0 {
		ok = true
	}
	return
}

// IsChineseID determine whether it is Chinese ID Card Number
func (u *utils) IsChineseID(s string) (ok bool, err error) {
	if len(s) != 18 {
		return
	}
	var sum int
	var i int
	for i = 1; i < len(s); i++ {
		var now int
		now, err = strconv.Atoi(string(s[len(s)-1-i]))
		if err != nil {
			return
		}
		var w int
		w = int(math.Pow(2, float64(i+1-1))) % 11
		sum += now * w
	}
	v := (12 - (sum % 11)) % 11
	if v == 10 {
		if strings.ToLower(string(s[len(s)-1])) != "x" {
			return
		}
		ok = true
		return
	}
	if string(s[len(s)-1]) != strconv.Itoa(v) {
		return
	}
	ok = true
	return
}

// IsChineseWords determine whether it is Chinese words
// Notice: NOT ALL
func (u *utils) IsChineseWords(words string) (ok bool, err error) {
	// every rune is chinese
	for _, c := range words {
		if !unicode.Is(unicode.Scripts["Han"], c) {
			return false, nil
		}
	}
	return true, nil
}

//Marshal 序列化为json
func (u *utils) Marshal(res interface{}) ([]byte, error) {
	return jsonIterator.Marshal(res)
}

//Unmarshal 反序列化为go 内存对象
func (u *utils) Unmarshal(message []byte, in interface{}) error {
	return jsonIterator.Unmarshal(message, in)
}

// bson Marshal 序列化为json
func (u *utils) BsonMarshal(res interface{}) ([]byte, error) {
	return bson.Marshal(res)
}

// bson Unmarshal 反序列化为go 内存对象
func (u *utils) BsonUnmarshal(bytes []byte, in interface{}) error {
	return bson.Unmarshal(bytes, in)
}

func (u *utils) NewDecoder(reader io.Reader) *jsoniter.Decoder {
	return jsoniter.NewDecoder(reader)
}

func (u *utils) NewEncoder(writer io.Writer) *jsoniter.Encoder {
	return jsoniter.NewEncoder(writer)
}

//ParseUrlEscapeToMap解析 a=123&b=456  到 map[a:123,b:456]
func (u *utils) ParseUrlEscapeToMap(params string) (map[string]string, error) {
	res, err := url.QueryUnescape(params)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	split := strings.Split(res, "&")
	for _, value := range split {
		sps := strings.Split(value, "=")
		m[sps[0]] = sps[1]
	}
	return m, nil
}

//GetUrlFormedMap 从一个map中返回a=123&b=456
func (u *utils) GetUrlFormedMap(source map[string]string) (urlEncoded string) {
	urlEncoder := url.Values{}
	for key, value := range source {
		urlEncoder.Add(key, value)
	}
	urlEncoded = urlEncoder.Encode()
	return
}

// StringToMap 字符串类型的json转成map
func (u *utils) StringToMap(str string) map[string]interface{} {
	var val interface{}
	err := jsonIterator.Unmarshal([]byte(str), &val)
	if err != nil {
		return nil
	}
	return val.(map[string]interface{})
}

//StructToMap 结构体转map[string]interface{}
func (u *utils) StructToMap(input interface{}) map[string]interface{} {
	var m map[string]interface{}
	b, _ := jsonIterator.Marshal(input)
	jsonIterator.Unmarshal(b, &m)
	return m
}

//MapToStruct map[string]interface{}转结构体
func (u *utils) MapToStruct(input interface{}, out interface{}) {
	b, _ := jsonIterator.Marshal(input)
	jsonIterator.Unmarshal(b, &out)
}

// GrpcServiceMethodConverts converts a gRPC method to a Go method
// Input:
// Foo.Bar, /Foo/Bar, /package.Foo/Bar, /a.package.Foo/Bar
// Output:
// [Foo, Bar]
func (u *utils) GrpcServiceMethodConverts(m string) (string, string, error) {
	if len(m) == 0 {
		return "", "", fmt.Errorf("malformed method name: %q", m)
	}

	// grpc method
	if m[0] == '/' {
		// [ , Foo, Bar]
		// [ , package.Foo, Bar]
		// [ , a.package.Foo, Bar]
		parts := strings.Split(m, "/")
		if len(parts) != 3 || len(parts[1]) == 0 || len(parts[2]) == 0 {
			return "", "", fmt.Errorf("malformed method name: %q", m)
		}
		service := strings.Split(parts[1], ".")
		return service[len(service)-1], parts[2], nil
	}

	// non grpc method
	parts := strings.Split(m, ".")

	// expect [Foo, Bar]
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed method name: %q", m)
	}

	return parts[0], parts[1], nil
}

//GetTimestamp
func (u *utils) GetTimestamp(f int) int64 {
	var result int64
	switch f {
	case 10:
		result = time.Now().Unix()
	case 13:
		result = time.Now().UnixNano() / 1e6
	case 19:
		result = time.Now().UnixNano()
	}
	return result
}

//IsPrivateIP 是否是内网IP
func (u *utils) IsPrivateIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}

// Extract returns a real ip
func (u *utils) Extract(addr string) (string, error) {
	// if addr specified then its returned
	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]") {
		return addr, nil
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("Failed to get interface addresses! Err: %v", err)
	}

	var ipAddr []byte

	for _, rawAddr := range addrs {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}

		if ip.To4() == nil {
			continue
		}

		if !u.IsPrivateIP(ip.String()) {
			continue
		}

		ipAddr = ip
		break
	}

	if ipAddr == nil {
		return "", fmt.Errorf("No private IP address found, and explicit IP not provided")
	}

	return net.IP(ipAddr).String(), nil
}

// IPs returns all known ips
func (u *utils) IPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ipAddrs []string

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			ipAddrs = append(ipAddrs, ip.String())
		}
	}

	return ipAddrs
}

//GetIntranetIP
func (u *utils) GetIntranetIP() string {
	if ip_cache != "" {
		return ip_cache
	}

	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return "127.0.0.1"
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip_cache = ipnet.IP.String()
						break
					}
				}
			}
		}
	}

	return ip_cache
}

//GetIntranetIPByName
func (u *utils) GetIntranetIPByName(name string) string {
	if ip_byName != "" {
		return ip_byName
	}

	iface, err := net.InterfaceByName(name)
	if err != nil {
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ip_byName = ipnet.IP.String()
			break
		}
	}
	return ip_byName
}

// IPToUInt32
func (u *utils) IPToUInt32(ipnr net.IP) uint32 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}

// UInt32ToIP
func (u *utils) UInt32ToIP(intIP uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

//ParseDns
func (u *utils) ParseDns(strDns string) bool {
	ns, err := net.LookupHost(strDns)
	if err != nil {
		fmt.Printf("error: %v, failed to parse %v\n", err, strDns)
		return false
	}

	if len(ns) <= 0 {
		return false
	}

	return true
}

//-------------

func (u *utils) GetUUIDV4() (uuidHex string) {
	uuidV4, _ := uuid.NewV4()
	uuidHex = hex.EncodeToString(uuidV4.Bytes())
	return
}

func (u *utils) GetMD5Base64(bytes []byte) (base64Value string) {
	md5Ctx := md5.New()
	md5Ctx.Write(bytes)
	md5Value := md5Ctx.Sum(nil)
	base64Value = base64.StdEncoding.EncodeToString(md5Value)
	return
}

func (u *utils) ByteToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func (u *utils) Base64ToByte(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

func (u *utils) GetGMTLocation() (*time.Location, error) {
	if loadLocationFromTZData != nil && tZData != nil {
		return loadLocationFromTZData("GMT", tZData)
	} else {
		return time.LoadLocation("GMT")
	}
}

func (u *utils) GetTimeInFormatISO8601() (timeStr string) {
	gmt, err := u.GetGMTLocation()

	if err != nil {
		panic(err)
	}
	return time.Now().In(gmt).Format("2006-01-02T15:04:05Z")
}

func (u *utils) GetTimeInFormatRFC2616() (timeStr string) {
	gmt, err := u.GetGMTLocation()

	if err != nil {
		panic(err)
	}
	return time.Now().In(gmt).Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

func (u *utils) InitStructWithDefaultTag(bean interface{}) {
	configType := reflect.TypeOf(bean)
	for i := 0; i < configType.Elem().NumField(); i++ {
		field := configType.Elem().Field(i)
		defaultValue := field.Tag.Get("default")
		if defaultValue == "" {
			continue
		}
		setter := reflect.ValueOf(bean).Elem().Field(i)
		switch field.Type.String() {
		case "int":
			intValue, _ := strconv.ParseInt(defaultValue, 10, 64)
			setter.SetInt(intValue)
		case "time.Duration":
			intValue, _ := strconv.ParseInt(defaultValue, 10, 64)
			setter.SetInt(intValue)
		case "string":
			setter.SetString(defaultValue)
		case "bool":
			boolValue, _ := strconv.ParseBool(defaultValue)
			setter.SetBool(boolValue)
		}
	}
}

// 当前时间的时间戳
func (u *utils) NowUnix() int {
	return int(time.Now().In(SysTimeLocation).Unix())
}

// 将unix时间戳格式化为yyyymmdd H:i:s格式字符串
func (u *utils) FormatFromUnixTime(t int64) string {
	if t > 0 {
		return time.Unix(t, 0).Format(SysTimeform)
	} else {
		return time.Now().Format(SysTimeform)
	}
}

// 将unix时间戳格式化为yyyymmdd格式字符串
func (u *utils) FormatFromUnixTimeShort(t int64) string {
	if t > 0 {
		return time.Unix(t, 0).Format(SysTimeformShort)
	} else {
		return time.Now().Format(SysTimeformShort)
	}
}

// 转化为yyyymmdd H:i:s格式字符串
func (u *utils) FormatUnixTime(year int, month int, day int) string {
	return time.Now().AddDate(year, month, day).Format(SysTimeform)
}

// 转化为yyyymmdd格式字符串
func (u *utils) FormatUnixTimeShort(year int, month int, day int) string {
	return time.Now().AddDate(year, month, day).Format(SysTimeformShort)
}

// 转化为yyyymm格式字符串
func (u *utils) FormatUnixTimeYm(year int, month int, day int) string {
	return time.Now().AddDate(year, month, day).Format(TimeformYm)
}

// 转化14位字符时间为标准时间格式 20190724151558 -> 2019-07-24 15:15:58
func (u *utils) FormatStringToStandTimeString(str string) string {
	if len(str) != 14 {
		return ""
	}
	sp := strings.Split(str, "")
	sb := strings.Builder{}
	sb.WriteString(strings.Join(sp[:4], ""))
	sb.WriteString("-")
	sb.WriteString(strings.Join(sp[4:6], ""))
	sb.WriteString("-")
	sb.WriteString(strings.Join(sp[6:8], ""))
	sb.WriteString(" ")
	sb.WriteString(strings.Join(sp[8:10], ""))
	sb.WriteString(":")
	sb.WriteString(strings.Join(sp[10:12], ""))
	sb.WriteString(":")
	sb.WriteString(strings.Join(sp[12:14], ""))
	return sb.String()
}

// 将字符串转成时间
func (u *utils) ParseTime(str string) (time.Time, error) {
	//return time.ParseInLocation(SysTimeform, str, SysTimeLocation)
	return dateparse.ParseIn(str, SysTimeLocation)
}

// 将字符串转成时间
//func (u *utils) ParseAny(str string) (time.Time, error) {
//	return dateparse.ParseIn(str, SysTimeLocation)
//}

// 判断unix时间是不是昨天
func (u *utils) IsYesToday(t int64) bool {
	ymd := u.FormatFromUnixTimeShort(t)

	yesDay := time.Now().AddDate(0, 0, -1)
	yd := u.FormatFromUnixTimeShort(yesDay.Unix())
	//fmt.Println(ymd,yd)
	if ymd == yd {
		return true
	} else {
		return false
	}
}

// 判断2个unix时间t1 是不是 t2的昨天
func (u *utils) IsYesTodayByTime(t1, t2 int64) bool {
	ymd := u.FormatFromUnixTimeShort(t1)
	t2s := u.FormatFromUnixTime(t2)
	parseTime, err := u.ParseTime(t2s)
	if err != nil {
		return false
	}

	yesDay := parseTime.AddDate(0, 0, -1)
	yd := u.FormatFromUnixTimeShort(yesDay.Unix())

	if ymd == yd {
		return true
	} else {
		return false
	}
}

func (u *utils) GetAnyDayZeroUnixTimestamp(t int) int64 {
	timeStr := time.Now().Format(SysTimeformShort)
	t2, _ := time.ParseInLocation(SysTimeformShort, timeStr, time.Local)
	return t2.AddDate(0, 0, t).Unix()
}

// 得到指定最大值的一个随机数
func (u *utils) Random(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if max < 1 {
		return r.Int()
	} else {
		return r.Intn(max)
	}
}

// RandomString 生成num个随机0-9,A-Z随机数
func (u *utils) RandomString(num int) string {
	chars := []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
	res := strings.Builder{}
	for i := 0; i < num; i++ {
		id := u.Random(62)
		res.WriteString(chars[id])
	}
	return res.String()
}

// RandomInt 生成num个随机0-9随机数
func (u *utils) RandomInt64(num int) (int64, error) {
	chars := []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
	res := strings.Builder{}
	for i := 0; i < num; i++ {
		id := u.Random(9)
		res.WriteString(chars[id])
	}
	return strconv.ParseInt(res.String(), 10, 64)
}

// 生成指定范围内的一个随机数
func (u *utils) RandRangeInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min+1)
}

// 生成指定范围内的一个随机数
func (u *utils) RandRangeInt64(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Int63n(max-min+1)
}

// 对字符串进行签名
func (u *utils) CreateSign(str string) string {
	signSecret := []byte("0123456789abcdef")
	str = string(signSecret) + str
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return sign
}

// addslashes() 函数返回在预定义字符之前添加反斜杠的字符串。
// 预定义字符是：
// 单引号（'）
// 双引号（"）
// 反斜杠（\）
func (u *utils) Addslashes(str string) string {
	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

// stripslashes() 函数删除由 addslashes() 函数添加的反斜杠。
func (u *utils) Stripslashes(str string) string {
	dstRune := []rune{}
	strRune := []rune(str)
	strLenth := len(strRune)
	for i := 0; i < strLenth; i++ {
		if strRune[i] == []rune{'\\'}[0] {
			i++
		}
		dstRune = append(dstRune, strRune[i])
	}
	return string(dstRune)
}

// 将字符串的IP转化为数字
func (u *utils) Ip4toInt(ip string) int64 {
	bits := strings.Split(ip, ".")
	if len(bits) == 4 {
		b0, _ := strconv.Atoi(bits[0])
		b1, _ := strconv.Atoi(bits[0])
		b2, _ := strconv.Atoi(bits[0])
		b3, _ := strconv.Atoi(bits[0])
		var sum int64
		sum += int64(b0) << 24
		sum += int64(b1) << 16
		sum += int64(b2) << 8
		sum += int64(b3)
		return sum
	} else {
		return 0
	}
}

// 得到当前时间到下一天零点的延时
func (u *utils) NextDayDuration() time.Duration {
	year, month, day := time.Now().Add(time.Hour * 24).Date()
	next := time.Date(year, month, day, 0, 0, 0, 0, SysTimeLocation)
	return next.Sub(time.Now())
}

// 从接口类型安全获取到int64
func (u *utils) GetInt64(i interface{}, d int64) int64 {
	if i == nil {
		return d
	}
	switch i.(type) {
	case string:
		num, err := strconv.Atoi(i.(string))
		if err != nil {
			return d
		} else {
			return int64(num)
		}
	case []byte:
		bits := i.([]byte)
		if len(bits) == 8 {
			return int64(binary.LittleEndian.Uint64(bits))
		} else if len(bits) <= 4 {
			num, err := strconv.Atoi(string(bits))
			if err != nil {
				return d
			} else {
				return int64(num)
			}
		}
	case uint:
		return int64(i.(uint))
	case uint8:
		return int64(i.(uint8))
	case uint16:
		return int64(i.(uint16))
	case uint32:
		return int64(i.(uint32))
	case uint64:
		return int64(i.(uint64))
	case int:
		return int64(i.(int))
	case int8:
		return int64(i.(int8))
	case int16:
		return int64(i.(int16))
	case int32:
		return int64(i.(int32))
	case int64:
		return i.(int64)
	case float32:
		return int64(i.(float32))
	case float64:
		return int64(i.(float64))
	}
	return d
}

// 从接口类型安全获取到字符串类型
func (u *utils) GetString(str interface{}, d string) string {
	if str == nil {
		return d
	}
	switch str.(type) {
	case string:
		return str.(string)
	case []byte:
		return string(str.([]byte))
	}
	return fmt.Sprintf("%v", str)
}

// 从map中得到指定的key
func (u *utils) GetInt64FromMap(dm map[string]interface{}, key string, dft int64) int64 {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return u.GetInt64(data, dft)
}

// 从map中得到指定的key
func (u *utils) GetInt64FromStringMap(dm map[string]string, key string, dft int64) int64 {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return u.GetInt64(data, dft)
}

// 从map中得到指定的key
func (u *utils) GetStringFromMap(dm map[string]interface{}, key string, dft string) string {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return u.GetString(data, dft)
}

// 从map中得到指定的key
func (u *utils) GetStringFromStringMap(dm map[string]string, key string, dft string) string {
	data, ok := dm[key]
	if !ok {
		return dft
	}
	return data
}

//map[string]interface{} 按照给定顺序序列化输出
func (u *utils) MarshalMap(mm map[string]interface{}) (string, error) {
	return mapEncoder(mm)
}

// map根据key排序放入list并排序
//func (u *utils)
// RemoteIp 返回远程客户端的 IP，如 192.168.1.1
func (u *utils) RemoteIp(ctx iris.Context) string {
	req := ctx.Request()
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
