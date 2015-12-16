package activity

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
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

	res.UseTheme("activities")

	res.GetAdmin().RegisterFuncMap("get_activities", func(context *admin.Context, types ...string) []QorActivity {
		activities, _ := GetActivities(context, types...)
		return activities
	})

	res.GetAdmin().RegisterFuncMap("formatted_datetime", func(datetime time.Time) string {
		return datetime.Format("Jan 2 15:04")
	})

	router := res.GetAdmin().GetRouter()
	router.Post(fmt.Sprintf("/%v/(.*?)/!activity", res.ToParam()), CreateActivityHandler)
}
