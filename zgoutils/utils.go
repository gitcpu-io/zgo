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
	"time"
)

var Utils Utilser

func init() {
	Utils = NewUtils()
}

type Utilser interface {
	Md5(body string) string
	Sha1(s string) string
	GetTimestamp(n int) int64
	Marshal(res interface{}) (string, error)
	ParseDns(strDns string) bool
	GetIntranetIp() string

	GetUUIDV4() string
	GetMD5Base64([]byte) string
	GetGMTLocation() (*time.Location, error)
	GetTimeInFormatISO8601() string
	GetTimeInFormatRFC2616() string
	GetUrlFormedMap(map[string]string) string
	InitStructWithDefaultTag(interface{})
}

var loadLocationFromTZData func(name string, data []byte) (*time.Location, error) = nil

var tZData []byte = nil

type utils struct{}

func NewUtils() Utilser {
	return &utils{}
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

//Marshal
func (u *utils) Marshal(res interface{}) (string, error) {
	jsonIterator := jsoniter.ConfigCompatibleWithStandardLibrary
	s, err := jsonIterator.Marshal(res)
	return string(s), err
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

var (
	ip_cache = ""
)

//GetIntranetIp
func (u *utils) GetIntranetIp() string {
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
