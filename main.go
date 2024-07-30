package main

import (
	"talkspace-api/app/configs"
	"talkspace-api/app/databases"
	"talkspace-api/app/routes"

	// "talkspace-api/docs"
	"talkspace-api/middlewares"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	// swaggerFiles "github.com/swaggo/files"
)

// @title TalkSpace API
// @version 1.0
// @description TalkSpace : Mental Health Care System
// @termsOfService http://swagger.io/terms/
// @host localhost:8000
// @BasePath /

func main() {
	godotenv.Load()
	config, err := configs.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load configuration: %v", err)
	}

	// docs.SwaggerInfo.Title = "TalkSpace API"
	// docs.SwaggerInfo.Description = "TalkSpace : Mental Health Care System"
	// docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:8000"
	// docs.SwaggerInfo.BasePath = "/"
	// docs.SwaggerInfo.Schemes = []string{"http", "https"}

	pdb := databases.ConnectPostgreSQL()
	es := databases.ConnectElasticsearch()
	rdb := databases.ConnectRedis()

	defer rdb.Close()

	e := echo.New()

	middlewares.RemoveTrailingSlash(e)
	middlewares.Logger(e)
	middlewares.RateLimiter(e)
	middlewares.Recover(e)
	middlewares.CORS(e)

	routes.SetupRoutes(e, pdb, es)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	host := config.SERVER.SERVER_HOST
	port := config.SERVER.SERVER_PORT
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "8000"
	}
	address := host + ":" + port

	logrus.Info("server is running on address %s...", address)
	if err := e.Start(address); err != nil {
		logrus.Fatalf("error starting server: %v", err)
	}
}
