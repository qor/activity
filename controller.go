package activity

import (
	"github.com/qor/qor/admin"
	"github.com/qor/qor/responder"
	"net/http"
)

func CreateActivityHandler(context *admin.Context) {
	var activity = QorActivity{
		Type:    context.Request.Form.Get("QorResource.Type"),
		Action:  context.Request.Form.Get("QorResource.Action"),
		Content: context.Request.Form.Get("QorResource.Content"),
		Note:    context.Request.Form.Get("QorResource.Note"),
	}
	userName := context.CurrentUser.DisplayName()
	activity.SetCreatedBy(userName)
	activity.SetUpdatedBy(userName)
	result, err := context.FindOne()
	if err == nil {
		context.Result = result
		err = CreateActivity(activity, context)
	}
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

func UpdateActivityHandler(context *admin.Context) {
	result, err := context.FindOne()
	context.AddError(err)
	if !context.HasError() {
		if context.AddError(context.Resource.Decode(context.Context, result)); !context.HasError() {
			context.AddError(context.Resource.CallSave(result, context.Context))
		}
	}

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
