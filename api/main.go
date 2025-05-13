package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kat-lego/acc-laptime-tracker/api/handlers"
	"github.com/kat-lego/acc-laptime-tracker/api/middleware"
	"github.com/kat-lego/acc-laptime-tracker/pkg/repos"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var logger *zap.Logger

func setupLogger() *zap.Logger {
	log, err := zap.NewProduction()
	if err != nil {
		panic("cannot initialize zap logger: " + err.Error())
	}
	return log
}

func main() {
	logger = setupLogger()
	defer logger.Sync()

	cosmosConn := os.Getenv("ACC_COSMOS_CONNECTION_STRING")
	cosmosDatabase := os.Getenv("ACC_COSMOS_DATABASE")
	cosmosContainer := os.Getenv("ACC_COSMOS_CONTAINER")
	repo, err := repos.NewCosmosSessionRepo(cosmosConn, cosmosDatabase, cosmosContainer)
	if err != nil {
		logger.Error("failed to connect to cosmos")
		return
	}

	origins := strings.Split(os.Getenv("ACC_CORS_ORIGINS"), ",")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.RateLimiter(rate.Every(time.Second/5), 10))

	router.GET("/api/sessions", handlers.GetSessionsHandler(repo, logger))

	port := os.Getenv("PORT")
	if port == "" {
		logger.Error("port not specified in the environment")
		return
	}

	port = fmt.Sprintf(":%s", port)

	logger.Info("starting server", zap.String("addr", port))
	router.Run(port)
}
