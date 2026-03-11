package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	core "github.com/sphireinc/core/v1"
)

func main() {
	app, err := core.New()
	if err != nil {
		log.Fatal(err)
	}

	api := app.Router.Group("/api/v1")
	api.Get("/our-custom-route", handler(app))
	api.Get("/non-mantis-route", nonMantisHandler(app))

	go func() {
		if err := app.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		log.Fatal(err)
	}
}

func handler(app *core.Config) func(ctx *core.Context) error {
	return func(ctx *core.Context) error {
		body := core.Res{
			Body:       []byte(`{"user":{"uuid":3}}`),
			BodyString: "hello",
		}
		return core.HandleResponseJSON(ctx, body.Byte(), app.S.OK)
	}
}

func nonMantisHandler(app *core.Config) func(ctx *core.Context) error {
	return func(ctx *core.Context) error {
		return core.HandleResponseJSON(ctx, []byte(`{}`), app.S.OK)
	}
}
