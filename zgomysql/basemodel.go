package zgomysql

// 基类 所有
type BaseModel struct {
	Id uint32 `json:"id"`
}

type MysqlBaser interface {
	GetID() uint32
	SetID(id uint32)
	TableName() string
	Omit() string
}

func (b *BaseModel) GetID() uint32 {
	return b.Id
}
func (b *BaseModel) SetID(id uint32) {
	b.Id = id
}

func (b *BaseModel) GetDB() string {
	return ""
}

func (b *BaseModel) Omit() string {
	return ""
}
