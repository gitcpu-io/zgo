##第一步：安装pkg-config、cmake和openssl
brew install pkg-config

brew install cmake

brew install openssl

##第二步下载seabolt到你机器任意目录
git clone https://github.com/neo4j-drivers/seabolt

####配置环境变量，如果不存在~/.bash_profile，就建立：

####注意：1.0.2r是你本机openssl的版本，换到你的可能是 1.0.2s
####注意：你的目录   -- 就是seabolt的下载目录

--------------请务必替换1.0.2r 和 你的目录-----------

export OPENSSL_ROOT_DIR=/usr/local/Cellar/openssl/1.0.2r

export OPENSSL_LIBRARIES=/usr/local/Cellar/openssl/1.0.2r/lib

export PKG_CONFIG_PATH=$HOME/你的目录/seabolt/build/dist/share/pkgconfig

export C_INCLUDE_PATH=$HOME/你的目录/seabolt/build/dist/include

export DYLD_LIBRARY_PATH=$HOME/你的目录/seabolt/build/dist/lib

export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/lib/pkgconfig

export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/local/lib/pkgconfig

export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:/usr/local/lib

##使用source生效
source ~/.bash_profile

##第三步：在seabolt目录下执行
./make_release.sh


##进入到seabolt目录的中 把安装好的.pc文件copy过去
cd seabolt

sudo cp build/share/pkgconfig/seabolt17.pc build/share/pkgconfig/seabolt17-static.pc /usr/local/lib/pkgconfig

sudo cp build/lib/libseabolt17.1.dylib /usr/local/lib

##在你的项目所在目录执行下面这行
go build -ldflags "-r $(pkg-config --variable=libdir seabolt17)"



