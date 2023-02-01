package main

import (
	"github.com/sphireinc/core/v1"
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

func handler(ctx *core.Context) error {
	output := `{
		"user": {
			"uuid": 3
		}
	}`
	body := core.Res{
		Body:       []byte(output),
		BodyString: "hello",
	}
	return core.HandleResponseJSON(ctx, body.Byte(), App.S.OK)
}

func nonMantisHandler(ctx *core.Context) error {
	return core.HandleResponseJSON(ctx, []byte(`{}`), App.S.OK)
}
