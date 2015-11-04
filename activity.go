package activity

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/qor/admin"
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

func RegisterActivityMeta(res *admin.Resource) {
	res.Meta(&admin.Meta{
		Name: "Activities",
		Type: "activities",
	})
}
