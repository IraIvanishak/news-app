package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IraIvanishak/news-app/internal/adapter/http/handlers"
	"github.com/IraIvanishak/news-app/internal/adapter/http/middlewares"
	"github.com/IraIvanishak/news-app/internal/adapter/http/routes"
	"github.com/IraIvanishak/news-app/internal/adapter/repositories"
	"github.com/IraIvanishak/news-app/internal/core/ports"
	"github.com/IraIvanishak/news-app/internal/core/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	logLevelStr := viper.GetString("LOG_LEVEL")

	var logLevel zapcore.Level
	switch logLevelStr {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	default:
		logLevel = zap.DebugLevel
	}

	// Create logger with development config for better error handling
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(logLevel)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	return logger
}

func NewMongoClient(logger *zap.Logger) (*mongo.Client, error) {
	mongoURI := viper.GetString("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure MongoDB client options
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", zap.Error(err), zap.String("uri", mongoURI))
		return nil, err
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Error("Failed to ping MongoDB", zap.Error(err))
		return nil, err
	}

	logger.Info("Connected to MongoDB successfully", zap.String("uri", mongoURI))
	return client, nil
}

func NewFiberApp() *fiber.App {
	// Create template engine
	engine := html.New("./web/templates", ".html")

	// Create Fiber app with template engine
	app := fiber.New(fiber.Config{
		Views:        engine,
		ServerHeader: "NewsApp",
		AppName:      "NewsApp",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())

	// Static files
	app.Static("/static", "./web/static")

	return app
}

// StartServer starts the HTTP server
func StartServer(
	lifecycle fx.Lifecycle,
	app *fiber.App,
	logger *zap.Logger,
) {
	serverPort := viper.GetString("LISTEN_PORT")

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					logger.Info("Starting server", zap.String("port", serverPort))
					if err := app.Listen(":" + serverPort); err != nil && err != http.ErrServerClosed {
						logger.Fatal("Server startup failed", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("Stopping server")
				return app.Shutdown()
			},
		},
	)
}

func main() {
	app := fx.New(
		// env
		fx.Invoke(viper.AutomaticEnv),
		// Logging
		fx.Provide(NewLogger),

		// Database
		fx.Provide(NewMongoClient),
		fx.Provide(func(client *mongo.Client) *mongo.Database {
			return client.Database(viper.GetString("MONGO_DATABASE"))
		}),

		fx.Provide(
			fx.Annotate(repositories.NewPostRepository, fx.As(new(ports.PostRepository))),
			fx.Annotate(services.NewPostService, fx.As(new(ports.PostService))),
			fx.Annotate(handlers.NewPostHandler, fx.As(new(ports.PostHandlers))),
		),
		fx.Invoke(routes.PostRoutes),

		fx.Invoke(func(app *fiber.App) {
			middlewares.RoutLoggerMiddlewareInitializer(app)
			middlewares.RecoverMiddlewareInitializer(app)
		}),

		// Web Server
		fx.Provide(NewFiberApp),
		fx.Invoke(StartServer),
	)

	// Run the application
	if err := app.Err(); err != nil {
		log.Fatal(err)
	}
	app.Run()
}
