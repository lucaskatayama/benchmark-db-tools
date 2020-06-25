# benchmark-db-tools

This benchmark is based on this post
https://medium.com/@rocketlaunchr.cloud/how-to-benchmark-dbq-vs-sqlx-vs-gorm-e814caacecb5

# Running

```
$ docker-compose up -d
```

```
$ go test ./benchmarks -run=XXX -bench=.
```

# Results
```
goos: darwin
goarch: amd64
pkg: github.com/lucaskatayama/benchmark-db/benchmarks
BenchmarkAll/dbq__limit:5-8                 1077           1124689 ns/op
BenchmarkAll/sqlx_limit:5-8                 1038           1026922 ns/op
BenchmarkAll/gorm_limit:5-8                 1068           1102079 ns/op
BenchmarkAll/std__limit:5-8                 1000           1081778 ns/op
============
BenchmarkAll/dbq__limit:50-8                1086           1175158 ns/op
BenchmarkAll/sqlx_limit:50-8                 858           1417824 ns/op
BenchmarkAll/gorm_limit:50-8                 734           1566976 ns/op
BenchmarkAll/std__limit:50-8                 890           1357056 ns/op
============
BenchmarkAll/dbq__limit:500-8                550           2224541 ns/op
BenchmarkAll/sqlx_limit:500-8                516           2295445 ns/op
BenchmarkAll/gorm_limit:500-8                334           3648389 ns/op
BenchmarkAll/std__limit:500-8                510           2176116 ns/op
============
BenchmarkAll/dbq__limit:10000-8               57          20568259 ns/op
BenchmarkAll/sqlx_limit:10000-8               50          21023772 ns/op
BenchmarkAll/gorm_limit:10000-8               24          46494037 ns/op
BenchmarkAll/std__limit:10000-8               51          19803001 ns/op
============
PASS
ok      github.com/lucaskatayama/benchmark-db/benchmarks        22.494s
```