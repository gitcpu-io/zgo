package zgo

import (
	"github.com/rubinus/zgo/logic/zgo_utils"
	"time"
)

var Utils utilser

func init() {
	Utils = zgo_utils.NewUtils()
}

type utilser interface {
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
