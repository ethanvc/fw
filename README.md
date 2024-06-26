FastWriter implements buffered I/O, which is fast than bufio package and have
better P99 performance.
Run benchmarks:
```shell
go test -bench . -benchmem -benchtime 3s
```

If you need file rotation function, use
[lumberjack](https://github.com/natefinch/lumberjack)
as the underline writer.
Similar library:
- [BufferedWriteSyncer](https://pkg.go.dev/go.uber.org/zap@v1.27.0/zapcore#BufferedWriteSyncer)

Benchmark result:
```shell
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/ethanvc/fw
Benchmark_Lumberjack-12                   300721              4065 ns/op               0 B/op          0 allocs/op
Benchmark_FileWriter-12                   418149              2816 ns/op               0 B/op          0 allocs/op
Benchmark_MemoryMapWriter-12             1271971               949.8 ns/op             0 B/op          0 allocs/op
Benchmark_BatchWriter-12                 1272800               943.2 ns/op             0 B/op          0 allocs/op
Benchmark_MultiBufferWriter-12            781482              1769 ns/op              24 B/op          0 allocs/op
Benchmark_FastWriter-12                  1385382               864.9 ns/op             1 B/op          0 allocs/op
Benchmark_NopWriter-12                   1295154               927.8 ns/op             0 B/op          0 allocs/op
Benchmark_BufioWriter-12                 1246978               967.4 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/ethanvc/fw   14.944s
```