package core

import (
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"os"
)

type Config struct {
	Application   Application   `json:"application"`
	Components    Components    `json:"components"`
	Server        Server        `json:"server"`
	Router        Router        `json:"router"`
	Log           Log           `json:"log"`
	Environment   Environment   `json:"environment"`
	Persistence   Persistence   `json:"persistence"`
	Communication Communication `json:"communication"`
	Context       *routing.Context
}

// Load takes our Config object and loads our environment defined JSON config
func (c *Config) Load() {
	file, err := os.ReadFile(c.Environment.Location)
	if err != nil {
		panic("Cannot load config file " + c.Environment.Location)
	}
	err = json.Unmarshal([]byte(file), &c)
	if err != nil {
		panic("Error unmarshalling config file " + c.Environment.Location)
	}
}

// Run loads routes and starts the Server.
func (c *Config) Run() {
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
