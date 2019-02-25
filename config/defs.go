package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
)

var (
	Env string
	Nsq map[string][]string
	Es map[string][]string
)

func InitConfig(e string) {
	initConfig(e)
}

func initConfig(e string) {
	_, f, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	cf := fmt.Sprintf("%s/%s.json", filepath.Dir(f), e)
	//fmt.Println(cf)
	viper.SetConfigFile(cf)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println("zgo engine is started on the ... ", viper.GetString("env"))

	//nsq地址
	Env = viper.GetString("env")
	Nsq = viper.GetStringMapStringSlice("nsq")
	Es = viper.GetStringMapStringSlice("es")

}
