package config

import (
	domain "github.com/darksuei/suei-intelligence/internal/domain/cache"
)

type CacheConfig struct {
	CacheType        domain.CacheType `default:"memory"`
	RedisAddr  string             `required:"false"`
	RedisPassword    string             `required:"false"`
	RedisDB int             `required:"false"`
}
