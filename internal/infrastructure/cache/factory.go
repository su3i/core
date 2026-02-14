package cache

import (
	"log"
	"sync"

	"github.com/darksuei/suei-intelligence/internal/config"
	domain "github.com/darksuei/suei-intelligence/internal/domain/cache"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/cache/memory"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/cache/redis"
	"github.com/kelseyhightower/envconfig"
)

var (
	instance domain.Cache
	once     sync.Once
)

// GetCache returns a singleton cache instance
func GetCache() domain.Cache {
	var cfg config.CacheConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to load cache config: %v", err)
	}

	once.Do(func() {
		switch cfg.CacheType {
			case domain.CacheTypeRedis:
				instance = redis.NewCache(&cfg)
			case domain.CacheTypeMemory:
				instance = memory.NewCache()
			default:
				instance = memory.NewCache()
		}
	})
	return instance
}