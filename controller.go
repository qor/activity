package activity

import (
	"net/http"

	"github.com/qor/qor/admin"
	"github.com/qor/responder"
)

type controller struct {
	ActivityResource *admin.Resource
}

func (ctrl controller) CreateActivityHandler(context *admin.Context) {
	result, err := context.FindOne()
	activityResource := ctrl.ActivityResource
	newActivity := &QorActivity{}
	if err == nil {
		context.Result = result
		if context.AddError(activityResource.Decode(context.Context, newActivity)); !context.HasError() {
			newActivity.ResourceType = context.Resource.ToParam()
			newActivity.ResourceID = getPrimaryKey(context)
			newActivity.CreatorName = context.CurrentUser.DisplayName()
			context.AddError(activityResource.CallSave(newActivity, context.Context))
		}
	}
	context.AddError(err)

	redirect_to := context.Request.Referer()
	if context.HasError() {
		responder.With("html", func() {
			context.Flash(context.Error(), "error")
			http.Redirect(context.Writer, context.Request, redirect_to, http.StatusFound)
		}).With("json", func() {
			context.JSON("edit", map[string]interface{}{"errors": context.GetErrors()})
		}).Respond(context.Request)
	} else {
		responder.With("html", func() {
			context.Flash(string(context.Admin.T(context.Context, "activity.successfully_created", "Activity was successfully created")), "success")
			http.Redirect(context.Writer, context.Request, redirect_to, http.StatusFound)
		}).With("json", func() {
			context.Resource = activityResource
			context.JSON("show", newActivity)
		}).Respond(context.Request)
	}
}

func (ctrl controller) UpdateActivityHandler(context *admin.Context) {
	c := context.Admin.NewContext(context.Writer, context.Request)
	c.ResourceID = ctrl.ActivityResource.GetPrimaryValue(context.Request)
	c.Resource = ctrl.ActivityResource
	c.Searcher = &admin.Searcher{Context: c}
	result, err := c.FindOne()

	context.AddError(err)
	if !context.HasError() {
		if context.AddError(c.Resource.Decode(c.Context, result)); !context.HasError() {
			context.AddError(context.Resource.CallSave(result, c.Context))
		}
	}

	redirect_to := context.Request.Referer()
	if context.HasError() {
		context.Writer.WriteHeader(admin.HTTPUnprocessableEntity)
		responder.With("html", func() {
			http.Redirect(context.Writer, context.Request, redirect_to, http.StatusFound)
		}).With("json", func() {
			context.JSON("edit", map[string]interface{}{"errors": context.GetErrors()})
		}).Respond(context.Request)
	} else {
		responder.With("html", func() {
			context.Flash(string(context.Admin.T(context.Context, "activity.successfully_updated", "Activity was successfully updated")), "success")
			http.Redirect(context.Writer, context.Request, redirect_to, http.StatusFound)
		}).With("json", func() {
			c.JSON("edit", result)
		}).Respond(context.Request)
	}
}
