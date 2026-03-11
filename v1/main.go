package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/go-echarts/statsview"
	routing "github.com/qiangxue/fasthttp-routing"
	mantis "github.com/sphireinc/mantis/http"
	"github.com/valyala/fasthttp"
)

type Config struct {
	Application   `json:"application"`
	Components    `json:"components"`
	Server        `json:"server"`
	Middleware    `json:"middleware"`
	Router        `json:"router"`
	Log           `json:"log"`
	Environment   `json:"environment"`
	Persistence   `json:"persistence"`
	Communication `json:"communication"`
	S             mantis.Status
	Runtime       Runtime `json:"-"`
}

type Context = routing.Context
type Res = mantis.Response

func (c *Config) Start() error {
	c.Router.app = c

	if c.Components.StatsView && c.Runtime.Stats == nil {
		view := statsview.New()
		c.Runtime.Stats = view

		go func() {
			_ = view.Start()
		}()
	}

	if c.Server.Address == "" {
		c.Server.Address = "127.0.0.1"
		c.logInfo("no server address in config, defaulting to 127.0.0.1")
	}

	if c.Server.Port == "" {
		c.Server.Port = "8080"
		c.logInfo("no server port in config, defaulting to 8080")
	}

	if c.Router.Router == nil {
		c.Router.Router = routing.New()
		wireDefaultHandlers(c.Router.Router)
	}

	if !c.Runtime.routesLoaded {
		c.Router.load()
		c.Runtime.routesLoaded = true
	}

	addr := net.JoinHostPort(c.Server.Address, c.Server.Port)

	server := &fasthttp.Server{
		Handler:      c.Router.Router.HandleRequest,
		ReadTimeout:  c.Server.ReadTimeout.Or(0),
		WriteTimeout: c.Server.WriteTimeout.Or(0),
	}

	c.Runtime.Server = server

	c.logInfo("serving app - BON VOYAGE!")
	return server.ListenAndServe(addr)
}

// Run keeps backward compatibility with existing callers.
func (c *Config) Run() error {
	return c.Start()
}

func (c *Config) Stop(ctx context.Context) error {
	serverErr := shutdownServer(ctx, c.Runtime.Server)

	statsErr := stopComponent(c.Runtime.Stats)
	redisErr := stopComponent(&c.Persistence.Redis)
	mysqlErr := stopComponent(&c.Persistence.MySQL)
	bigCacheErr := stopComponent(&c.Persistence.BigCache)
	memCacheErr := stopComponent(&c.Persistence.MemCache)

	return joinErrors(serverErr, statsErr, redisErr, mysqlErr, bigCacheErr, memCacheErr)
}

// Load takes our Config object and loads our environment defined JSON config
func (c *Config) Load() error {
	file, err := os.ReadFile(c.Environment.Location)
	if err != nil {
		return fmt.Errorf("cannot read config file %s: %w", c.Environment.Location, err)
	}

	normalized := normalizeConfigJSON(file)
	if err := json.Unmarshal(normalized, c); err != nil {
		return fmt.Errorf("error loading config file %s: %w", c.Environment.Location, err)
	}

	c.Router.app = c

	// Fill statuses via interface to reduce direct coupling.
	if filler, ok := interface{}(&c.S).(StatusFiller); ok {
		filler.Fill()
	} else {
		c.S.Fill()
	}

	return nil
}
