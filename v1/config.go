package core

import (
	mantisCache "github.com/sphireinc/mantis/cache"
	mantisDatabase "github.com/sphireinc/mantis/database"
	mantisLog "github.com/sphireinc/mantis/log"
	cache "github.com/victorspringer/http-cache"

	"github.com/qiangxue/fasthttp-routing"
	"time"
)

var EnvironmentVariableName = "SPK_ENV"

type Application struct {
	UUID    string `json:"uuid,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Runtime string `json:"runtime,omitempty"`
}

type Components struct {
	BigCache  bool `json:"bigCache,omitempty"`
	MemCache  bool `json:"memCache,omitempty"`
	Memory    bool `json:"memory,omitempty"`
	Redis     bool `json:"redis,omitempty"`
	MySQL     bool `json:"mysql,omitempty"`
	HTTPCache bool `json:"HTTPCache,omitempty"`
	Log       bool `json:"log,omitempty"`
	StatsView bool `json:"stats_view"`
}

type Server struct {
	Address        string         `json:"address,omitempty"`
	Port           string         `json:"port,omitempty"`
	WriteTimeout   time.Duration  `json:"writeTimeout,omitempty"`
	ReadTimeout    time.Duration  `json:"readTimeout,omitempty"`
	MemCacheTime   time.Duration  `json:"memCacheTime,omitempty"`
	ResponseConfig ResponseConfig `json:"response_config"`
}

type ResponseConfig struct {
	RequestId    bool `json:"request_id"`
	SessionToken bool `json:"session_token"`
}

type Middleware struct {
	Limiter struct {
		Enabled bool `json:"enabled"`
		Time    int  `json:"time"`
		Max     int  `json:"max"`
	} `json:"limiter"`
	LogRequest struct {
		Enabled bool `json:"enabled"`
	} `json:"log_request"`
	BasicHeaders struct {
		Enabled    bool `json:"enabled"`
		ShowServer bool `json:"show_server"`
	} `json:"basic_headers"`
	AdminOnly struct {
		Enabled bool `json:"enabled"`
	} `json:"admin_only"`
	BasicAuth struct {
		Enabled bool `json:"enabled"`
	} `json:"basic_auth"`
}

type Router struct {
	Routes    []route `json:"routes,omitempty"`
	Router    *routing.Router
	httpCache *cache.Client
}

type route struct {
	Name        string            `json:"name,omitempty"`
	Method      string            `json:"method,omitempty"`
	URI         string            `json:"uri,omitempty"`
	Middlewares []routing.Handler `json:"-"`
	Handler     routing.Handler   `json:"-"`
}

type Log struct {
	Writer   *mantisLog.Log `json:"log,omitempty"`
	Location string         `json:"location,omitempty"`
}

type Environment struct {
	Environment string `json:"environment,omitempty"`
	Location    string `json:"location,omitempty"`
}

type Persistence struct {
	BigCache mantisCache.BigCache `json:"bigCache"`
	MemCache mantisCache.MemCache `json:"memCache"`
	Memory   *mantisCache.Memory  `json:"memory"`
	Redis    mantisDatabase.Redis `json:"redis"`
	MySQL    mantisDatabase.MySQL `json:"mysql"`
	Neo4j    mantisDatabase.Neo4j `json:"neo4J"`
}

type Communication struct {
	Email Emailer `json:"email"`
}

type Emailer struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
}
