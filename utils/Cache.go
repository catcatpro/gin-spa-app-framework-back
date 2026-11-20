package utils

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache 泛型缓存封装
type Cache[T any] struct {
	inner *cache.Cache
}

// New 创建缓存实例
func New[T any](defaultExpiration, cleanupInterval time.Duration) *Cache[T] {
	return &Cache[T]{
		inner: cache.New(defaultExpiration, cleanupInterval),
	}
}

// Set 设置缓存（使用默认过期时间）
func (c *Cache[T]) Set(key string, value T) {
	c.inner.Set(key, value, cache.DefaultExpiration)
}

// SetWithTTL 设置缓存并指定过期时间
func (c *Cache[T]) SetWithTTL(key string, value T, ttl time.Duration) {
	c.inner.Set(key, value, ttl)
}

// SetForever 设置永不过期
func (c *Cache[T]) SetForever(key string, value T) {
	c.inner.Set(key, value, cache.NoExpiration)
}

// Get 获取缓存值
func (c *Cache[T]) Get(key string) (T, bool) {
	var zero T
	val, found := c.inner.Get(key)
	if !found {
		return zero, false
	}

	result, ok := val.(T)
	if !ok {
		return zero, false
	}
	return result, true
}

// GetWithExpiration 获取缓存值及其过期时间点
func (c *Cache[T]) GetWithExpiration(key string) (T, time.Time, bool) {
	var zero T
	val, expiration, found := c.inner.GetWithExpiration(key)
	if !found {
		return zero, time.Time{}, false
	}

	result, ok := val.(T)
	if !ok {
		return zero, time.Time{}, false
	}
	return result, expiration, true
}

// GetRemainingTTL 获取剩余过期时间（Duration）
func (c *Cache[T]) GetRemainingTTL(key string) (time.Duration, bool) {
	_, expiration, found := c.inner.GetWithExpiration(key)
	if !found {
		return 0, false
	}

	// 检查是否永不过期
	if expiration.IsZero() {
		return 0, true // 0 表示永不过期
	}

	remaining := time.Until(expiration)
	if remaining <= 0 {
		return 0, false // 已过期
	}
	return remaining, true
}

// IsExpired 检查是否已过期
func (c *Cache[T]) IsExpired(key string) bool {
	_, expiration, found := c.inner.GetWithExpiration(key)
	if !found {
		return true
	}

	// 永不过期
	if expiration.IsZero() {
		return false
	}

	return time.Now().After(expiration)
}

// GetOrSet 获取缓存，如果不存在则设置并返回
func (c *Cache[T]) GetOrSet(key string, factory func() T) T {
	if val, found := c.Get(key); found {
		return val
	}

	val := factory()
	c.Set(key, val)
	return val
}

// GetOrSetWithTTL 获取缓存，如果不存在则设置（带过期时间）
func (c *Cache[T]) GetOrSetWithTTL(key string, ttl time.Duration, factory func() T) T {
	if val, found := c.Get(key); found {
		return val
	}

	val := factory()
	c.SetWithTTL(key, val, ttl)
	return val
}

// Delete 删除缓存
func (c *Cache[T]) Delete(key string) {
	c.inner.Delete(key)
}

// DeleteAll 删除所有缓存（清空）
func (c *Cache[T]) DeleteAll() {
	c.inner.Flush()
}

// Exists 检查键是否存在（未过期）
func (c *Cache[T]) Exists(key string) bool {
	_, found := c.inner.Get(key)
	return found
}

// Count 获取缓存条目数量
func (c *Cache[T]) Count() int {
	return c.inner.ItemCount()
}

// Keys 获取所有键（仅用于调试）
func (c *Cache[T]) Keys() []string {
	items := c.inner.Items()
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	return keys
}

// Items 获取所有缓存项（包含过期时间）
func (c *Cache[T]) Items() map[string]cache.Item {
	return c.inner.Items()
}

// 使用示例
// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {
// 	// 创建缓存
// 	cache := New[string](5*time.Minute, 1*time.Minute)

// 	// 设置缓存，10秒后过期
// 	cache.SetWithTTL("key1", "Hello World", 10*time.Second)

// 	// 方法1: 获取值和过期时间点
// 	val, expiration, found := cache.GetWithExpiration("key1")
// 	if found {
// 		fmt.Printf("值: %s, 过期时间: %v\n", val, expiration)
// 	}

// 	// 方法2: 获取剩余过期时间（推荐）
// 	remaining, found := cache.GetRemainingTTL("key1")
// 	if found {
// 		fmt.Printf("剩余时间: %v\n", remaining)
// 	}

// 	// 方法3: 检查是否过期
// 	isExpired := cache.IsExpired("key1")
// 	fmt.Printf("是否过期: %v\n", isExpired)

// 	// 等待过期
// 	time.Sleep(11 * time.Second)

// 	// 再次检查
// 	remaining, found = cache.GetRemainingTTL("key1")
// 	fmt.Printf("过期后剩余时间: %v, 存在: %v\n", remaining, found)

// 	// 测试永不过期
// 	cache.SetForever("key2", "Forever Value")
// 	remaining, found = cache.GetRemainingTTL("key2")
// 	fmt.Printf("永久缓存: 剩余时间=%v, 存在=%v\n", remaining, found)

// 	isExpired = cache.IsExpired("key2")
// 	fmt.Printf("永久缓存是否过期: %v\n", isExpired)
// }
