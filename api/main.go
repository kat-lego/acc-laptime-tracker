package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	repo := repos.NewPostgresAccSessionsRepo(os.Getenv("ACCLTRCR_POSTGRES_CONNECTION_STRING"))

	router := gin.Default()
	router.Use(middleware.RateLimiter(rate.Every(time.Second/5), 10))
	router.GET("/api/sessions", handlers.GetSessionsHandler(repo))

	logger.Info("starting server", zap.String("addr", ":8080"))
	router.Run(":8080")
}
