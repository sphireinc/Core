# Sphire Core

Core is a Go framework for fast API prototyping built on:
- `fasthttp` (server)
- `github.com/qiangxue/fasthttp-routing` (router)

<div align="center">
  <img src="assets/core_logo.png" width="400px" alt="logo" />
</div>

## Setup

Core is configured via a JSON file loaded into the `Config` struct.

### Environment variables

- `SPC_ENV` (preferred): selects the config file `<env>.json` (default: `dev`)
- `SPC_CONFIG` (preferred): absolute path to a config file (overrides `SPC_ENV`)

Backward compatible fallbacks are still accepted:
- `SPK_ENV`
- `SPK_CONFIG`

### Minimal start

```go
package main

import (
  "log"
  core "github.com/sphireinc/core/v1"
)

func main() {
  app, err := core.New()
  if err != nil {
    log.Fatal(err)
  }

  if err := app.Run(); err != nil {
    log.Fatal(err)
  }
}
```

Example app code lives under `cmd/example/`

### Routing

Handlers use this signature:

```go
func handler(ctx *core.Context) error {
  return core.HandleResponseJSON(ctx, []byte(`{}`), 200)
}
```

Register routes before Run():

```go
app.Router.Get("/our-custom-route", handler)
```

### Profiling / Stats View

When enabled, Core uses github.com/go-echarts/statsview.

### Default endpoint:

- :18066/debug/statsview
