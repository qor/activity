package activity

import (
	"net/http"

	"github.com/qor/admin"
	"github.com/qor/responder"
)

type controller struct {
	ActivityResource *admin.Resource
}

func (ctrl controller) GetActivity(context *admin.Context) {
	activities, _ := GetActivities(context, "-tag")
	activityResource := ctrl.ActivityResource

	if context.HasError() {
		responder.With("json", func() {
			context.JSON("edit", map[string]interface{}{"errors": context.GetErrors()})
		}).Respond(context.Request)
	} else {
		responder.With("json", func() {
			context.Resource = activityResource
			context.JSON("index", activities)
		}).Respond(context.Request)
	}
}

func (ctrl controller) CreateActivity(context *admin.Context) {
	result, err := context.FindOne()
	activityResource := ctrl.ActivityResource
	newActivity := &QorActivity{}
	if err == nil {
		if context.AddError(activityResource.Decode(context.Context, newActivity)); !context.HasError() {
			context.AddError(CreateActivity(context, newActivity, result))
		}
	}
	context.AddError(err)

	redirectTo := context.Request.Referer()
	if context.HasError() {
		responder.With("html", func() {
			context.Flash(context.Error(), "error")
			http.Redirect(context.Writer, context.Request, redirectTo, http.StatusFound)
		}).With("json", func() {
			context.JSON("edit", map[string]interface{}{"errors": context.GetErrors()})
		}).Respond(context.Request)
	} else {
		responder.With("html", func() {
			context.Flash(string(context.Admin.T(context.Context, "activity.successfully_created", "Activity was successfully created")), "success")
			http.Redirect(context.Writer, context.Request, redirectTo, http.StatusFound)
		}).With("json", func() {
			context.Resource = activityResource
			context.JSON("show", newActivity)
		}).Respond(context.Request)
	}
}

func (ctrl controller) UpdateActivity(context *admin.Context) {
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

	redirectTo := context.Request.Referer()
	if context.HasError() {
		context.Writer.WriteHeader(admin.HTTPUnprocessableEntity)
		responder.With("html", func() {
			http.Redirect(context.Writer, context.Request, redirectTo, http.StatusFound)
		}).With("json", func() {
			context.JSON("edit", map[string]interface{}{"errors": context.GetErrors()})
		}).Respond(context.Request)
	} else {
		responder.With("html", func() {
			context.Flash(string(context.Admin.T(context.Context, "activity.successfully_updated", "Activity was successfully updated")), "success")
			http.Redirect(context.Writer, context.Request, redirectTo, http.StatusFound)
		}).With("json", func() {
			c.JSON("show", result)
		}).Respond(context.Request)
	}
}
