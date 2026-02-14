package datasource

type DatasourceRepository interface {
	Find(projectId uint) (*[]Datasource, error)
	FindOne(datasourceId uint, projectId uint) (*Datasource, error)
	Create(payload *Datasource) (*Datasource, error)
	Update(payload *Datasource) error
	SoftDelete(datasourceId uint, projectId uint) error
	HardDelete(datasourceId uint, projectId uint) error
}