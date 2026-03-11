package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/qiangxue/fasthttp-routing"
	mantisCache "github.com/sphireinc/mantis/cache"
	mantisDatabase "github.com/sphireinc/mantis/database"
	mantisLog "github.com/sphireinc/mantis/log"
	cache "github.com/victorspringer/http-cache"
)

const EnvironmentVariableName = "SPC_ENV"
const LegacyEnvironmentVariableName = "SPK_ENV"

const ConfigPathEnvironmentVariableName = "SPC_CONFIG"
const LegacyConfigPathEnvironmentVariableName = "SPK_CONFIG"

type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*d = 0
		return nil
	}

	// Preferred form: "5s", "250ms", etc.
	if data[0] == '"' {
		var raw string
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}
		if raw == "" {
			*d = 0
			return nil
		}

		parsed, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("invalid duration %q: %w", raw, err)
		}

		*d = Duration(parsed)
		return nil
	}

	// Backward-compatible form: numbers are treated as seconds.
	var seconds float64
	if err := json.Unmarshal(data, &seconds); err != nil {
		return err
	}

	*d = Duration(time.Duration(seconds * float64(time.Second)))
	return nil
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	if d == nil || time.Duration(*d) <= 0 {
		return json.Marshal("")
	}
	return json.Marshal(time.Duration(*d).String())
}

func (d *Duration) Std() time.Duration {
	if d == nil {
		return 0
	}
	return time.Duration(*d)
}

func (d *Duration) Or(defaultValue time.Duration) time.Duration {
	if d == nil || time.Duration(*d) <= 0 {
		return defaultValue
	}
	return time.Duration(*d)
}

type Logger interface {
	Info(string, ...string)
}

type MemoryStore interface {
	Get(key uint64) (interface{}, bool)
	Set(key uint64, value interface{}, expiresAt time.Time)
}

type Initializer interface {
	Init() error
}

type Connector interface {
	Connect() error
}

type StatusFiller interface {
	Fill()
}

type Application struct {
	UUID    string `json:"uuid,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Runtime string `json:"runtime,omitempty"`
}

type Components struct {
	BigCache  bool `json:"big_cache,omitempty"`
	MemCache  bool `json:"mem_cache,omitempty"`
	Memory    bool `json:"memory,omitempty"`
	Redis     bool `json:"redis,omitempty"`
	MySQL     bool `json:"mysql,omitempty"`
	HTTPCache bool `json:"http_cache,omitempty"`
	Log       bool `json:"log,omitempty"`
	StatsView bool `json:"stats_view,omitempty"`
}

type Server struct {
	Address        string         `json:"address,omitempty"`
	Port           string         `json:"port,omitempty"`
	WriteTimeout   Duration       `json:"write_timeout,omitempty"`
	ReadTimeout    Duration       `json:"read_timeout,omitempty"`
	MemCacheTTL    Duration       `json:"mem_cache_ttl,omitempty"`
	ResponseConfig ResponseConfig `json:"response_config,omitempty"`
}

type ResponseConfig struct {
	RequestID    bool `json:"request_id,omitempty"`
	SessionToken bool `json:"session_token,omitempty"`
}

type Middleware struct {
	Limiter struct {
		Enabled bool     `json:"enabled,omitempty"`
		Time    Duration `json:"time,omitempty"`
		Max     int      `json:"max,omitempty"`
	} `json:"limiter,omitempty"`
	LogRequest struct {
		Enabled bool `json:"enabled,omitempty"`
	} `json:"log_request,omitempty"`
	BasicHeaders struct {
		Enabled    bool `json:"enabled,omitempty"`
		ShowServer bool `json:"show_server,omitempty"`
	} `json:"basic_headers,omitempty"`
	AdminOnly struct {
		Enabled bool `json:"enabled,omitempty"`
	} `json:"admin_only,omitempty"`
	BasicAuth struct {
		Enabled bool `json:"enabled,omitempty"`
	} `json:"basic_auth,omitempty"`
}

type Router struct {
	Routes    []route         `json:"routes,omitempty"`
	Router    *routing.Router `json:"router,omitempty"`
	httpCache *cache.Client
	app       *Config
}

type route struct {
	Method      string            `json:"method,omitempty"`
	URI         string            `json:"uri,omitempty"`
	Middlewares []routing.Handler `json:"-"`
	Handler     routing.Handler   `json:"-"`
}

type LoggerOptions struct {
	PrintToTerm bool `json:"print_to_term,omitempty"`
	Overwrite   bool `json:"overwrite,omitempty"`
}

type Log struct {
	Writer   *mantisLog.Log `json:"-"`
	Location string         `json:"location,omitempty"`
	Options  LoggerOptions  `json:"options,omitempty"`
}

type Environment struct {
	Environment string `json:"environment,omitempty"`
	Location    string `json:"location,omitempty"`
}

type InMemoryConfig struct {
	Capacity int      `json:"capacity,omitempty"`
	TTL      Duration `json:"ttl,omitempty"`
}

type Persistence struct {
	InMemory InMemoryConfig       `json:"in_memory,omitempty"`
	Memory   *mantisCache.Memory  `json:"-"`
	BigCache mantisCache.BigCache `json:"big_cache,omitempty"`
	MemCache mantisCache.MemCache `json:"mem_cache,omitempty"`
	Redis    mantisDatabase.Redis `json:"redis,omitempty"`
	MySQL    mantisDatabase.MySQL `json:"mysql,omitempty"`
	Neo4j    mantisDatabase.Neo4j `json:"neo4j,omitempty"`
}

type Communication struct {
	Email Emailer `json:"email,omitempty"`
}

type Emailer struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
}
