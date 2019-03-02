package zgoutils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/satori/go.uuid"
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
