package zgofile

import (
  "os"
  "strings"
  "testing"
)

func TestFileLocal(t *testing.T) {
  f := New("/tmp")
  for i := 0; i < 10000; i++ {
    go func() {
      input := strings.NewReader("hello world")
      pn, err := f.Put("/x/a.txt", input)
      if err != nil {
        t.Error(err)
      }
      t.Log("put bytes num:", pn)

      gn, err := f.Get("/x/a.txt", os.Stdout)
      if err != nil {
        t.Error(err)
      }
      t.Log("get bytes num:", gn)

      if gn != pn {
        t.Error("lose some things")
      }

      size, err := f.Size("/x/a.txt")
      if err != nil {
        t.Error(err)
      }
      t.Log("file size:", size)
      if gn != size {
        t.Error("error file size")
      }
    }()
  }

}
