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


## To make changes to memcache source and run full tamale
* `cd memcache`
* Edit memcache source as desired
* `make`
* `cd ..`
* `./run.sh`


# References
* [Facebook Memcache](http://www.cs.utah.edu/~stutsman/cs6963/public/papers/memcached.pdf)
* [Saving Cash by Using Less Cache](http://www.cs.cmu.edu/~harchol/Papers/HotCloud12.pdf)
