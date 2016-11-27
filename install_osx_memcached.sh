# https://www.hacksparrow.com/install-memcached-on-mac-os-x.html
# http://mac-dev-env.patrickbougie.com/memcached/

# https://github.com/memcached/memcached/wiki/TutorialCachingStory


# brew install libevent

wget http://www.memcached.org/files/memcached-1.4.31.tar.gz
tar -xzvf memcached-1.4.31.tar.gz

rm memcached-1.4.31.tar.gz
mv memcached-1.4.31 memcached

cd memcached
./configure --prefix=/usr/local/memcached --with-libevent=/usr/local/lib/libevent

make
sudo make install

#./memcached -l 127.0.0.1 -vv
