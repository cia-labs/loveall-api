package handlers

import (
	"log"
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
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RefreshHandler(c *gin.Context) {
	refreshTokenString := c.GetHeader("Refresh-Token")
	if refreshTokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is required"})
		c.Abort()
		return
	}

	// Check if the refresh token is valid and active

	// // Verify the refresh token
	// refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return jwtSecret, nil
	// })
	// Parse refresh token
	secretKey := "your-secret-key"
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// Return the secret key used to sign the token.
		return []byte(secretKey), nil

		// TODO : fix this
		// return token, nil
		// return jwtSecret, nil
	})

	if err != nil || !refreshToken.Valid {
		allowExpire := true
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				allowExpire = false
			}
		}
		if allowExpire {
			log.Println("dun:", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token eerr"})
			c.Abort()
			return
		}
	}

	// claims, ok := refreshToken.Claims.(*Claims)
	// if !ok || !refreshToken.Valid || claims.Refresh == false {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
	// 	return
	// }

	// Generate a new access token
	var user models.User
	// userID, ok := refreshToken.Claims.(jwt.MapClaims)["user_id"]
	email, ok := refreshToken.Claims.(jwt.MapClaims)["user_email"]
	log.Println(refreshToken.Claims)
	log.Println("UID", email)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}
	userErr := database.Db.Debug().Model(&models.User{}).Where("email = ?", email).First(&user).Error
	if userErr != nil {
		// log.Println(userErr.Error.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	newAccessToken, newRefreshToken, err := GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken, "refresh_token": newRefreshToken})
	return

	// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token end"})
}

// Login godoc
// @Summary Authenticate user and create a session
// @Description Authenticate a user and create a session
// @Tags authentication
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 200 {object} handlers.LoginResponse
// @Failure 401
// @Router /login [post]
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
		log.Println("DEBUG: Compare Hash error:", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// If the email and password are valid, generate a JWT token.
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"user_id":    user.ID,
	// 	"user_email": user.Email,
	// 	"user_role":  user.Role,
	// 	// "exp":        time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours.
	// 	"exp": time.Now().Add(time.Second * 10).Unix(), // Expires in 24 hours.
	// })
	// Replace "your-secret-key" with your actual secret key.
	// tokenString, err := token.SignedString([]byte("your-secret-key"))
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Internal server error",
	// 	})
	// 	return
	// }

	// refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"user_id":    user.ID,
	// 	"user_email": user.Email,
	// 	"user_role":  user.Role,
	// 	"exp":        time.Now().Add(time.Hour * 24 * 7).Unix(), // Expires in 7 days.
	// })

	// refreshTokenString, refreshErr := refreshToken.SignedString([]byte("your-secret-key-refresh"))
	// if refreshErr != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Internal server error",
	// 	})
	// 	return
	// }
	tokenString, refreshTokenString, tokenGenErr := GenerateTokenPair(user)
	if tokenGenErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	// Return the token in the response.
	c.JSON(http.StatusOK, LoginResponse{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	})
}

func GenerateTokenPair(user models.User) (string, string, error) {
	jwtSecret := []byte("your-secret-key")
	// Create access token with a short expiration time
	// expirationTime := time.Now().Add(24 * time.Hour)
	expirationTime := time.Now().Add(5 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"user_email": user.Email,
		"user_role":  user.Role,

		"exp": expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Create refresh token with a longer expiration time
	// refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour)
	refreshExpirationTime := time.Now().Add(6 * time.Minute)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"user_email": user.Email,
		"user_role":  user.Role,
		"exp":        refreshExpirationTime.Unix(),
		// "exp":        expirationTime.Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Store the active refresh token
	// activeRefreshTokens = append(activeRefreshTokens, RefreshToken{
	// 	Username:       username,
	// 	RefreshToken:   refreshTokenString,
	// 	ExpirationTime: refreshExpirationTime,
	// })

	return tokenString, refreshTokenString, nil
}
