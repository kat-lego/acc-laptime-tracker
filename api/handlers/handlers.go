package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kat-lego/acc-laptime-tracker/pkg/repos"
	"go.uber.org/zap"
)

func GetSessionsHandler(repo *repos.CosmosSessionRepo, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters: limit and offset
		limitParam := c.DefaultQuery("limit", "10")
		offsetParam := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}

		offset, err := strconv.Atoi(offsetParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}

		sessions, err := repo.GetSessions(limit, offset)
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
