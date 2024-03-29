package core

import (
	"encoding/json"
	"github.com/go-echarts/statsview"
	routing "github.com/qiangxue/fasthttp-routing"
	mantis "github.com/sphireinc/mantis/http"
	"github.com/valyala/fasthttp"
	"os"
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
}

type Context = routing.Context
type Res = mantis.Response

// Run loads routes and starts the Server.
func (c *Config) Run() {
	if c.Components.StatsView {
		go func() {
			_ = statsview.New().Start()
		}()
	}

	// if c has no Address, default to localhost
	if c.Server.Address == "" {
		c.Server.Address = "127.0.0.1"
		c.Server.Port = "8080"
		c.Log.Writer.Info("no server address in config, defaulting to 127.0.0.1:8080")
	}

	// Load our routes
	c.Router.load()

	c.Log.Writer.Info("serving app - BON VOYAGE!")
	if err := fasthttp.ListenAndServe(":"+c.Server.Port, c.Router.Router.HandleRequest); err != nil {
		c.Log.Writer.HandleFatalError(err)
	}
}

// Load takes our Config object and loads our environment defined JSON config
func (c *Config) Load() {
	file, err := os.ReadFile(c.Environment.Location)
	if err != nil {
		panic("Cannot find config file " + c.Environment.Location)
	}
	err = json.Unmarshal([]byte(file), &c)
	if err != nil {
		panic("Error loading config file " + c.Environment.Location)
	}

	// Populate our statuses
	c.S.Fill()
}
