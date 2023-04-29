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

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"

	"acsp/internal/config"
	"acsp/internal/constants"
	"acsp/internal/handler"
	"acsp/internal/infrastructure/db"
	awsS3 "acsp/internal/infrastructure/s3"
	"acsp/internal/logging"
	"acsp/internal/repository"
	"acsp/internal/service"
)

// @title ACSP Backend
// @version 1.0
// @description Backend for AITU Corporate Self-Study Portal.

// @host squid-app-8kray.ondigitalocean.app
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

	// Initialize S3 Wasabi bucket
	// s3Cfg, err := awsConfig.LoadDefaultConfig(context.Background())
	// if err != nil {
	// 	appLogger.Fatal("Error occurred when initializing S3 config: ", zap.Error(err))
	// }
	appLogger.Info("Initializing S3 configuration")
	s3Session, err := awsS3.NewSessionBuilder().
		WithAWSConfig(appConfig.Bucket).
		NewSession()

	// cfg := aws.Config{
	// 	Region: aws.String("eu-central-2"),
	// 	Credentials: credentials.NewStaticCredentials(
	// 		"3V7EQMDA1DCT1IHTTRZ9",
	// 		"GAQyVDS3hcKA5ueufQEQK20ibDuh2eNHkaMzVGLL",
	// 		"",
	// 	),
	// 	Endpoint: aws.String("https://s3.eu-central-2.wasabisys.com"),
	// }
	appLogger.Info("Creating new session")
	// sess, err := session.NewSession(&cfg)
	// if err != nil {
	// 	appLogger.Fatal("Error occurred when initializing S3 session: ", zap.Error(err))
	// }
	// user, err := user.Current()
	// if err != nil {
	// 	appLogger.Fatal("Error occurred when getting current user: ", zap.Error(err))
	// }
	//
	// homedir := user.HomeDir
	// desk := homedir + "\\Desktop\\me.jpg"
	// f, err := os.Open(desk)
	// if err != nil {
	// 	appLogger.Fatal("Error occurred when opening file: ", zap.Error(err))
	// }
	//
	// svc := s3.New(sess)
	//
	// _, err = svc.PutObject(&s3.PutObjectInput{
	// 	Bucket: aws.String("acsp-bucket"),
	// 	Key:    aws.String("test"),
	// 	Body:   f,
	// 	ACL:    aws.String("public-read"),
	// })
	// if err != nil {
	// 	appLogger.Fatal("Error occurred when uploading file to S3: ", zap.Error(err))
	// }
	// toSave, err := svc.GetObject(&s3.GetObjectInput{
	// 	Bucket: aws.String("acsp-bucket"),
	// 	Key:    aws.String("test"),
	// })
	//
	// dst, err := os.Create(desk)
	// // save toSave file in desktop directory
	// _, err = io.Copy(dst, toSave.Body)
	// if err != nil {
	// 	appLogger.Fatal("Error occurred when creating file: ", zap.Error(err))
	// }
	// defer func(dst *os.File) {
	// 	err := dst.Close()
	// 	if err != nil {
	// 		appLogger.Fatal("Error occurred when closing file: ", zap.Error(err))
	// 	}
	// }(dst)

	// client := s3.NewFromConfig(s3Cfg)
	// uploader := manager.NewUploader(client)

	// Initializing context with timeout and logger for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), constants.ContextTimeoutSeconds*time.Second)
	ctx = logging.ContextWithLogger(ctx, appLogger)

	appLogger.Info("Initializing database client")

	// Initializing database client
	dbClient, err := db.NewDBClient(ctx, cancel, postgresConfig)
	if err != nil {
		appLogger.Fatal("Error when initializing Postgres Client", zap.Error(err))
	}
	defer func(dbClient *sqlx.DB) {
		appLogger.Info("Closing database client")

		err := dbClient.Close()
		if err != nil {
			appLogger.Fatal("Error when closing Postgres Client", zap.Error(err))
		}
	}(dbClient)

	appLogger.Info("Initializing redis client")

	// Initializing redis client
	redisClient, err := db.NewClientRedis(ctx, cancel, redisConfig)
	if err != nil {
		appLogger.Fatal("Error when initializing Redis Client", zap.Error(err))
	}

	defer func(redisClient *redis.Client) {
		appLogger.Info("Closing redis client")

		err := redisClient.Close()
		if err != nil {
			appLogger.Fatal("Error when closing Redis Client", zap.Error(err))
		}
	}(redisClient)

	// Initializing database engine with database client and redis client
	dbEngine := db.NewDBEngine(dbClient, *redisClient)

	appLogger.Info("Initializing router and middlewares")

	// Initializing fiber app with fiber config and logger
	app := fiber.New(fiberConfig)

	// Initializing built-in logger and recover middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
	}))

	// Initializing app repository, service and handler
	appRepository := repository.NewRepository(dbEngine.DB, s3Session)
	appService := service.NewService(appRepository, &dbEngine.Cache, *appConfig.Auth)
	appHandler := handler.NewHandler(appService)

	appLogger.Info("Initializing app routes and handlers")

	// Initializing routes with fiber app
	app = appHandler.InitRoutesFiber(app)

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

	// Shutdown the server with a timeout of 10 seconds and log the error if any
	if err := app.Shutdown(); err != nil {
		panic(err)
	} else {
		appLogger.Info("Application shutdown")
	}
}
