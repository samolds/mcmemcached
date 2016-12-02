#!/bin/sh

#PORT=8000
#
#if [[ $# -eq 1 ]]; then
#  PORT=$1
#fi

#./memcached/memcached -vvvv -d -p 11211 &
#./memcached/memcached -vvvv -d -p 11212 &

./memcached/memcached -d -p 11211 &
./memcached/memcached -d -p 11212 &

sleep 1 # wait for 1 second for memcache servers to start

go run client/client.go
