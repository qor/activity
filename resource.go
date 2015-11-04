package activity

import (
	"fmt"
	"strings"

	"github.com/qor/qor"
	"github.com/qor/qor/admin"
)

type Resource struct {
	*admin.Resource
	Context *qor.Context
}

func New(res *admin.Resource, context *qor.Context) *Resource {
	return &Resource{Resource: res, Context: context}
}

func (res *Resource) getPrimaryKey(record interface{}) string {
	db := res.Context.GetDB()

	var primaryValues []string
	for _, field := range db.NewScope(record).PrimaryFields() {
		primaryValues = append(primaryValues, fmt.Sprint(field.Field.Interface()))
	}
	return strings.Join(primaryValues, "::")
}

func (res *Resource) CreateActivity(record interface{}, activity Activity) error {
	db := res.Context.GetDB()
	activity.ResourceType = res.ToParam()
	activity.ResourceID = res.getPrimaryKey(record)
	return db.Save(&activity).Error
}

func (res *Resource) GetActivities(record interface{}, types ...string) ([]Activity, error) {
	var activities []Activity
	db := res.Context.GetDB().Where("resource_id = ? AND resource_type = ?", res.getPrimaryKey(record), res.ToParam())
	if len(types) > 0 {
		db = db.Where("type IN (?)", types)
	}
	err := db.Find(&activities).Error
	return activities, err
}
