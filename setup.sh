cd memcached
./configure --prefix=/usr/local/memcached --with-libevent=/usr/local/lib/libevent

make
sudo make install

#./memcached -l 127.0.0.1 -vv

export GOPATH=`pwd`
go get github.com/bradfitz/gomemcache/memcache
