package core

import (
	"context"
	"errors"
	"strings"

	"github.com/valyala/fasthttp"
)

type Runtime struct {
	Logger       Logger           `json:"-"`
	Memory       MemoryStore      `json:"-"`
	Server       *fasthttp.Server `json:"-"`
	Stats        interface{}      `json:"-"`
	routesLoaded bool
}

type noopLogger struct{}

func (noopLogger) Info(a string, b ...string) {
	_ = a
	_ = b
}

func (c *Config) logger() Logger {
	if c == nil {
		return nil
	}
	if c.Runtime.Logger != nil {
		return c.Runtime.Logger
	}
	if c.Log.Writer != nil {
		return c.Log.Writer
	}
	return nil
}

func (c *Config) memoryStore() MemoryStore {
	if c == nil {
		return nil
	}
	if c.Runtime.Memory != nil {
		return c.Runtime.Memory
	}
	return c.Persistence.Memory
}

func shutdownServer(ctx context.Context, server *fasthttp.Server) error {
	if server == nil {
		return nil
	}

	done := make(chan error, 1)
	go func() {
		done <- server.Shutdown()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func stopComponent(component interface{}) error {
	if component == nil {
		return nil
	}

	switch c := component.(type) {
	case interface{ Shutdown() error }:
		return c.Shutdown()
	case interface{ Close() error }:
		return c.Close()
	case interface{ Disconnect() error }:
		return c.Disconnect()
	case interface{ Stop() error }:
		return c.Stop()
	case interface{ Shutdown() }:
		c.Shutdown()
		return nil
	case interface{ Close() }:
		c.Close()
		return nil
	case interface{ Disconnect() }:
		c.Disconnect()
		return nil
	case interface{ Stop() }:
		c.Stop()
		return nil
	default:
		return nil
	}
}

func joinErrors(errs ...error) error {
	parts := make([]string, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			parts = append(parts, err.Error())
		}
	}
	if len(parts) == 0 {
		return nil
	}
	return errors.New(strings.Join(parts, "; "))
}
