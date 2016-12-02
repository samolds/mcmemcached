ARCH=""

# install memcached deps
if type "apt-get" > /dev/null 2>&1 && ! type "brew" > /dev/null 2>&1; then
  # apt-get exists but brew does not
  sudo apt-get install libevent-dev
  ARCH="unix"
elif ! type "apt-get" > /dev/null 2>&1 && type "brew" > /dev/null 2>&1; then
  # apt-get does not exist but brew does
  brew install libevent
  ARCH="osx"
else
  echo "Either 'apt-get' or 'brew' are required for this installation to proceed"
  exit 1
fi


# install memcached
cd memcached

if [ "$ARCH" = "osx" ]; then
  ./configure --prefix=/usr/local/memcached --with-libevent=/usr/local/lib/libevent
elif [ "$ARCH" = "osx" ]; then
  ./configure --prefix=/usr/local/memcached
else
  exit 1
fi

make
make test
sudo make install
cd ..


# download and install go
wget https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.7.3.linux-amd64.tar.gz
rm -f go1.7.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPATH=`pwd`


# download  and install go deps
go get github.com/bradfitz/gomemcache/memcache
