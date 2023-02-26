package core

import (
	"fmt"
	"github.com/qiangxue/fasthttp-routing"
	mantisCache "github.com/sphireinc/mantis/cache"
	mantisLog "github.com/sphireinc/mantis/log"
	mantisUUID "github.com/sphireinc/mantis/uuid"
	"os"
	"time"
)

var App *Config

// New creates a new App based on Config. It Loads() the config file based
// on the ENV that is set, then calls Factory() to stand up our application.
func New() *Config {
	env := os.Getenv(EnvironmentVariableName)
	if env == "" {
		env = "dev"
	}

	App = &Config{
		Environment: Environment{
			Environment: env,
			Location:    env + ".json",
		},
	}

	App.Load()
	App.Factory()

	return App
}

// Factory creates our Application and instantiates our services
func (c *Config) Factory() {
	if c.Application.UUID == "" {
		c.Application.UUID = (mantisUUID.New()).String()
	}
	c.Application.Version = Version
	c.Application.Runtime = time.Now().UTC().Format(time.RFC3339)

	if c.Components.Log {
		c.Log.Writer, _ = mantisLog.New(c.Log.Location, c.Log.Writer.PrintToTerm, c.Log.Writer.Overwrite)
	}

	ignite := func(name string, componentStatus bool, init func() error) {
		if componentStatus {
			c.Log.Writer.Info(fmt.Sprintf("igniting %s", name))
			if err := init(); err != nil {
				App.Log.Writer.Fatal(err.Error())
			}
		}
	}

	ignite("Redis", c.Components.Redis, c.Persistence.Redis.Init)
	ignite("MySQL", c.Components.MySQL, c.Persistence.MySQL.Connect)
	ignite("BigCache", c.Components.BigCache, c.Persistence.BigCache.Init)
	ignite("MemCache", c.Components.MemCache, c.Persistence.MemCache.Init)

	c.Persistence.Memory = mantisCache.NewMemoryCache(100000, "30s")

	// Initiate our Router
	c.Router.Router = routing.New()
	c.Router.New("/status", Status, c.Methods(), nil)
	c.Router.New("/teapot", TeaPot, c.Methods(), nil)

	if err := c.validate(); err != nil {
		c.Log.Writer.Fatal("error validating Core application", err.Error())
		panic(err.Error())
	}

	c.Log.Writer.Info(fmt.Sprintf("initializing %s %s @ %s", c.Application.Name, c.Application.Version, c.Application.Runtime))
	c.Log.Writer.Info("App ID: " + c.Application.UUID)
}
