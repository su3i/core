package datasource

import (
	"gorm.io/gorm"
)

type Datasource struct {
	gorm.Model

	SourceType      string `gorm:"not null"`
	SourceID      string `gorm:"not null"`
	ProjectID		uint   `gorm:"not null;index"` // <- foreign key to Project
	CreatedBy       map[string]string 		 `gorm:"type:jsonb;serializer:json;default:'{}'"`
}