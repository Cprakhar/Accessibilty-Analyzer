package api

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// RegisterRoutes sets up authentication routes
func RegisterAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", LoginHandler)
		auth.POST("/register", RegisterHandler)
		auth.POST("/logout", LogoutHandler)
		auth.GET("/me", AuthMiddleware(), MeHandler)
	}
}

// LoginHandler handles user login
func LoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input", "error": err.Error()})
		return
	}
	user, err := services.FindUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		utils.LogAction("", "login", "failure", "email not found: "+req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid email or password"})
		return
	}
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		utils.LogAction(user.ID, "login", "failure", "invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid email or password"})
		return
	}
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.LogAction(user.ID, "login", "failure", "token generation error")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}
	utils.LogAction(user.ID, "login", "success", "user logged in")
	c.JSON(http.StatusOK, gin.H{"success": true, "token": token, "user": gin.H{"_id": user.ID, "email": user.Email, "name": user.Name}})
}

// RegisterHandler handles user registration
func RegisterHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input", "error": err.Error()})
		return
	}
	_, err := services.FindUserByEmail(c.Request.Context(), req.Email)
	if err == nil {
		utils.LogAction("", "register", "failure", "email already registered: "+req.Email)
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Email already registered"})
		return
	}
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.LogAction("", "register", "failure", "hash error")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password"})
		return
	}
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
	}
	err = services.CreateUser(c.Request.Context(), user)
	if err != nil {
		utils.LogAction("", "register", "failure", "db error")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create user"})
		return
	}
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.LogAction(user.ID, "register", "failure", "token generation error")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}
	utils.LogAction(user.ID, "register", "success", "user registered")
	c.JSON(http.StatusOK, gin.H{"success": true, "token": token, "user": gin.H{"_id": user.ID, "email": user.Email, "name": user.Name}})
}

// LogoutHandler handles user logout
func LogoutHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	userID := ""
	if claims != nil {
		userClaims := claims.(jwt.MapClaims)
		userID = userClaims["user_id"].(string)
	}
	utils.LogAction(userID, "logout", "success", "user logged out")
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logged out."})
}

// MeHandler returns current user info
func MeHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(jwt.MapClaims)
	email := userClaims["email"].(string)
	user, err := services.FindUserByEmail(c.Request.Context(), email)
	if err != nil {
		utils.LogAction("", "me", "failure", "user not found: "+email)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "_id": user.ID, "email": user.Email, "name": user.Name, "createdAt": user.CreatedAt})
}

// AuthMiddleware checks JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			utils.LogAction("", "auth", "failure", "missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing Authorization header"})
			return
		}
		tokenStr := header
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			utils.LogAction("", "auth", "failure", "invalid or expired token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired token"})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
