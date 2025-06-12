package api

import (
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterReportRoutes(router *gin.Engine) {
	reports := router.Group("/api/reports")
	reports.Use(AuthMiddleware())
	{
		reports.GET("", ListReportsHandler)
		reports.GET(":id", GetReportHandler)
		reports.DELETE(":id", DeleteReportHandler)
		reports.GET(":id/suggestions", GetSuggestionsHandler)
	}
}

func getUserIDFromClaims(c *gin.Context) (primitive.ObjectID, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return primitive.NilObjectID, false
	}
	userClaims := claims.(jwt.MapClaims)
	userID, err := primitive.ObjectIDFromHex(userClaims["user_id"].(string))
	if err != nil {
		return primitive.NilObjectID, false
	}
	return userID, true
}

func ListReportsHandler(c *gin.Context) {
	userID, ok := getUserIDFromClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	reports, err := services.ListReportsByUser(c.Request.Context(), userID)
	if err != nil {
		utils.LogAction(userID.Hex(), "list_reports", "failure", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch reports"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": reports})
}

func GetReportHandler(c *gin.Context) {
	userID, ok := getUserIDFromClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	reportID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid report id"})
		return
	}
	report, err := services.GetReportByID(c.Request.Context(), reportID)
	if err != nil || report.UserID != userID {
		utils.LogAction(userID.Hex(), "get_report", "failure", "not found or forbidden")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Report not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": report})
}

func DeleteReportHandler(c *gin.Context) {
	userID, ok := getUserIDFromClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	reportID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid report id"})
		return
	}
	report, err := services.GetReportByID(c.Request.Context(), reportID)
	if err != nil || report.UserID != userID {
		utils.LogAction(userID.Hex(), "delete_report", "failure", "not found or forbidden")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Report not found"})
		return
	}
	err = services.DeleteReportByID(c.Request.Context(), reportID)
	if err != nil {
		utils.LogAction(userID.Hex(), "delete_report", "failure", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete report"})
		return
	}
	utils.LogAction(userID.Hex(), "delete_report", "success", "deleted report "+reportID.Hex())
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Report deleted."})
}

func GetSuggestionsHandler(c *gin.Context) {
	userID, ok := getUserIDFromClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	reportID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid report id"})
		return
	}
	report, err := services.GetReportByID(c.Request.Context(), reportID)
	if err != nil || report.UserID != userID {
		utils.LogAction(userID.Hex(), "get_suggestions", "failure", "not found or forbidden")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Report not found"})
		return
	}
	suggestions, err := services.GetSuggestionsByReportID(c.Request.Context(), reportID)
	if err != nil {
		utils.LogAction(userID.Hex(), "get_suggestions", "failure", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch suggestions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": suggestions})
}
