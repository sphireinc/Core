package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
)

type Handler []routing.Handler

const configContextKey = "__core_config"

var (
	MethodGet    = []string{"GET"}
	MethodPut    = []string{"PUT"}
	MethodPost   = []string{"POST"}
	MethodDelete = []string{"DELETE"}
)

func (c *Config) Methods() []string {
	return []string{"GET", "PUT", "POST", "DELETE"}
}

func (c *Config) Method(method string) []string {
	return []string{method}
}

func (c *Config) M(method string) []string {
	return c.Method(method)
}

func (c *Config) logInfo(msg string) {
	logger := c.logger()
	if logger != nil {
		logger.Info(msg)
	}
}

func configFromContext(ctx *Context) *Config {
	if ctx == nil {
		return nil
	}

	val := ctx.RequestCtx.UserValue(configContextKey)
	if val == nil {
		return nil
	}

	cfg, _ := val.(*Config)
	return cfg
}

func injectConfig(cfg *Config) routing.Handler {
	return func(ctx *Context) error {
		if cfg != nil {
			ctx.RequestCtx.SetUserValue(configContextKey, cfg)
		}
		return nil
	}
}

func currentEnvironment() string {
	if env := strings.TrimSpace(os.Getenv(EnvironmentVariableName)); env != "" {
		return env
	}
	if env := strings.TrimSpace(os.Getenv(LegacyEnvironmentVariableName)); env != "" {
		return env
	}
	return "dev"
}

func explicitConfigPath() string {
	if p := strings.TrimSpace(os.Getenv(ConfigPathEnvironmentVariableName)); p != "" {
		return p
	}
	if p := strings.TrimSpace(os.Getenv(LegacyConfigPathEnvironmentVariableName)); p != "" {
		return p
	}
	return ""
}

// configLocations returns deterministic config search locations.
// Priority is:
// 1. explicit env override handled in findConfigFile()
// 2. current working directory
// 3. executable directory
func configLocations() []string {
	locations := []string{"."}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		if exeDir != "." {
			locations = append(locations, exeDir)
		}
	}

	return locations
}

// findConfigFile searches for the configuration file in deterministic locations.
// Priority:
// 1. ConfigPathEnvironmentVariableName
// 2. working directory
// 3. executable directory
func findConfigFile(locations []string, configFilename string) (string, error) {
	if explicit := explicitConfigPath(); explicit != "" {
		resolved, err := filepath.Abs(explicit)
		if err != nil {
			return "", fmt.Errorf("cannot resolve config path %q: %w", explicit, err)
		}
		if _, err := os.Stat(resolved); err != nil {
			return "", fmt.Errorf("configuration file %q not found: %w", resolved, err)
		}
		return resolved, nil
	}

	for _, loc := range locations {
		configFile := filepath.Join(loc, configFilename)
		resolved, err := filepath.Abs(configFile)
		if err != nil {
			continue
		}
		if _, err := os.Stat(resolved); err == nil {
			return resolved, nil
		}
	}

	return "", fmt.Errorf(
		"configuration file %q not found; set %s (preferred) or %s to an absolute path, or place %s in the working directory",
		configFilename,
		ConfigPathEnvironmentVariableName,
		LegacyConfigPathEnvironmentVariableName,
		configFilename,
	)
}

func normalizeConfigJSON(raw []byte) []byte {
	replacer := strings.NewReplacer(
		`"bigCache"`, `"big_cache"`,
		`"memCache"`, `"mem_cache"`,
		`"HTTPCache"`, `"http_cache"`,
		`"mySQL"`, `"mysql"`,
		`"neo4J"`, `"neo4j"`,
		`"writeTimeout"`, `"write_timeout"`,
		`"readTimeout"`, `"read_timeout"`,
		`"memCacheTime"`, `"mem_cache_ttl"`,
		`"requestId"`, `"request_id"`,
		`"sessionToken"`, `"session_token"`,
		`"printToTerm"`, `"print_to_term"`,
		`"inMemory"`, `"in_memory"`,
	)
	return []byte(replacer.Replace(string(raw)))
}

func joinRoute(prefix, uri string) string {
	if prefix == "" {
		if uri == "" {
			return "/"
		}
		if strings.HasPrefix(uri, "/") {
			return uri
		}
		return "/" + uri
	}

	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	prefix = strings.TrimRight(prefix, "/")

	if uri == "" || uri == "/" {
		return prefix
	}
	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}
	return prefix + uri
}
