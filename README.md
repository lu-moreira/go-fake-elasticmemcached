# go-fake-elasticmemcached

A golang TCP server which simulates AWS's ElasticCache backed by memcached. 

Inspired heavily by https://github.com/stevenjack/fake_elasticache and https://github.com/dazoakley/fake-elasticache

# Requirements

- docker installed and up running
- [Preferred] `memcached` docker image already downloaded

# Usage


> Start TCP server at localhost:11210

```
go run ./cmd/* -numnodes 10
```

- numnodes, defines the number of memcached docker containers the server must start previously to accept any new connections.

> Connecting and get memcached nodes information (as most of AWS clients does)

```
telnet localhost 11210
config get cluster
```

> A sample response can be 
```
CONFIG cluster 0 249
1
localhost|127.0.0.1|1120 localhost|127.0.0.1|1121 localhost|127.0.0.1|1122 localhost|127.0.0.1|1123 localhost|127.0.0.1|1124 localhost|127.0.0.1|1125 localhost|127.0.0.1|1126 localhost|127.0.0.1|1127 localhost|127.0.0.1|1128 localhost|127.0.0.1|1129

END
```
