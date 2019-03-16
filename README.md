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
