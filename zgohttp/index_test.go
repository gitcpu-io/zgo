package zgohttp

import (
  "fmt"
  "testing"
)

/*
@Time : 2019-06-18 19:22
@Author : rubinus.chu
@File : index_test
@project: zgo
*/

func TestZgohttp_Get(t *testing.T) {
  zh := zgohttp{}
  bytes, err := zh.Get("http://www.zhuge.com")
  if err != nil {
    fmt.Println(err)
  }
  fmt.Printf("%s", bytes)
}
