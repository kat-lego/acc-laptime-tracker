package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kat-lego/acc-laptime-tracker/pkg/repos"
)

func GetSessionsHandler(repo *repos.PostgresAccSessionRepo) gin.HandlerFunc {
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

		sessions, count, err := repo.GetSessionsWithCount(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"sessions": sessions,
			"total":    count,
		})
	}
}
