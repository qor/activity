package activity

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/audited"
	"github.com/qor/media_library"
	"github.com/qor/qor"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/validations"
)

// QorActivity default model used to save resource's activities
type QorActivity struct {
	gorm.Model
	Action       string
	Content      string `gorm:"type:text"`
	Note         string `gorm:"type:text"`
	Type         string
	ResourceType string
	ResourceID   string
	CreatorName  string
	audited.AuditedModel
}

// Register register activity feature for an qor resource
func Register(res *admin.Resource) {
	var (
		qorAdmin         = res.GetAdmin()
		activityResource = qorAdmin.GetResource("QorActivity")
	)

	if activityResource == nil {
		// Auto run migration before add resource
		res.GetAdmin().Config.DB.AutoMigrate(&QorActivity{})

		activityResource = qorAdmin.AddResource(&QorActivity{}, &admin.Config{Invisible: true})
		activityResource.Meta(&admin.Meta{Name: "Action", Type: "hidden", Valuer: func(value interface{}, ctx *qor.Context) interface{} {
			act := value.(*QorActivity)
			if act.Action == "" {
				act.Action = "comment on"
			}
			return activityResource.GetAdmin().T(ctx, "activity."+act.Action, act.Action)
		}})
		activityResource.Meta(&admin.Meta{Name: "UpdatedAt", Type: "hidden", Valuer: func(value interface{}, ctx *qor.Context) interface{} {
			return value.(*QorActivity).UpdatedAt.Format("Jan 2 15:04")
		}})
		activityResource.Meta(&admin.Meta{Name: "URL", Valuer: func(value interface{}, ctx *qor.Context) interface{} {
			return fmt.Sprintf("/admin/%v/%v/!%v/%v/edit", res.ToParam(), res.GetPrimaryValue(ctx.Request), activityResource.ToParam(), value.(*QorActivity).ID)
		}})

		assetManager := qorAdmin.GetResource("AssetManager")
		if assetManager == nil {
			assetManager = qorAdmin.AddResource(&media_library.AssetManager{}, &admin.Config{Invisible: true})
		}

		activityResource.Meta(&admin.Meta{Name: "Content", Type: "rich_editor", Resource: assetManager})
		activityResource.Meta(&admin.Meta{Name: "Note", Type: "string", Resource: assetManager})
		activityResource.EditAttrs("Action", "Content", "Note")
		activityResource.ShowAttrs("ID", "Action", "Content", "Note", "URL", "UpdatedAt", "CreatorName")
		activityResource.AddValidator(func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
			if meta := metaValues.Get("Content"); meta != nil {
				if name := utils.ToString(meta.Value); strings.TrimSpace(name) == "" {
					return validations.NewError(record, "Content", "Content can't be blank")
				}
			}
			return nil
		})
	}

	admin.RegisterViewPath("github.com/qor/activity/views")
	res.UseTheme("activity")

	qorAdmin.RegisterFuncMap("get_activities", func(context *admin.Context, types ...string) []QorActivity {
		activities, _ := GetActivities(context, types...)
		return activities
	})

	qorAdmin.RegisterFuncMap("formatted_datetime", func(datetime time.Time) string {
		return datetime.Format("Jan 2 15:04")
	})

	qorAdmin.RegisterFuncMap("formatted_content", func(content string) template.HTML {
		return template.HTML(content)
	})

	qorAdmin.RegisterFuncMap("activity_resource", func() *admin.Resource {
		return qorAdmin.GetResource("QorActivity")
	})

	router := res.GetAdmin().GetRouter()
	ctrl := controller{ActivityResource: activityResource}
	router.Get(fmt.Sprintf("/%v/%v/!qor_activities", res.ToParam(), res.ParamIDName()), ctrl.GetActivityHandler)
	router.Post(fmt.Sprintf("/%v/%v/!qor_activities", res.ToParam(), res.ParamIDName()), ctrl.CreateActivityHandler)
	router.Post(fmt.Sprintf("/%v/%v/!qor_activities/%v/edit", res.ToParam(), res.ParamIDName(), activityResource.ParamIDName()), ctrl.UpdateActivityHandler)
}
