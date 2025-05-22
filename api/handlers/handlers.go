package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kat-lego/acc-laptime-tracker/pkg/repos"
	"go.uber.org/zap"
)

func GetSessionsHandler(repo *repos.FirebaseSessionRepo, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions, err := repo.GetRecentSessions(context.Background())
		if err != nil {
			logger.Sugar().Errorf("Failed to fetch sessions: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"sessions": sessions,
			"total":    len(sessions),
		})
	}
}
