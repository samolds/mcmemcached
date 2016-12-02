#!/bin/sh

PORT1=11211
PORT2=11212

if [[ $# -eq 1 ]]; then
  PORT1=$1
elif [[ $# -eq 2 ]]; then
  PORT1=$1
  PORT2=$2
fi

# SUPER VERBOSE
#./memcached/memcached -vvvv -d -p $PORT1
#./memcached/memcached -vvvv -d -p $PORT2

./memcached/memcached -d -p $PORT1
./memcached/memcached -d -p $PORT2

sleep 2 # wait for 2 second for memcache servers to start

go run client/client.go


echo "\n\nYou may want to kill the memcache servers now to prevent zombie processes"
ps axl | grep memcached
echo "You can easily kill these processes with:"
echo "$ kill -9 <pid (second column)>"
