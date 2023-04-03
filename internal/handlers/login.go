package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/madeinatria/love-all-backend/internal/database"
	"github.com/madeinatria/love-all-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(c *gin.Context) {

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// TODO: use controller
	var user models.User
	userErr := database.Db.Model(&models.User{}).Where("email = ?", req.Email).First(&user)
	if userErr.Error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare the password hash with the given password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// If the email and password are valid, generate a JWT token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"user_email": user.Email,
		"user_role":  user.Role,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours.
	})
	// Replace "your-secret-key" with your actual secret key.
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Return the token in the response.
	c.JSON(http.StatusOK, LoginResponse{
		Token: tokenString,
	})
}
