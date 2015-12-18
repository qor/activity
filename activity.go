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
	Content      string
	Note         string
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
	if qorAdmin.GetResource("QorActivity") == nil {
		assetManager := qorAdmin.AddResource(&media_library.AssetManager{}, &admin.Config{Invisible: true})
		activity := qorAdmin.AddResource(QorActivity{}, &admin.Config{Invisible: true})
		activity.Meta(&admin.Meta{Name: "Action", Type: "hidden", Valuer: func(value interface{}, ctx *qor.Context) interface{} {
			return "comment on"
		}})
		activity.Meta(&admin.Meta{Name: "Content", Type: "rich_editor", Resource: assetManager})
		activity.EditAttrs("Action", "Content", "Note")
		activity.AddValidator(func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
			if meta := metaValues.Get("Content"); meta != nil {
				if name := utils.ToString(meta.Value); strings.TrimSpace(name) == "" {
					return validations.NewError(record, "Content", "Content can't be blank")
				}
			}
			return nil
		})
	}

	res.UseTheme("activities")

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

	qorAdmin.RegisterFuncMap("new_activity", func() QorActivity {
		return QorActivity{}
	})

	qorAdmin.RegisterFuncMap("is_edit_or_show", func(context *admin.Context) bool {
		return context.Action == "edit" || context.Action == "show"
	})

	router := res.GetAdmin().GetRouter()
	router.Post(fmt.Sprintf("/%v/(.*?)/!%v", res.ToParam(), qorAdmin.GetResource("QorActivity").ToParam()), CreateActivityHandler)
	router.Post(fmt.Sprintf("/%v/(.*?)/!%v/(.*?)/edit", res.ToParam(), qorAdmin.GetResource("QorActivity").ToParam()), UpdateActivityHandler)
}
