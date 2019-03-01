package zgozoneinfo

type zoneinfo struct{}

func NewZoneInfo() ZoneInfoer {
	return &zoneinfo{}
}

func (z *zoneinfo) GetTZData(name string) ([]byte, bool) {
	data, ok := files["zoneinfo/"+name]
	return data, ok
}

type ZoneInfoer interface {
	GetTZData(string) ([]byte, bool)
}
