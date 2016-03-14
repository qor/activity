# Activity

Activity is dependant on [QOR Admin](https://github.com/qor/admin). It provides QOR Admin with an activity tracking feature for any Resource.

Applying Activity to a Resource will add `Comment` and `Track` data/state changes within the admin interface.

[![GoDoc](https://godoc.org/github.com/qor/activity?status.svg)](https://godoc.org/github.com/qor/activity)

## Usage

```go
import "github.com/qor/service"

func main() {
  Admin := service.New(&qor.Config{DB: db})
  order := Admin.AddResource(&models.Order{})

  // Register Activity for Order resource
  activity.Register(order)
}
```

[Online Demo](http://demo.getqor.com/admin/orders)

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).
