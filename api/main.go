package main

import (
	"fmt"
	"os"
	"time"

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

	router := gin.Default()
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
