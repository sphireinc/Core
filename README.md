# Sphire Core
Core is a framework which allows for fast API prototyping. Based off of Fiber.

## Setup

The (pun fully intended) *core* of Sphire Core is the `Config` struct. This is the entrypoint
of the framework, and must be instantiated as follows:

    var App core.Config = core.New()

Within `core.New()` (located in `factory.go`), we look at the `SPC_ENV` environment variable
which must coincide with the configuration JSON filename. If `SPC_ENV` is set to `dev` for instance
then the configuration JSON filename must be `dev.json`. By default, it uses `dev.json`. 

Within `New()`, the `Load()` function is called, which loads our JSON configuration into
the `Config` struct. It then proceeds to call `Factory()` which instantiates our enabled services.

It is then possible to simply run `App.Run()` to start listening on our designated address:port.

The framework can run in the simple runner:

    package main
    
    import (
        core "github.com/sphireinc/core/v1"
    )

    var App := core.New()

    func main() {
        App.Run()
    }

That is all that is required.

## Routing

Routing in Core is simple. We create a function with this signature:

    func handler(ctx *routing.Context) error {
        body := core.Res{}
        return core.HandleResponseJSON(ctx, body.Byte(), App.S.OK)
    }

Then, we add it to our `app` before calling `app.Run()`:

    package main
    
    import (
        core "github.com/sphireinc/core/v1"
    )
    
    var App = core.New()
    
    func main() {
        App.Router.Get("/our-custom-route", handler)
        App.Router.Get("/non-mantis-route", nonMantisHandler)
        App.Run()
    }

    func handler(ctx *core.Context) error {
        body := core.Res{
            Body:       []byte(`{"x": 3}`),
            BodyString: string("hello"),
        }
        return core.HandleResponseJSON(ctx, body.Byte(), App.S.OK)
    }

    func nonMantisHandler(ctx *core.Context) error {
        return core.HandleResponseJSON(ctx, []byte(`{}`), App.S.OK)
    }

That is all there is to it. Things like MySQL and Redis are set up automatically
when they find a configuration, and hang off of the `App` struct (like the Mantis logger).


## Profiling / Stats View

Core utilizes the StatsView package (github.com/go-echarts/statsview) to display a quick
statistical view of memory usage etc. This can be viewed at `:18066/debug/statsview`


## To Do

1. Implement FastWS (https://github.com/fasthttp/fastws)
2. Implement Proper CORS, CSRF, Limiter, etc (https://github.com/gofiber/fiber/tree/master/middleware)
3. Implement Goth (https://github.com/markbates/goth)