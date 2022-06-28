package main

import (
	routing "github.com/qiangxue/fasthttp-routing"
	core "github.com/sphireinc/core/src"
	mantisHttp "github.com/sphireinc/mantis/http"
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

func handler(ctx *routing.Context) error {
	body := mantisHttp.Response{
		Body: []byte(`{}`),
	}
	return core.HandleResponseJSON(ctx, body.Byte(), mantisHttp.StatusOK)
}

func nonMantisHandler(ctx *routing.Context) error {
	return core.HandleResponseJSON(ctx, []byte(`{}`), mantisHttp.StatusOK)
}
