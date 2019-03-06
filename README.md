#zgo Engine

###测试方法使用：进入到比如zgonsq目录下执行，生成相应的.out，并通过go tool pprof查看

// 查看测试代码覆盖率
go test -coverprofile=c.out
go tool cover -html=c.out

// 查看cpu使用
go test -bench . -cpuprofile cpu.out
go tool pprof cpu.out

// 查看内存使用
go test -memprofile mem.out
go tool pprof mem.out

执行pprof后，然后输入web  或是quit 保证下载了svg
https://graphviz.gitlab.io/_pages/Download/Download_source.html
下载graphviz-2.40.1后进入目录
./configure
make
make install


docker-compose up -d

[![Go Report Card](https://goreportcard.com/badge/github.com/zgo-io/zgo?style=flat-square)](https://goreportcard.com/report/github.com/zgo-io/zgo)

**Note**: The `master` branch 


![zgo Logo](logos/zgo-horizontal-color.svg)

being:

* *Simple*: well-defined, user-facing API (gRPC)
* *Secure*: automatic TLS with optional client cert authentication


## Community meetings

*Community meeting will resume at 11:00 am on Thursday, January 10th, 2019.*


```
One tap mobile
+14086380986,,916003437# US (San Jose)
+16465588665,,916003437# US (New York)

Dial by location
        +1 408 638 0986 US (San Jose)
        +1 646 558 8665 US (New York)

Meeting ID: 916 003 437
```


## Getting started

### Getting zgo

The easiest way to get zgo is to use one of the pre-built release binaries which are available for OSX, Linux, Windows, and Docker on the [release page][github-release].

[dl-build]: ./Documentation/dl_build.md#build-the-latest-version

### Running zgo

First start a single-member cluster

If zgo is installed using the [pre-built release binaries][github-release], run it from the installation location as below:

```bash
/tmp/zgo-download-test/
```

## Contact

- Mailing list: [zgo-dev](https://groups.google.com/forum/?hl=en#!forum/zgo-dev)
- IRC: #[zgo](irc://irc.freenode.org:6667/#zgo) on freenode.org
