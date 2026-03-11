package core

import (
	"fmt"
	"strconv"
)

func (c *Config) validate() error {
	if c.Server.Port != "" {
		port, err := strconv.Atoi(c.Server.Port)
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("server.port must be a valid TCP port, got %q", c.Server.Port)
		}
	}

	if c.Server.ReadTimeout.Std() < 0 {
		return fmt.Errorf("server.read_timeout cannot be negative")
	}
	if c.Server.WriteTimeout.Std() < 0 {
		return fmt.Errorf("server.write_timeout cannot be negative")
	}
	if c.Server.MemCacheTTL.Std() < 0 {
		return fmt.Errorf("server.mem_cache_ttl cannot be negative")
	}
	if c.Middleware.Limiter.Time.Std() < 0 {
		return fmt.Errorf("middleware.limiter.time cannot be negative")
	}
	if c.Persistence.InMemory.TTL.Std() < 0 {
		return fmt.Errorf("persistence.in_memory.ttl cannot be negative")
	}
	if c.Persistence.InMemory.Capacity < 0 {
		return fmt.Errorf("persistence.in_memory.capacity cannot be negative")
	}
	if c.Middleware.Limiter.Enabled && c.Middleware.Limiter.Max < 1 {
		return fmt.Errorf("middleware.limiter.max must be greater than 0 when limiter is enabled")
	}

	return nil
}
