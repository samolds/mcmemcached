#!/bin/bash

EXPERIMENT='1z'
PORT1=11211
PORT2=11212
PORT3=11213
PORT4=11214

if [ $# -eq 1 ]; then
  EXPERIMENT=$1
elif [ $# -eq 4 ]; then
  PORT1=$1
  PORT2=$2
  PORT3=$3
  PORT4=$4
elif [ $# -eq 5 ]; then
  EXPERIMENT=$1
  PORT1=$2
  PORT2=$3
  PORT3=$4
  PORT4=$5
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

if [ "$EXPERIMENT" = "1z" ]; then
  go run src/experiments/1naive/client_zipf.go > data/time_stats_$EXPERIMENT.csv
elif [ "$EXPERIMENT" = "1n" ]; then
  go run src/experiments/1naive/client_n.go > data/time_stats_$EXPERIMENT.csv
elif [ "$EXPERIMENT" = "2" ]; then
  go run src/experiments/2coldstart/client.go > data/time_stats_$EXPERIMENT.csv
elif [ "$EXPERIMENT" = "3" ]; then
  go run src/experiments/3querycold/client.go > data/time_stats_$EXPERIMENT.csv
elif [ "$EXPERIMENT" = "4" ]; then
  go run src/experiments/4hotkey/client.go > data/time_stats_$EXPERIMENT.csv
fi

# kills all 4 memcache servers (as well as any other processes named "memcached"
killall -9 memcached
