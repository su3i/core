package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/darksuei/suei-intelligence/internal/domain/metadata"
)

var DB *gorm.DB

type metadataRepository struct{}

func (r *metadataRepository) FindOne() (*metadata.Metadata, error) {
	var _metadata metadata.Metadata

	if err := DB.Unscoped().First(&_metadata).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &_metadata, nil
}

func (r *metadataRepository) Create(payload *metadata.Metadata) (*metadata.Metadata, error) {
	_metadata := metadata.Metadata{BootstrapToken: payload.BootstrapToken}

	err := DB.Create(&_metadata).Error

	if err != nil {
		return nil, errors.New("failed to create metadata: " + err.Error())
	}

	return &_metadata, nil
}

func (r *metadataRepository) Update(payload *metadata.Metadata) error {
	err := DB.Updates(payload).Error

	if err != nil {
		return errors.New("failed to update metadata: " + err.Error())
	}

	return nil
}

func NewMetadataRepository(_db *gorm.DB) metadata.MetadataRepository {
	DB = _db
	return &metadataRepository{}
}
