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
value, err := memo.Get()
if err != nil {
    // error handle
    .... 
}
fn(value * 2)
```

# Benchmark
GetTCPInfo SYSCALL

```
BenchmarkTcpInfo
BenchmarkTcpInfo/cached
BenchmarkTcpInfo/cached       	32753490	        36.19 ns/op
BenchmarkTcpInfo/raw
BenchmarkTcpInfo/raw          	  583994	      2284 ns/op
```

SYSCALL spend 2284ns+ while cached value only spend 36ns. (63x faster)

The high frequency reading SYSCALL can be optimized by the cache function.
