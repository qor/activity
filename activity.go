package activity

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/qor/audited"
)

type Activity struct {
	gorm.Model
	Comment      string
	Type         string
	ResourceType string
	ResourceID   string
	audited.AuditedModel
}
