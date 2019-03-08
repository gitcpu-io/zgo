package zgoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/satori/go.uuid"
	"io"
	"math/rand"
	"mygo/lottery/conf"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ip_cache = ""
)

var (
	privateBlocks []*net.IPNet
)

func init() {
	for _, b := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10"} {
		if _, block, err := net.ParseCIDR(b); err == nil {
			privateBlocks = append(privateBlocks, block)
		}
	}
}

var Utils Utilser

func init() {
	Utils = NewUtils()
}

type Utilser interface {
	//md5 对字符串md5
	Md5(s string) string
	Sha1(s string) string
	//获取当前时间时间戳，n= 10/13/19 位时间戳
	GetTimestamp(n int) int64
	//Marshal 序列化为json
	Marshal(in interface{}) ([]byte, error)
	//Unmarshal 反序列化为go 内存对象
	Unmarshal(message []byte, in interface{}) error
	NewDecoder(reader io.Reader) *jsoniter.Decoder
	NewEncoder(writer io.Writer) *jsoniter.Encoder
	//结构体转map[string]interface{}
	StructToMap(interface{}) map[string]interface{}
	// GrpcServiceMethod converts a gRPC method to a Go method
	GrpcServiceMethodConverts(m string) (string, string, error)

	ParseDns(strDns string) bool
	IPs() []string
	//是否是内网ip
	IsPrivateIP(ipAddr string) bool
	Extract(addr string) (string, error)
	//获取内网ip
	GetIntranetIP() string

	GetUUIDV4() string
	GetMD5Base64([]byte) string
	GetGMTLocation() (*time.Location, error)
	GetTimeInFormatISO8601() string
	GetTimeInFormatRFC2616() string
	//从一个map中返回a=123&b=456
	GetUrlFormedMap(map[string]string) string
	InitStructWithDefaultTag(interface{})

	NowUnix() int
	FormatFromUnixTime(t int64) string
	FormatFromUnixTimeShort(t int64) string
	ParseTime(str string) (time.Time, error)
	Random(max int) int
	CreateSign(str string) string
	Encrypt(key, text []byte) ([]byte, error)
	Decrypt(key, text []byte) ([]byte, error)
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
}

var loadLocationFromTZData func(name string, data []byte) (*time.Location, error) = nil

var tZData []byte = nil

var jsonIterator = jsoniter.ConfigCompatibleWithStandardLibrary

type utils struct{}

func NewUtils() Utilser {
	return &utils{}
}

//Marshal 序列化为json
func (u *utils) Marshal(res interface{}) ([]byte, error) {
	return jsonIterator.Marshal(res)
}

//Unmarshal 反序列化为go 内存对象
func (u *utils) Unmarshal(message []byte, in interface{}) error {
	return jsonIterator.Unmarshal(message, in)
}

func (u *utils) NewDecoder(reader io.Reader) *jsoniter.Decoder {
	return jsoniter.NewDecoder(reader)
}

func (u *utils) NewEncoder(writer io.Writer) *jsoniter.Encoder {
	return jsoniter.NewEncoder(writer)
}

//StructToMap 结构体转map[string]interface{}
func (u *utils) StructToMap(input interface{}) map[string]interface{} {
	var m map[string]interface{}
	b, _ := jsonIterator.Marshal(input)
	jsonIterator.Unmarshal(b, &m)
	return m
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

//Md5
func (u *utils) Md5(body string) string {
	md5 := md5.New()
	md5.Write([]byte(body))
	return hex.EncodeToString(md5.Sum(nil))
}

//Sha1
func (u *utils) Sha1(s string) string {
	r := sha1.Sum([]byte(s))
	return hex.EncodeToString(r[:])
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

//GetUrlFormedMap 从一个map中返回a=123&b=456
func (u *utils) GetUrlFormedMap(source map[string]string) (urlEncoded string) {
	urlEncoder := url.Values{}
	for key, value := range source {
		urlEncoder.Add(key, value)
	}
	urlEncoded = urlEncoder.Encode()
	return
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
	return int(time.Now().In(conf.SysTimeLocation).Unix())
}

// 将unix时间戳格式化为yyyymmdd H:i:s格式字符串
func (u *utils) FormatFromUnixTime(t int64) string {
	if t > 0 {
		return time.Unix(t, 0).Format(conf.SysTimeform)
	} else {
		return time.Now().Format(conf.SysTimeform)
	}
}

// 将unix时间戳格式化为yyyymmdd格式字符串
func (u *utils) FormatFromUnixTimeShort(t int64) string {
	if t > 0 {
		return time.Unix(t, 0).Format(conf.SysTimeformShort)
	} else {
		return time.Now().Format(conf.SysTimeformShort)
	}
}

// 将字符串转成时间
func (u *utils) ParseTime(str string) (time.Time, error) {
	return time.ParseInLocation(conf.SysTimeform, str, conf.SysTimeLocation)
}

// 得到一个随机数
func (u *utils) Random(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if max < 1 {
		return r.Int()
	} else {
		return r.Intn(max)
	}
}

// 对字符串进行签名
func (u *utils) CreateSign(str string) string {
	signSecret := []byte("0123456789abcdef")
	str = string(signSecret) + str
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return sign
}

// 对一个字符串进行加密
func (u *utils) Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	//if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	//	return nil, err
	//}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// 对一个字符串进行解密
func (u *utils) Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
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
	next := time.Date(year, month, day, 0, 0, 0, 0, conf.SysTimeLocation)
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
	return fmt.Sprintf("%s", str)
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
