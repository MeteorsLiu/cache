# Time-based Cache

It's like `useMemo` in React.

It caches the return value of function.

# Example
```
func getValue() (int, err) {
    // Some slow operation like SYSCALL
    return doSlow()
}

// cache the return value of getValue(), valid until 5s.
memo := cache.NewCache(getValue, 5*time.Second)
// use this value
fn(memo * 2)
```

# Benchmark
GetTCPInfo SYSCALL

```
BenchmarkCache
BenchmarkCache/cached
BenchmarkCache/cached         	32197653	        36.19 ns/op
BenchmarkCache/raw
BenchmarkCache/raw            	  834799	      1433 ns/op
BenchmarkTcpInfo
BenchmarkTcpInfo/cached
BenchmarkTcpInfo/cached       	32753490	        36.19 ns/op
BenchmarkTcpInfo/raw
BenchmarkTcpInfo/raw          	  583994	      2284 ns/op
BenchmarkUntil
BenchmarkUntil/until
BenchmarkUntil/until          	36392731	        32.43 ns/op
BenchmarkUntil/after
BenchmarkUntil/after          	22826443	        52.03 ns/op
```

SYSCALL spend 1400ns+ while cached value only spend 36ns.
The high frequency reading SYSCALL can be optimized by the cache function.
