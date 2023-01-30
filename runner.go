package main

import (
	"github.com/sphireinc/core/v1"
	mantis "github.com/sphireinc/mantis/http"
)

var App = core.New()

/**
 * This is purely an example application to aid in the development of this framework
 */

func main() {
	App.Router.Get("/our-custom-route", handler)
	App.Router.Get("/non-mantis-route", nonMantisHandler)
	App.Run()
}

func handler(ctx App.Context) error {
	body := mantis.Response{
		Body: []byte(`{}`),
	}
	return core.HandleResponseJSON(ctx, body.Byte(), mantis.StatusOK)
}

func nonMantisHandler(ctx App.Context) error {
	return core.HandleResponseJSON(ctx, []byte(`{}`), mantis.StatusOK)
}
