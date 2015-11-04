package activity

import (
	"fmt"
	"strings"

	"github.com/qor/qor"
	"github.com/qor/qor/admin"
)

type Resource struct {
	*admin.Resource
	Record  interface{}
	Context *qor.Context
}

func New(res *admin.Resource, record interface{}, context *qor.Context) *Resource {
	return &Resource{Resource: res, Record: record, Context: context}
}

func (res Resource) getPrimaryKey() string {
	db := res.Context.GetDB()

	var primaryValues []string
	for _, field := range db.NewScope(res.Record).PrimaryFields() {
		primaryValues = append(primaryValues, fmt.Sprint(field.Field.Interface()))
	}
	return strings.Join(primaryValues, "::")
}

func (res Resource) CreateActivity(activity QorActivity) error {
	db := res.Context.GetDB()
	activity.ResourceType = res.ToParam()
	activity.ResourceID = res.getPrimaryKey()
	return db.Save(&activity).Error
}

func (res Resource) GetActivities(types ...string) ([]QorActivity, error) {
	var activities []QorActivity
	db := res.Context.GetDB().Where("resource_id = ? AND resource_type = ?", res.getPrimaryKey(), res.ToParam())

	var inTypes, notInTypes []string
	for _, t := range types {
		if strings.HasPrefix(t, "-") {
			notInTypes = append(notInTypes, strings.TrimPrefix(t, "-"))
		} else {
			inTypes = append(inTypes, t)
		}
	}

	if len(inTypes) > 0 {
		db = db.Where("type IN (?)", inTypes)
	}

	if len(notInTypes) > 0 {
		db = db.Where("type NOT IN (?)", notInTypes)
	}

	err := db.Find(&activities).Error
	return activities, err
}
