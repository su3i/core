package etl

import (
	"log"
	"sync"

	"github.com/darksuei/suei-intelligence/internal/config"
	domain "github.com/darksuei/suei-intelligence/internal/domain/etl"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/etl/airbyte"
	"github.com/kelseyhightower/envconfig"
)

var (
	instance domain.ETL
	once     sync.Once
)

// GetInstance returns a singleton ETL instance
func GetInstance() domain.ETL {
	var cfg config.AirbyteConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to load cache config: %v", err)
	}

	return airbyte.Initialize(&cfg)
}