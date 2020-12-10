package route

import (
	"firefly/app/web/controller"
	"firefly/bootstrap"
	"github.com/kataras/iris/v12/mvc"
)

func Configure(app *bootstrap.Bootstrapper) {
	index := mvc.New(app.Party("/"))
	index.Handle(new(controller.Index))
}
