package activity

import "github.com/qor/qor/admin"

type Resource struct {
	*admin.Resource
}

func New(res *admin.Resource) *Resource {
	return &Resource{res}
}
