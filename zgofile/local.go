package zgofile

import (
	"bufio"
	"io"
	"os"
	"path"
)

type Local struct {
	home string
}

func New(home ...string) *Local {
	sf := ""
	if len(home) != 0 {
		sf = home[0]
	}
	return &Local{sf}
}

func (file *Local) Get(name string, output io.Writer) (n int64, err error) {
	src := path.Join(file.home, name)
	fd, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	input := bufio.NewReader(fd)
	defer fd.Close()
	return io.Copy(output, input)
}

func (file *Local) Put(name string, input io.Reader) (n int64, err error) {
	dst := path.Join(file.home, name)
	parent := path.Dir(dst)
	_, err = os.Stat(parent)
	if os.IsNotExist(err) {
		err = os.MkdirAll(parent, 0700)
		if err != nil {
			return 0, err
		}
	}
	fd, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return 0, err
	}
	output := bufio.NewWriter(fd)
	defer func() {
		output.Flush()
		fd.Close()
	}()
	return io.Copy(output, input)
}

func (file *Local) Append(name string, input io.Reader) (n int64, err error) {
	dst := path.Join(file.home, name)
	parent := path.Dir(dst)
	_, err = os.Stat(parent)
	if os.IsNotExist(err) {
		err = os.MkdirAll(parent, 0700)
		if err != nil {
			return 0, err
		}
	}
	fd, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return 0, err
	}
	fd.Seek(0, 2)
	output := bufio.NewWriter(fd)
	defer func() {
		output.Flush()
		fd.Close()
	}()
	return io.Copy(output, input)
}

func (file *Local) Size(name string) (n int64, err error) {
	dst := path.Join(file.home, name)
	fi, err := os.Stat(dst)
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}
