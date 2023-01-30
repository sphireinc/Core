package core

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

type Handler []routing.Handler

var (
	MethodGet    = []string{"GET"}
	MethodPut    = []string{"PUT"}
	MethodPost   = []string{"POST"}
	MethodDelete = []string{"DELETE"}
)

func (c *Config) Methods() []string {
	return []string{"GET,PUT,POST,DELETE"}
}

func (c *Config) Method(method string) []string {
	return []string{method}
}

func (c *Config) M(method string) []string {
	return c.Method(method)
}
