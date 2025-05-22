package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kat-lego/acc-laptime-tracker/pkg/repos"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

func GetSessionsHandler(
	repo *repos.FirebaseSessionRepo,
	ccache *cache.Cache,
	logger *zap.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const cacheKey = "recent_sessions"

		if cached, found := ccache.Get(cacheKey); found {
			logger.Sugar().Info("Serving sessions from cache")
			c.JSON(http.StatusOK, gin.H{
				"sessions": cached,
				"cached":   true,
			})
			return
		}

		sessions, err := repo.GetRecentSessions(context.Background())
		if err != nil {
			logger.Sugar().Errorf("Failed to fetch sessions: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
			return
		}

		ccache.Set(cacheKey, sessions, 5*time.Minute)

		c.JSON(http.StatusOK, gin.H{
			"sessions": sessions,
			"total":    len(sessions),
			"cached":   false,
		})
	}
}
