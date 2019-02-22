package zgo_db_redis

import (
	"fmt"
	"github.com/mediocregopher/radix"
	"time"
)

var client *radix.Pool

func init() {
	customConnFunc := func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialTimeout(10*time.Second), radix.DialSelectDB(9), radix.DialAuthPass(""),
		)
	}
	c, err := radix.NewPool("tcp", "127.0.0.1:6379", 100, radix.PoolConnFunc(customConnFunc))
	if err != nil {
		fmt.Println("redis ", err)
	}
	client = c
}
