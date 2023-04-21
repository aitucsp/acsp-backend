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

	"go.uber.org/zap"

	"acsp/internal/config"
	"acsp/internal/constants"
	"acsp/internal/handler"
	"acsp/internal/infrastructure/db"
	"acsp/internal/logging"
	"acsp/internal/repository"
	"acsp/internal/service"
)

// @title ACSP Backend
// @version 1.0
// @description Backend for AITU Corporate Self-Study Portal.

// @host https://squid-app-8kray.ondigitalocean.app
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	fallbackLogger := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC|log.Lmsgprefix)

	baseProvider, err := config.NewDotenvProvider("base.env")
	if err != nil {
		fallbackLogger.Println("couldn't create config provider:", err)

		return
	}

	configProvider := config.NewProviderChain(baseProvider)
	appConfig, err := config.NewBaseConfig(configProvider)
	if err != nil {
		fallbackLogger.Println("couldn't create config:", err)

		return
	}

	//	Initializing logger
	appLogger, err := logging.NewBuilder().
		WithFallbackLogger(fallbackLogger).
		WithLoggerConfig(appConfig.Logger).
		WithHostConfig(appConfig.Host).
		NewLogger()

	if err != nil {
		fallbackLogger.Println("couldn't create logger:", err)

		return
	}

	// Initializing logger with pid and hostname
	pid := os.Getpid()
	appLogger = appLogger.With(zap.Int("pid", pid))

	// Initializing logger with hostname
	hostname, err := os.Hostname()
	if err != nil {
		appLogger.Error("couldn't get hostname", zap.Error(err))
	} else {
		appLogger = appLogger.With(zap.String("host", hostname))
	}

	appLogger.Info("starting")

	// Initializing fiber config
	fiberConfig := config.FiberConfig(appConfig)

	// Initializing database config
	postgresConfig, err := db.LoadPostgresConfig(configProvider)
	if err != nil {
		appLogger.Fatal("Error occurred when initializing database config: ", zap.Error(err))
	}

	// Initializing redis config
	redisConfig, err := db.LoadRedisConfig(configProvider)
	if err != nil {
		appLogger.Fatal("Error occurred when initializing database config: ", zap.Error(err))
	}

	// Initializing context with timeout and logger for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), constants.ContextTimeoutSeconds*time.Second)
	ctx = logging.ContextWithLogger(ctx, appLogger)

	// Initializing database client
	dbClient, err := db.NewDBClient(ctx, cancel, postgresConfig)
	if err != nil {
		appLogger.Error("Error when initializing Postgres Client", zap.Error(err))
	}

	// Initializing redis client
	redisClient, err := db.NewClientRedis(ctx, cancel, redisConfig)
	if err != nil {
		appLogger.Error("Error when initializing Redis Client", zap.Error(err))
		os.Exit(1)
	}

	// Initializing database engine with database client and redis client
	dbEngine := db.NewDBEngine(dbClient, *redisClient)

	// Initializing fiber app with fiber config and logger
	app := fiber.New(fiberConfig)
	app.Use(logger.New())

	// Initializing app repository, service and handler
	appRepository := repository.NewRepository(dbEngine.DB)
	appService := service.NewService(appRepository, &dbEngine.Cache, *appConfig.Auth)
	appHandler := handler.NewHandler(appService)

	// Initializing routes
	appHandler.InitRoutesFiber(app)

	// Initializing graceful shutdown with fiber app, port and logger
	go start(app, appConfig.HTTP.Port, appLogger)

	// Initializing graceful shutdown with context, fiber app and logger
	stopChannel, closeChannel := createChannel()
	defer closeChannel()

	// Waiting for stop signal
	appLogger.Info("Notified ", zap.Any("Channel ", <-stopChannel))

	// Graceful shutdown with context, fiber app and logger
	shutdown(ctx, app, appLogger)
}

// start starts the server and listens for shutdown signals to gracefully stop the server
func start(server *fiber.App, port string, appLogger *zap.Logger) {
	appLogger.Info("Application started")
	if err := server.Listen(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	} else {
		appLogger.Info("Application stopped gracefully")
	}
}

// createChannel creates a channel to listen for shutdown signals
func createChannel() (chan os.Signal, func()) {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopChannel, func() {
		close(stopChannel)
	}
}

// shutdown gracefully shuts down the server
func shutdown(ctx context.Context, app *fiber.App, appLogger *zap.Logger) {
	ctx, cancel := context.WithTimeout(ctx, constants.ContextTimeoutSeconds*time.Second)
	defer cancel()

	// Shutdown the server with a timeout of 5 seconds and log the error if any
	if err := app.Shutdown(); err != nil {
		panic(err)
	} else {
		appLogger.Info("Application shutdown")
	}
}
