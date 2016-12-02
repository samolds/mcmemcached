#!/bin/bash

ARCH=""

# install memcached deps
if type "apt-get" > /dev/null 2>&1 && ! type "brew" > /dev/null 2>&1; then
  # apt-get exists but brew does not
  echo ""
  echo "sudo apt-get -y update"
  sudo apt-get -y update
  echo ""
  echo "sudo apt-get -y install libevent-dev"
  sudo apt-get -y install libevent-dev
  ARCH="unix"
elif ! type "apt-get" > /dev/null 2>&1 && type "brew" > /dev/null 2>&1; then
  # apt-get does not exist but brew does
  echo ""
  echo "brew install libevent"
  brew install libevent
  ARCH="osx"
else
  echo ""
  echo "Either 'apt-get' or 'brew' are required for this installation to proceed"
  exit 1
fi


# install memcached
cd memcached

echo ""
echo "./configure --prefix=/usr/local/memcached"
if [ "$ARCH" = "osx" ]; then
  ./configure --prefix=/usr/local/memcached --with-libevent=/usr/local/lib/libevent
elif [ "$ARCH" = "unix" ]; then
  # update timestamps for weird bug explained here: http://stackoverflow.com/a/33279062
  touch aclocal.m4 configure
  touch Makefile.am
  touch Makefile.in

  ./configure --prefix=/usr/local/memcached
else
  exit 1
fi

echo ""
echo "make && sudo make install"
make
#make test
sudo make install
cd ..


# download and install go if you don't already have it
if ! type "go" > /dev/null 2>&1; then
  echo ""
  echo "installing go"
  wget https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.7.3.linux-amd64.tar.gz
  rm -f go1.7.3.linux-amd64.tar.gz
  echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
fi
echo "export GOPATH=`pwd`" >> ~/.bashrc
export GOPATH=`pwd`


# download  and install go deps
echo ""
echo "go get github.com/bradfitz/gomemcache/memcache"
/usr/local/go/bin/go get github.com/bradfitz/gomemcache/memcache
