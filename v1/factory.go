package core

import (
	"fmt"
	"time"

	"github.com/qiangxue/fasthttp-routing"
	mantisCache "github.com/sphireinc/mantis/cache"
	mantisLog "github.com/sphireinc/mantis/log"
	mantisUUID "github.com/sphireinc/mantis/uuid"
)

func New() (*Config, error) {
	env := currentEnvironment()
	configFile := env + ".json"

	locations := configLocations()
	resolvedConfigFile, err := findConfigFile(locations, configFile)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Environment: Environment{
			Environment: env,
			Location:    resolvedConfigFile,
		},
	}
	cfg.Router.app = cfg

	if err := cfg.Load(); err != nil {
		return nil, err
	}

	if err := cfg.Factory(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) igniteComponent(name string, enabled bool, component interface{}) error {
	if !enabled {
		return nil
	}

	c.logInfo(fmt.Sprintf("igniting %s", name))

	switch comp := component.(type) {
	case Initializer:
		if err := comp.Init(); err != nil {
			return fmt.Errorf("initialize %s: %w", name, err)
		}
	case Connector:
		if err := comp.Connect(); err != nil {
			return fmt.Errorf("initialize %s: %w", name, err)
		}
	default:
		return fmt.Errorf("initialize %s: unsupported component type %T", name, component)
	}

	return nil
}

func (c *Config) initMemoryStore() {
	capacity := c.Persistence.InMemory.Capacity
	if capacity <= 0 {
		capacity = 100000
	}

	ttl := c.Persistence.InMemory.TTL.Or(30 * time.Second)
	c.Persistence.Memory = mantisCache.NewMemoryCache(int64(capacity), ttl.String())
	c.Runtime.Memory = c.Persistence.Memory
}

// Factory creates our Application and instantiates our services
func (c *Config) Factory() error {
	c.Router.app = c

	if c.Application.UUID == "" {
		c.Application.UUID = (mantisUUID.New()).String()
	}
	c.Application.Version = ResolveVersion()
	c.Application.Runtime = time.Now().UTC().Format(time.RFC3339)

	c.Runtime.Logger = noopLogger{}

	if c.Components.Log {
		printToTerm := c.Log.Options.PrintToTerm
		overwrite := c.Log.Options.Overwrite

		if c.Log.Writer != nil {
			printToTerm = c.Log.Writer.PrintToTerm
			overwrite = c.Log.Writer.Overwrite
		}

		writer, err := mantisLog.New(c.Log.Location, printToTerm, overwrite)
		if err != nil {
			return fmt.Errorf("initialize log: %w", err)
		}

		c.Log.Writer = writer
		c.Runtime.Logger = writer
	}

	if err := c.igniteComponent("Redis", c.Components.Redis, c.Persistence.Redis); err != nil {
		return err
	}
	if err := c.igniteComponent("MySQL", c.Components.MySQL, c.Persistence.MySQL); err != nil {
		return err
	}
	if err := c.igniteComponent("BigCache", c.Components.BigCache, c.Persistence.BigCache); err != nil {
		return err
	}
	if err := c.igniteComponent("MemCache", c.Components.MemCache, c.Persistence.MemCache); err != nil {
		return err
	}

	c.initMemoryStore()

	c.Router.Router = routing.New()
	wireDefaultHandlers(c.Router.Router)

	c.Router.New("/status", Status, c.Methods(), nil)
	c.Router.New("/teapot", TeaPot, c.Methods(), nil)

	if err := c.validate(); err != nil {
		return fmt.Errorf("validate application: %w", err)
	}

	c.logInfo(fmt.Sprintf("initializing %s %s @ %s", c.Application.Name, c.Application.Version, c.Application.Runtime))
	c.logInfo("App ID: " + c.Application.UUID)

	return nil
}
