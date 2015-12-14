package activity

import (
	"github.com/qor/qor/admin"
	"github.com/qor/qor/responder"
	"net/http"
)

func CreateActivityHandler(context *admin.Context) {
	var activity = QorActivity{
		Type:    context.Request.Form.Get("type"),
		Comment: context.Request.Form.Get("comment"),
	}
	result, _ := context.FindOne()
	err := CreateActivity(activity, context)
	context.AddError(err)

	if context.HasError() {
		context.Writer.WriteHeader(admin.HTTPUnprocessableEntity)
		responder.With("html", func() {
			http.Redirect(context.Writer, context.Request, context.Request.PostFormValue("redirect_to"), http.StatusFound)
		}).With("json", func() {
			context.JSON("edit", map[string]interface{}{"errors": context.GetErrors()})
		}).Respond(context.Writer, context.Request)
	} else {
		responder.With("html", func() {
			context.Flash(string(context.T("resource_successfully_created", "Activity was successfully created")), "success")
			http.Redirect(context.Writer, context.Request, context.Request.PostFormValue("redirect_to"), http.StatusFound)
		}).With("json", func() {
			context.JSON("edit", result)
		}).Respond(context.Writer, context.Request)
	}
}
