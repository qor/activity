package activity

import "github.com/qor/qor/admin"

func CreateActivity(context *admin.Context) {
	if result, err := context.FindOne(); err == nil {
		resource := New(context.Resource, result, context.Context)
		var activity = QorActivity{
			Type:    context.Request.Form.Get("type"),
			Comment: context.Request.Form.Get("comment"),
		}
		resource.CreateActivity(activity)
	}
}
