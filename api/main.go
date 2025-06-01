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
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
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

	projectId := os.Getenv("ACC_FIREBASE_PROJECT_ID")
	database := os.Getenv("ACC_FIREBASE_DATABASE")
	collectionName := os.Getenv("ACC_FIREBASE_COLLECTION")
	repo, err := repos.NewFirebaseSessionRepo(projectId, database, collectionName)
	if err != nil {
		logger.Error("failed to connect to firebase")
		os.Exit(1)
	}

	origins := strings.Split(os.Getenv("ACC_CORS_ORIGINS"), ",")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(middleware.RateLimiter())

	ccache := cache.New(1*time.Minute, 10*time.Minute)
	router.GET("/api/sessions", handlers.GetSessionsHandler(repo, ccache, logger))

	port := os.Getenv("PORT")
	if port == "" {
		logger.Error("port not specified in the environment")
		return
	}

	port = fmt.Sprintf(":%s", port)

	logger.Info("starting server", zap.String("addr", port))
	router.Run(port)
}
