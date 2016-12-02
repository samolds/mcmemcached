# https://www.hacksparrow.com/install-memcached-on-mac-os-x.html
# http://mac-dev-env.patrickbougie.com/memcached/
# https://github.com/memcached/memcached/wiki/TutorialCachingStory

if ! type "wget" > /dev/null 2>&1; then
  echo "'wget' is required for this download to proceed"
  exit 1
fi

# download memcache source
wget http://www.memcached.org/files/memcached-1.4.31.tar.gz
tar -xzf memcached-1.4.31.tar.gz

rm memcached-1.4.31.tar.gz
mv memcached-1.4.31 memcached
