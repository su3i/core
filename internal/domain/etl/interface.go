package etl

type EtlType string

const (
	EtlTypeRedis   EtlType = "airbyte"
)

// Cache defines the minimal cache operations
type ETL interface {
	CreateSourceConnection(name string, configuration interface {}) error
}