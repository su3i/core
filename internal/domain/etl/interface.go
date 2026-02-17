package etl

type EtlType string

const (
	EtlTypeRedis   EtlType = "airbyte"
)

// Minimal ETL operations
type ETL interface {
	CreateSourceConnection(name string, configuration map[string]interface{}) (*string, error)
	DeleteSourceConnection(sourceId string) error
	TestSourceConnection(sourceId string) error
}