# McMemcached
Class Project for University of Utah [Distributed Systems](http://www.cs.utah.edu/~stutsman/cs6963)


## Dependencies
* `wget`
* `apt-get` (ubuntu) OR `brew` (osx)


## Getting client memcache running up against memcache server on fresh ubuntu

```sh
git clone https://cs6963.utah.systems/u0810837/memcached.git mcmemcached
cd mcmemcached
./setup.sh
source ~/.bashrc

./run.sh
```


## Experiments

There are 4 main experiments, located in `src/experiments/`. If you have the
memcache servers running, you can run them with:

```
go run src/experiments/<dir>/<filename>.go
```

OR, you can use the provided `./run.sh` shell script. If provided 1 of the
following arguments, it will spin up the memcache servers and run the associated
experiment against the servers:

| Command | Experiment Source |
| --- | --- |
| `./run.sh` | `src/experiments/1naive/client_zipf.go` |
| `./run.sh 1z` | `src/experiments/1naive/client_zipf.go` |
| `./run.sh 1n` | `src/experiments/1naive/client_n.go` |
| `./run.sh 2` | `src/experiments/2coldstart/client.go` |
| `./run.sh 3` | `src/experiments/3querycold/client.go` |
| `./run.sh 4` | `src/experiments/4hotkey/client.go` |


## To make changes to memcache source and run full tamale
* `cd memcache`
* Edit memcache source as desired
* `make`
* `cd ..`
* `./run.sh`


## To make plots with R:

```
setwd("~/Desktop")
ee<-read.csv("time_stats_1z.csv")
plot(x=ee$time, y=ee$mean_response_time, main="Experiment 1", ylab="Mean Response Time", xlab="Time (ms)", xlim=c(0,1200000), ylim=c(1,6))
```


# References
* [Facebook Memcache](http://www.cs.utah.edu/~stutsman/cs6963/public/papers/memcached.pdf)
* [Saving Cash by Using Less Cache](http://www.cs.cmu.edu/~harchol/Papers/HotCloud12.pdf)


# Contributors:

* [Keith Downie](https://github.com/kdownie)
* [Rehan Ghori](https://www.linkedin.com/in/mohammad-rehan-ghori-4402542)
* [Sam Olds](https://github.com/samolds)
