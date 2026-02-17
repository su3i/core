package datasource

import (
	"errors"

	"github.com/darksuei/suei-intelligence/internal/application/account"
	"github.com/darksuei/suei-intelligence/internal/application/project"
	"github.com/darksuei/suei-intelligence/internal/config"
	"github.com/darksuei/suei-intelligence/internal/domain/datasource"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/database"
)

func NewDatasource(key string, sourceType string, sourceId string, createdByEmail string, cfg *config.DatabaseConfig) (*datasource.Datasource, error) {
	_datasourceRepository := database.NewDatasourceRepository(cfg)

	createdByAccount, err := account.RetrieveAccount(createdByEmail, cfg)

	if err != nil {
		return nil, errors.New("Failed to get account")
	}

	createdBy := map[string]string{
		"Email": createdByEmail,
		"Name": createdByAccount.Name,
	}

	_project, err := project.RetrieveProject(key, cfg)

	if err != nil {
		return nil, errors.New("Invalid project key")
	}

	_datasource := &datasource.Datasource{
		SourceType: sourceType,
		SourceID: sourceId,
		ProjectID: _project.ID,
		CreatedBy: createdBy,
	}

	return _datasourceRepository.Create(_datasource)
}

func RetrieveDatasource(datasourceID uint, key string, cfg *config.DatabaseConfig) (*datasource.Datasource, error) {
	_datasourceRepository := database.NewDatasourceRepository(cfg)

	_project, err := project.RetrieveProject(key, cfg)

	if err != nil {
		return nil, errors.New("Invalid project key")
	}

	return _datasourceRepository.FindOne(datasourceID, _project.ID)
}

func RetrieveDatasources(key string, cfg *config.DatabaseConfig) (*[]datasource.Datasource, error) {
	_datasourceRepository := database.NewDatasourceRepository(cfg)

	_project, err := project.RetrieveProject(key, cfg)

	if err != nil {
		return nil, errors.New("Invalid project key")
	}

	return _datasourceRepository.Find(_project.ID)
}

func SoftDeleteDatasource(datasourceID uint, key string, cfg *config.DatabaseConfig) error {
	_datasourceRepository := database.NewDatasourceRepository(cfg)

	_project, err := project.RetrieveProject(key, cfg)

	if err != nil {
		return errors.New("Invalid project key")
	}

	return _datasourceRepository.SoftDelete(datasourceID, _project.ID)
}

func HardDeleteDatasource(datasourceID uint, key string, cfg *config.DatabaseConfig) error {
	_datasourceRepository := database.NewDatasourceRepository(cfg)

	_project, err := project.RetrieveProject(key, cfg)

	if err != nil {
		return errors.New("Invalid project key")
	}

	return _datasourceRepository.HardDelete(datasourceID, _project.ID)
}