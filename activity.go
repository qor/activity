package activity

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/media_library"
	"github.com/qor/qor"
	"github.com/qor/qor/admin"
	"github.com/qor/qor/audited"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/qor/validations"
)

type QorActivity struct {
	gorm.Model
	Action       string
	Content      string `sql:"size:5000"`
	Note         string `sql:"size:2000"`
	Type         string
	ResourceType string
	ResourceID   string
	CreatorName  string
	audited.AuditedModel
}

func Register(res *admin.Resource) {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		admin.RegisterViewPath(path.Join(gopath, "src/github.com/qor/activity/views"))
	}

	qorAdmin := res.GetAdmin()
	assetManager := qorAdmin.GetResource("AssetManager")
	if assetManager == nil {
		assetManager = qorAdmin.AddResource(&media_library.AssetManager{}, &admin.Config{Invisible: true})
	}
	activityResource := qorAdmin.AddResource(&QorActivity{}, &admin.Config{Invisible: true})
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
	activityResource.Meta(&admin.Meta{Name: "Content", Type: "rich_editor", Resource: assetManager})
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
	router.Post(fmt.Sprintf("/%v/%v/!%v", res.ToParam(), res.ParamIDName(), qorAdmin.GetResource("QorActivity").ToParam()), ctrl.CreateActivityHandler)
	router.Post(fmt.Sprintf("/%v/%v/!%v/:activity_id/edit", res.ToParam(), res.ParamIDName(), qorAdmin.GetResource("QorActivity").ToParam()), ctrl.UpdateActivityHandler)
}
