
https://github.com/natefinch/lumberjack

go test -bench . -benchmem -benchtime 3s

remove file rotation function form fast writer to concentrate on real important things.

benchmark result:
```shell
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/ethanvc/fw
Benchmark_Lumberjack-12                   418022              3178 ns/op               0 B/op          0 allocs/op
Benchmark_FileWriter-12                   424155              2813 ns/op               0 B/op          0 allocs/op
Benchmark_MemoryMapWriter-12             1275559               945.6 ns/op             0 B/op          0 allocs/op
Benchmark_BatchWriter-12                 1277884               944.4 ns/op             0 B/op          0 allocs/op
Benchmark_MultiBufferWriter-12            829839              1763 ns/op              24 B/op          0 allocs/op
Benchmark_FastWriter-12                  1388536               870.0 ns/op             1 B/op          0 allocs/op
Benchmark_NopWriter-12                   1306736               922.0 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/ethanvc/fw   13.832s
```