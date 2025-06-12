package api

import (
	"backend/jobs"
	"backend/services"
	"backend/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterAnalyzeRoutes(router *gin.Engine) {
	analyze := router.Group("/api")
	{
		analyze.POST("/analyze", AuthMiddleware(), AnalyzeHandler)
	}
}

func AnalyzeHandler(c *gin.Context) {
	var req struct {
		URL  string `json:"url"`
		HTML string `json:"html"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || (req.URL == "" && req.HTML == "") {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Must provide url or html", "error": err})
		return
	}
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userClaims := claims.(jwt.MapClaims)
	userID, _ := primitive.ObjectIDFromHex(userClaims["user_id"].(string))

	report, err := services.CreateReport(context.Background(), userID, req.URL, req.HTML)
	if err != nil {
		utils.LogAction(userID.Hex(), "analyze", "failure", "failed to create report")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create report"})
		return
	}
	jobs.EnqueueAnalyzeJob(jobs.AnalyzeJob{
		ReportID: report.ID,
		URL:      req.URL,
		HTML:     req.HTML,
	})
	utils.LogAction(userID.Hex(), "analyze", "success", "enqueued analysis job")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analysis started",
		"data":    gin.H{"reportId": report.ID, "status": report.Status, "createdAt": report.CreatedAt},
		"error":   nil,
	})
}
