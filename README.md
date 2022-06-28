# Sphire Eroc
Eroc is a framework which allows for fast API prototyping. Based off of Fiber.

## Setup

The (pun fully intended) *core* of Sphire Eroc is the `Config` struct. This is the entrypoint
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
        core "github.com/sphireinc/core/src"
    )

    var App := core.New()

    func main() {
        App.Run()
    }

That is all that is required.

## Routing

Routing in Eroc is simple. We create a function with this signature:

    func handler(ctx *routing.Context) error {
        body := mantisHttp.Response{}
        return core.HandleResponseJSON(ctx, body.Byte(), mantisHttp.StatusOK)
    }

Then, we add it to our `app` before calling `app.Run()`:

    package main
    
    import (
        routing "github.com/qiangxue/fasthttp-routing"
        core "github.com/sphireinc/core/src"
        mantisHttp "github.com/sphireinc/mantis/http"
    )
    
    var App = core.New()
    
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

That is all there is to it. Things like MySQL and Redis are set up automatically
when they find a configuration, and hang off of the `App` struct (like the Mantis logger).


## To Do

1. Implement FastWS (https://github.com/fasthttp/fastws)
2. Implement Proper CORS, CSRF, Limiter, etc (https://github.com/gofiber/fiber/tree/master/middleware)
3. Implement Goth (https://github.com/markbates/goth)