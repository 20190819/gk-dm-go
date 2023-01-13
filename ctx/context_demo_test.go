package ctx

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Cache interface {
	Get(key string) (string, error)
}

type OtherCache interface {
	GetValue(ctx context.Context, key string) (any, error)
}

// CacheAdapter 适配器强调不同接口之间的进行适配
// 装饰器强调额外功能扩展
type CacheAdapter struct {
	Cache
}

func (c *CacheAdapter) GetValue(ctx context.Context, key string) (any, error) {
	return c.Cache.Get(key)
}

// 已有的，不是线程安全
type memoryMap struct {

	// 侵入式的写法
	// 需要测试这个类
	// 如果是三方依赖，可能都改不了
	//lock sync.RWMutex

	m map[string]string
}

func (mp *memoryMap) Get(key string) (string, error) {
	return mp.m[key], nil
}

type SafeMap struct {
	mm   *memoryMap
	lock sync.RWMutex
}

func (s *SafeMap) Get(key string) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.mm.Get(key)
}

func TestSafeMap(t *testing.T) {
	var sfm = &SafeMap{
		mm: &memoryMap{m: map[string]string{
			"url":    "baidu.com",
			"scheme": "https",
		}},
	}

	for i := 0; i < 10; i++ {
		value, _ := sfm.Get("scheme")
		assert.Equal(t, "https", value)
	}
}

func TestErrGroup(t *testing.T) {
	group, ctx := errgroup.WithContext(context.Background())
	var result int64
	for i := 0; i < 10; i++ {
		delta := i
		group.Go(func() error {
			atomic.AddInt64(&result, int64(delta))
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		t.Fatal(err)
	}
	fmt.Println("err: ", ctx.Err())
	fmt.Println(result)
}

func TestBusinessTimeout(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	end := make(chan struct{}, 1)
	go func() {
		MyBusiness()
		end <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("timeout")
	case <-end:
		fmt.Println("biz end")
	}

}

func MyBusiness() {
	time.Sleep(1500 * time.Millisecond)
	fmt.Println("hello,world")
}

func TestParentValueContext(t *testing.T) {
	ctx := context.Background()
	childCtx := context.WithValue(ctx, "map", map[string]string{"url": "baidu.com"})
	subChildCtx := context.WithValue(childCtx, "key1", "value1")

	m := subChildCtx.Value("map").(map[string]string)
	m["kk"] = "vv"

	val := subChildCtx.Value("key1")
	fmt.Println(val)
	val = childCtx.Value("map")
	fmt.Println(val)
}

func TestParentCtx(t *testing.T) {
	ctx := context.Background()
	dlCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Minute))
	ctx = context.WithValue(dlCtx, "key", 123)
	cancel()
	fmt.Println(ctx.Err())
}
