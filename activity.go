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
)

type QorActivity struct {
	gorm.Model
	Action       string
	Content      string
	Note         string
	Type         string
	ResourceType string
	ResourceID   string
	audited.AuditedModel
}

func RegisterActivityMeta(res *admin.Resource) {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		admin.RegisterViewPath(path.Join(gopath, "src/github.com/qor/activity/views"))
	}

	qorAdmin := res.GetAdmin()
	resource := qorAdmin.GetResource("QorActivity")
	if resource == nil {
		assetManager := qorAdmin.AddResource(&media_library.AssetManager{}, &admin.Config{Invisible: true})
		activity := qorAdmin.AddResource(QorActivity{}, &admin.Config{Invisible: true})
		activity.Meta(&admin.Meta{Name: "Action", Type: "hidden", Valuer: func(value interface{}, ctx *qor.Context) interface{} {
			return "comment"
		}})
		activity.Meta(&admin.Meta{Name: "Content", Type: "rich_editor", Resource: assetManager})
		activity.EditAttrs("Action", "Content", "Note")
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

	router := res.GetAdmin().GetRouter()
	router.Post(fmt.Sprintf("/%v/(.*?)/!activity", res.ToParam()), CreateActivityHandler)
}
