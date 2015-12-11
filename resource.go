package activity

import (
	"fmt"
	"strings"

	"github.com/qor/qor/admin"
)

func getPrimaryKey(context *admin.Context) string {
	db := context.GetDB()

	var primaryValues []string
	result, _ := context.FindOne()
	for _, field := range db.NewScope(result).PrimaryFields() {
		primaryValues = append(primaryValues, fmt.Sprint(field.Field.Interface()))
	}
	return strings.Join(primaryValues, "::")
}

func CreateActivity(activity QorActivity, context *admin.Context) error {
	db := context.GetDB()
	activity.ResourceType = context.Resource.ToParam()
	activity.ResourceID = getPrimaryKey(context)
	return db.Save(&activity).Error
}

func GetActivities(context *admin.Context, types ...string) ([]QorActivity, error) {
	var activities []QorActivity
	db := context.GetDB().Where("resource_id = ? AND resource_type = ?", getPrimaryKey(context), context.Resource.ToParam())

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
