#!/bin/bash

PORT1=11211
PORT2=11212
PORT3=11213
PORT4=11214

if [ $# -eq 4 ]; then
  PORT1=$1
  PORT2=$2
  PORT3=$3
  PORT4=$4
fi

# SUPER VERBOSE
#./memcached/memcached -vvvv -d -p $PORT1
#./memcached/memcached -vvvv -d -p $PORT2
#./memcached/memcached -vvvv -d -p $PORT3
#./memcached/memcached -vvvv -d -p $PORT4

./memcached/memcached -d -p $PORT1
./memcached/memcached -d -p $PORT2
./memcached/memcached -d -p $PORT3
./memcached/memcached -d -p $PORT4

sleep 2 # wait for 2 second for memcache servers to start

go run src/experiments/1-naive/client.go

# kills all 4 memcache servers (as well as any other processes named "memcached"
killall -9 memcached
