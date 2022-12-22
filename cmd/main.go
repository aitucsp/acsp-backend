package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"acsp/internal/config"
	"acsp/internal/handler"
	"acsp/internal/logs"
	"acsp/internal/repository"
	"acsp/internal/service"
)

const (
	configsDirectory      = "configs"
	contextTimeoutSeconds = 10
)

// @title           Go REST API
// @version         1.0
// @description Articles REST API for Go.

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Initializing a sugared logger
	err := logs.InitLogger()
	if err != nil {
		log.Fatalf("Logger error: %s", err.Error())
	}

	// Initializing the application configuration
	appConfig, err := config.Init(configsDirectory)
	if err != nil {
		logs.Log().Error(err.Error())
		os.Exit(1)
	}

	// Initializing fiber configuration
	fiberConfig := config.FiberConfig(appConfig)

	postgresConfig, err := config.LoadConfig(".")
	if err != nil {
		logs.Log().Fatalf("Error occurred when initializing database config: %s", err)
	}

	// Declaring a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeoutSeconds*time.Second)

	// Initializing mongo db client
	postgresClient, err := config.NewClientPostgres(ctx, cancel, postgresConfig)
	if err != nil {
		logs.Log().Error(err.Error())
		os.Exit(1)
	}

	app := fiber.New(fiberConfig)
	app.Use(logger.New())

	appRepository := repository.NewRepository(postgresClient)
	appService := service.NewService(appRepository)
	appHandler := handler.NewHandler(appService)

	appHandler.InitRoutesFiber(app)

	go start(app, appConfig.HTTP.Port)

	stopChannel, closeChannel := createChannel()
	defer closeChannel()

	logs.Log().Info("Notified ", <-stopChannel)
	shutdown(ctx, app)
}

func start(server *fiber.App, port string) {
	logs.Log().Info("Application started")
	if err := server.Listen(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	} else {
		logs.Log().Info("application stopped gracefully")
	}
}

func createChannel() (chan os.Signal, func()) {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopChannel, func() {
		close(stopChannel)
	}
}

func shutdown(ctx context.Context, app *fiber.App) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutSeconds*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		panic(err)
	} else {
		logs.Log().Info("Application shutdown")
	}
}
