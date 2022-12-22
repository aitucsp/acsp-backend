package config

import (
	"github.com/gofiber/fiber/v2"
)

func FiberConfig(appCfg *Config) fiber.Config {
	return fiber.Config{
		ReadTimeout:  appCfg.HTTP.ReadTimeout,
		WriteTimeout: appCfg.HTTP.WriteTimeout,
	}
}
