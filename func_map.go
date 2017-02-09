package activity

import (
	"fmt"
	"html/template"

	"github.com/qor/admin"
)


func activityJoinURL(path string, param string) string {
	return fmt.Sprintf("%s/!%s", path, param)
}

func registerFuncMap(a *admin.Admin) {
	funcMaps := template.FuncMap{
		"activity_join_url": activityJoinURL,
	}

	for key, value := range funcMaps {
		a.RegisterFuncMap(key, value)
	}
}
