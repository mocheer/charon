package agent

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/pluto/fs"
)

type AgentConfig struct {
	Name      string
	URL       string
	QueryArgs string
}

var config []AgentConfig

// Use
func Use(api fiber.Router) {
	router := api.Group("agent")
	//
	if fs.IsExist(agentConfigPath) {
		err := fs.ReadJSON(agentConfigPath, &config)
		if err == nil {
			for _, cf := range config {
				router.All("/"+cf.Name+"/*", func(c *fiber.Ctx) error {
					return proxyURL(c, cf.URL+c.Params("*")+cf.QueryArgs+c.Context().QueryArgs().String())
				})
			}
		}
	}
	//
	router.All("/*", proxyHandle)
}

// proxyHandle
func proxyHandle(c *fiber.Ctx) error {
	url := c.Params("*") + "?" + c.Context().QueryArgs().String()
	return proxyURL(c, url)
}
