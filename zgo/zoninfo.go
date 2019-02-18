package zgo

import "git.zhugefang.com/gocore/zgo/logic/zgo_resource"

var ZoneInfo zoneInfoer

func init() {
	ZoneInfo = zgo_resource.NewZoneInfo()
}

type zoneInfoer interface {
	GetTZData(string) ([]byte, bool)
}
