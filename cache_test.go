package cache

import (
	"crypto/rand"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestCache(t *testing.T) {
	var cnt int
	c := NewCache[int](func() (int, error) {
		cnt++
		return cnt, nil
	}, 5*time.Second)
	ct, _ := c.Get()
	if ct != 1 {
		t.Error("cache fail", ct, cnt)
	}
	time.Sleep(3 * time.Second)
	ct, _ = c.Get()
	if ct != 1 {
		t.Error("cache fail", ct, cnt)
	}
	time.Sleep(2 * time.Second)
	ct, _ = c.Get()
	if ct != 2 {
		t.Error("sleep cache fail", ct, cnt)
	}
}

func TestCacheConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	var cnt atomic.Int32
	c := NewCache[int32](func() (int32, error) {
		return cnt.Add(1), nil
	}, 5*time.Second)
	wg.Add(2)
	go func() {
		time.Sleep(6 * time.Second)
		ct, _ := c.Get()
		if ct != 2 {
			t.Error("cache fail 1", ct)
		}
		ct, _ = c.Get()
		if ct != 2 {
			t.Error("cache fail 1-2", ct)
		}
		wg.Done()
	}()

	go func() {
		time.Sleep(6 * time.Second)
		ct, _ := c.Get()
		if ct != 2 {
			t.Error("cache fail 2", ct)
		}
		ct, _ = c.Get()
		if ct != 2 {
			t.Error("cache fail 2-2", ct)
		}
		wg.Done()
	}()

	wg.Wait()
}

func BenchmarkCache(b *testing.B) {
	var buf [32]byte
	c := NewCache[int](func() (int, error) {
		_, err := rand.Read(buf[:])
		return 0, err
	}, time.Second)
	c.Get()
	b.ResetTimer()
	b.Run("cached", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Get()
		}
	})

	b.Run("raw", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rand.Read(buf[:])
		}
	})
}

func BenchmarkTcpInfo(b *testing.B) {
	c, err := net.Dial("tcp", "www.google.com:80")
	if err != nil {
		b.Error(err)
		return
	}
	defer c.Close()
	raw, _ := c.(*net.TCPConn).SyscallConn()

	getInfo := func() (info *unix.TCPInfo, err error) {
		raw.Control(func(fd uintptr) {
			info, err = unix.GetsockoptTCPInfo(int(fd), unix.IPPROTO_TCP, unix.TCP_INFO)
		})
		return
	}
	ca := NewCache[*unix.TCPInfo](getInfo, time.Second)
	ca.Get()

	b.ResetTimer()
	b.Run("cached", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ca.Get()
		}
	})

	b.Run("raw", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			getInfo()
		}
	})
}

func BenchmarkUntil(b *testing.B) {
	t1 := time.Now()
	b.Run("until", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			time.Until(t1)
		}
	})
	b.Run("after", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			time.Now().After(t1)
		}
	})
}
