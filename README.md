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

##zgo 测试环境
阿里云内网
10.45.146.41
阿里云公网
123.56.173.28

数字是端口号，供测试zgo admin使用，跑在docker里

2个mysql
3307
3308

2个mongo
27018
27019

2个redis
6380
6381

1个kafka
9202

1个nsq
4150
管理页面
http://123.56.173.28:4171/

1个etcd
2380
管理页面
http://47.93.163.209:9097/
