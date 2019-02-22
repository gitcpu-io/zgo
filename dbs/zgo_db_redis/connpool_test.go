package zgo_db_redis

import (
	"fmt"
	"github.com/mediocregopher/radix"
	"testing"
	"time"
)

func TestInit(t *testing.T) {

	r := NewRedisResource()
	client := r.GetRedisClient()

	for i := 0; i < 10000; i++ {
		go func() {
			var result string

			client.Do(radix.Cmd(&result, "get", "forbidden"))
			fmt.Println(result)
		}()
	}

	time.Sleep(5 * time.Second)

}
