# cup

![It's fine!](http://z.ndr.su/go_fine.jpg)

## Ideas

- [easyjson](https://github.com/mailru/easyjson)
- [fasthttp](https://github.com/valyala/fasthttp)
- [dep](https://github.com/golang/dep)

## Install

### Prepare go code

```
go get -d github.com/eaglemoor/hlc

make deps
make gen
```

### Run with docker

```
make run
```

### Run native

#### Prepare test data

```
make data-unpack-train
```

#### Run

```
go build -o ./bin/app && ./bin/app -data-path ./data/data -http :3000

or

make run
```

## Run test

```
make tester
```

## Benchmarks

### net/http

```
Running 10s test @ http://127.0.0.1:3000
  8 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.22ms  311.40us   8.50ms   81.79%
    Req/Sec    11.26k   598.64    15.62k    72.99%
  903830 requests in 10.10s, 1.13GB read
  Socket errors: connect 0, read 61, write 0, timeout 0
Requests/sec:  89473.87
Transfer/sec:    114.26MB
```

### fasthttp

```
Running 10s test @ http://127.0.0.1:3000
  8 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.94ms  229.96us   6.01ms   83.22%
    Req/Sec    12.90k   768.43    14.46k    68.77%
  1036356 requests in 10.10s, 1.31GB read
  Socket errors: connect 0, read 52, write 0, timeout 0
Requests/sec: 102596.05
Transfer/sec:    132.77MB
```

### fasthttp + fasthttprouter

```
Running 10s test @ http://127.0.0.1:3000
  8 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.99ms  251.82us   5.38ms   82.00%
    Req/Sec    12.60k   833.37    14.14k    68.28%
  1011510 requests in 10.10s, 1.28GB read
  Socket errors: connect 0, read 60, write 0, timeout 0
Requests/sec: 100144.41
Transfer/sec:    129.60MB
```