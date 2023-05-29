package config

import (
	"github.com/gofiber/fiber/v2"
)

// FiberConfig returns a fiber.Config.
func FiberConfig(appCfg *Config) fiber.Config {
	return fiber.Config{
		ReadTimeout:  appCfg.HTTP.ReadTimeout,
		WriteTimeout: appCfg.HTTP.WriteTimeout,
	}
}
