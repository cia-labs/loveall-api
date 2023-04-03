package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {

	secretKey := "your-secret-key"
	return func(c *gin.Context) {

		if c.Request.URL.Path == "/api/v1/users" && c.Request.Method == http.MethodPost {
			c.Next()
			return
		}
		if c.Request.URL.Path == "/api/v1/login" && c.Request.Method == http.MethodPost {
			log.Println("DEBUG: loginging")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "login detect",
			})
			c.Next()
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing authorization header",
			})
			return
		}

		// Check if the authorization header starts with "Bearer ".
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header",
			})
			return
		}

		// Extract the JWT token from the authorization header.
		tokenString := authHeader[len("Bearer "):]

		// Parse the JWT token and validate the signature.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing algorithm.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			// Return the secret key used to sign the token.
			return []byte(secretKey), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization token",
			})
			return
		}

		// Verify that the token is valid and has not expired.
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization token",
			})
			return
		}

		// Extract the user ID and role from the token's claims.
		userID, ok := token.Claims.(jwt.MapClaims)["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userEmail, ok := token.Claims.(jwt.MapClaims)["user_email"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userRole, ok := token.Claims.(jwt.MapClaims)["user_role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Set("user_id", userID)
		c.Set("user_email", userEmail)
		c.Set("user_role", userRole)

		// TODO: fix casin policy
		// obj := c.Request.URL.Path
		// act := c.Request.Method

		// // Casbin enforces policy
		// enfOk, enfErr := enforce(userRole, obj, act)
		// // ok, err := enforce(val.(string), obj, act)
		// if enfErr != nil {
		// 	log.Println(enfErr)
		// 	c.AbortWithStatusJSON(500, gin.H{"error": "error occurred when authorizing user"})
		// 	return
		// }
		// if !enfOk {
		// 	c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
		// 	return
		// }
		c.Next()

		// Check if the user is an admin, and if so, allow access to all endpoints.
		// if userRole == "admin" {
		// 	c.Next()
		// 	return
		// }

		// Check the endpoint's permission level, and allow access if it matches the user's role.
		// endpointPermission := getEndpointPermission(c.Request.Method, c.Request.URL.Path)
		// if endpointPermission == "" {
		// 	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Endpoint not found", "ep": c.Request.URL.Path})
		// 	return
		// }
		// if endpointPermission != userRole {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied", "ep": c.Request.URL.Path})
		// 	return
		// }
		// c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "test", "ep": c.Request.URL.Path})
		// c.Next()

	}
}

func getEndpointPermission(method, path string) string {
	// Map the endpoint paths to their required permission levels.
	permissions := map[string]string{
		"GET:/api/v1/users":       "admin",
		"POST:/api/v1/users":      "admin",
		"GET:/api/v1/users/:id":   "user",
		"PATCH:/api/v1/users/:id": "admin",
		// Add more endpoint permissions here as needed.
	}

	// Get the permission level for the requested endpoint.
	key := method + ":" + path
	return permissions[key]
}

func enforce(sub string, obj string, act string) (bool, error) {
	log.Println("input", "\n\n", sub, obj, act, "\n\n")
	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer(
		"/Users/surya.m/Documents/CIAECO/love-all-backend/internal/config/model.conf",
		"/Users/surya.m/Documents/CIAECO/love-all-backend/internal/config/policy.csv")

	// logger = zaplog.TempLogger.getLogger()

	enforcer.EnableLog(true)

	log.Println("enforcer err? : ", "\n\n", sub, obj, act, enforcer, err, "\n\n")
	if err != nil {
		return false, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}
	// Load policies from DB dynamically
	err = enforcer.LoadPolicy()
	log.Println("loadpolicy", "\n\n", sub, obj, act, err, "\n\n")
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	// Verify
	ok, err := enforcer.Enforce(sub, obj, act)
	log.Println("retun", "\n\n", sub, obj, act, ok, err, "\n\n")
	return ok, err
}
