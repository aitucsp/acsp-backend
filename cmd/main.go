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

// @host https://monkfish-app-pxfhy.ondigitalocean.app
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

	appLogger, err := logging.NewBuilder().
		WithFallbackLogger(fallbackLogger).
		WithLoggerConfig(appConfig.Logger).
		WithHostConfig(appConfig.Host).
		NewLogger()

	pid := os.Getpid()
	appLogger = appLogger.With(zap.Int("pid", pid))

	hostname, err := os.Hostname()
	if err != nil {
		appLogger.Error("couldn't get hostname", zap.Error(err))
	} else {
		appLogger = appLogger.With(zap.String("host", hostname))
	}

	appLogger.Info("starting")

	// Initializing fiber configuration
	fiberConfig := config.FiberConfig(appConfig)

	postgresConfig, err := db.LoadPostgresConfig(configProvider)
	if err != nil {
		appLogger.Fatal("Error occurred when initializing database config: ", zap.Error(err))
	}

	redisConfig, err := db.LoadRedisConfig(configProvider)
	if err != nil {
		appLogger.Fatal("Error occurred when initializing database config: ", zap.Error(err))
	}

	// Declaring a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), constants.ContextTimeoutSeconds*time.Second)
	ctx = logging.ContextWithLogger(ctx, appLogger)

	// Initializing postgresql database client
	dbClient, err := db.NewDBClient(ctx, cancel, postgresConfig)
	if err != nil {
		appLogger.Error("Error when initializing Postgres Client", zap.Error(err))
	}

	redisClient, err := db.NewClientRedis(ctx, cancel, redisConfig)
	if err != nil {
		appLogger.Error("Error when initializing Redis Client", zap.Error(err))
		os.Exit(1)
	}

	dbEngine := db.NewDBEngine(dbClient, *redisClient)

	app := fiber.New(fiberConfig)
	app.Use(logger.New())

	appRepository := repository.NewRepository(dbEngine.DB)
	appService := service.NewService(appRepository, &dbEngine.Cache, *appConfig.Auth)
	appHandler := handler.NewHandler(appService)

	appHandler.InitRoutesFiber(app)

	go start(app, appConfig.HTTP.Port, appLogger)

	stopChannel, closeChannel := createChannel()
	defer closeChannel()

	appLogger.Info("Notified ", zap.Any("Channel ", <-stopChannel))
	shutdown(ctx, app, appLogger)
}

func start(server *fiber.App, port string, appLogger *zap.Logger) {
	appLogger.Info("Application started")
	if err := server.Listen(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	} else {
		appLogger.Info("Application stopped gracefully")
	}
}

func createChannel() (chan os.Signal, func()) {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopChannel, func() {
		close(stopChannel)
	}
}

func shutdown(ctx context.Context, app *fiber.App, appLogger *zap.Logger) {
	ctx, cancel := context.WithTimeout(ctx, constants.ContextTimeoutSeconds*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		panic(err)
	} else {
		appLogger.Info("Application shutdown")
	}
}
