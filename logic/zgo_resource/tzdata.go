package zgo_resource

type zoneinfo struct{}

func NewZoneInfo() *zoneinfo {
	return &zoneinfo{}
}

func (z *zoneinfo) GetTZData(name string) ([]byte, bool) {
	data, ok := files["zoneinfo/"+name]
	return data, ok
}
