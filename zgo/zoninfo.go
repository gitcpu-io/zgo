package zgo

import "git.zhugefang.com/gocore/zgo.git/logic/zgo_resource"

var ZoneInfo zoneInfoer

func init() {
	ZoneInfo = zgo_resource.NewZoneInfo()
}

type zoneInfoer interface {
	GetTZData(string) ([]byte, bool)
}
