package activity

import (
	"github.com/qor/qor/admin"
	"github.com/qor/responder"
	"net/http"
	"regexp"
)

func CreateActivityHandler(context *admin.Context) {
	result, err := context.FindOne()
	activityResource := context.GetResource("QorActivity")
	if err == nil {
		context.Result = result
		newActivity := &QorActivity{}
		if context.AddError(activityResource.Decode(context.Context, newActivity)); !context.HasError() {
			newActivity.ResourceType = context.Resource.ToParam()
			newActivity.ResourceID = getPrimaryKey(context)
			newActivity.CreatorName = context.CurrentUser.DisplayName()
			context.AddError(activityResource.CallSave(newActivity, context.Context))
		}
	}
	context.AddError(err)

	if context.HasError() {
		responder.With("html", func() {
			context.Flash(context.Error(), "error")
			http.Redirect(context.Writer, context.Request, context.Request.PostFormValue("redirect_to"), http.StatusFound)
		}).With("json", func() {
			context.JSON("edit", map[string]interface{}{"errors": context.GetErrors()})
		}).Respond(context.Request)
	} else {
		responder.With("html", func() {
			context.Flash(string(context.Admin.T(context.Context, "resource_successfully_created", "Activity was successfully created")), "success")
			http.Redirect(context.Writer, context.Request, context.Request.PostFormValue("redirect_to"), http.StatusFound)
		}).With("json", func() {
			context.JSON("edit", result)
		}).Respond(context.Request)
	}
}

func UpdateActivityHandler(context *admin.Context) {
	c := context.Admin.NewContext(context.Writer, context.Request)
	c.Resource = context.GetResource("QorActivity")
	re := regexp.MustCompile("/(\\d?)/edit$")
	activityID := re.FindString(context.Request.URL.Path)
	c.ResourceID = activityID
	result, err := c.FindOne()

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
		}).Respond(context.Request)
	} else {
		responder.With("html", func() {
			context.Flash(string(context.Admin.T(context.Context, "resource_successfully_created", "Activity was successfully created")), "success")
			http.Redirect(context.Writer, context.Request, context.Request.PostFormValue("redirect_to"), http.StatusFound)
		}).With("json", func() {
			context.JSON("edit", result)
		}).Respond(context.Request)
	}
}
