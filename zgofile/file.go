package zgofile

import (
	"git.zhugefang.com/gocore/zgo/config"
	"io"
)

type File interface {
	Put(name string, input io.Reader) (int64, error)
	Get(name string, output io.Writer) (int64, error)
	Size(name string) (int64, error)
}

var (
	FileStore File
)

func InitFile() {
	switch config.Conf.File.Type {
	case "local":
		FileStore = New(config.Conf.File.Home)
		break
	case "s3":
		break
	default:
		FileStore = New(config.Conf.File.Home)
	}
}
