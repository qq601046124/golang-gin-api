package model

import (
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Cache 是本地缓存实例
type Cache struct {
	Self *bigcache.BigCache
}

// LocalCache 是个单例
var LocalCache *Cache

// Init 初始化本地缓存单例
func (c *Cache) Init() {

	LocalCache = &Cache{
		Self: getCache(),
	}
}

// GetCache 创建一个本地缓存
func getCache() *bigcache.BigCache {

	defaultExpiration, cleanupInterval := viper.GetInt("cache.default_expiration"), viper.GetInt("cache.cleanup_interval")
	config := bigcache.Config{
		Shards:             1024,                                           // 存储的条目数量，值必须是2的幂
		LifeWindow:         time.Duration(defaultExpiration) * time.Second, // 超时后条目被处理
		CleanWindow:        time.Duration(cleanupInterval) * time.Second,   //处理超时条目的时间范围
		MaxEntriesInWindow: 0,                                              // 在 Life Window 中的最大数量，
		MaxEntrySize:       0,                                              // 条目最大尺寸，以字节为单位
		HardMaxCacheSize:   0,                                              // 设置缓存最大值，以MB为单位，超过了不在分配内存。0表示无限制分配
	}
	var initErr error
	BigCache, initErr := bigcache.NewBigCache(config)
	if initErr != nil {
		logrus.Error(initErr)
		return nil
	}
	return BigCache
}
