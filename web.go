package main

import (
	"firefly/bootstrap"
	"firefly/route"
)

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("Firefly", "grail")
	app.Bootstrap()
	app.Configure(route.Configure)
	return app
}

func main() {
	app := newApp()
	app.Listen(":9464")
}
