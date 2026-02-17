package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/darksuei/suei-intelligence/internal/domain/datasource"
)

type datasourceRepository struct {
	db *gorm.DB
}

func (r *datasourceRepository) FindOne(datasourceID uint, projectId uint) (*datasource.Datasource, error) {
	var _datasource datasource.Datasource

	if err := r.db.Where(&datasource.Datasource{ProjectID: projectId, Model: gorm.Model{ID: datasourceID}}).First(&_datasource).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &_datasource, nil
}

func (r *datasourceRepository) Find(projectId uint) (*[]datasource.Datasource, error) {
	var _datasources []datasource.Datasource

	if err := r.db.Where(&datasource.Datasource{ProjectID: projectId}).Find(&_datasources).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &_datasources, nil
}

func (r *datasourceRepository) Create(payload *datasource.Datasource) (*datasource.Datasource, error) {
	_datasource := datasource.Datasource{
		SourceType: payload.SourceType,
		SourceID: payload.SourceID,
		ProjectID: payload.ProjectID,
		CreatedBy: payload.CreatedBy,
	}

	err := r.db.Create(&_datasource).Error

	if err != nil {
		return nil, errors.New("failed to create datasource: " + err.Error())
	}

	return &_datasource, nil
}

func (r *datasourceRepository) Update(payload *datasource.Datasource) error {
	err := r.db.Updates(payload).Error

	if err != nil {
		return errors.New("failed to update datasource: " + err.Error())
	}

	return nil
}

func (r *datasourceRepository) SoftDelete(datasourceID, projectID uint) error {
	return r.db.
		Where(&datasource.Datasource{
			Model: gorm.Model{ID: datasourceID},
			ProjectID: projectID,
		}).
		Delete(&datasource.Datasource{}).
		Error
}

func (r *datasourceRepository) HardDelete(datasourceID, projectID uint) error {
	return r.db.Unscoped().
		Where(&datasource.Datasource{
			Model: gorm.Model{ID: datasourceID},
			ProjectID: projectID,
		}).
		Delete(&datasource.Datasource{}).
		Error
}

func NewDatasourceRepository(db *gorm.DB) datasource.DatasourceRepository {
	return &datasourceRepository{db: db}
}