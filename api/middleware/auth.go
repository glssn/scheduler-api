package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api/models"
	"github.com/glssn/scheduler-api/initializers"
	"github.com/golang-jwt/jwt"
)

// contains checks if a string slice contains a given string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func RequireAuth(c *gin.Context) {
	// Get the authorization header from the request
	authHeader := c.Request.Header.Get("Authorization")

	// Check if the authorization header is a Bearer token
	if strings.HasPrefix(authHeader, "Bearer ") {
		// Extract the token from the authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check if the token is in the list of allowed tokens
		allowedTokens := strings.Split(os.Getenv("ALLOWED_TOKENS"), ",")
		if contains(allowedTokens, tokenString) {
			// If the token is in the list of allowed tokens, continue with the request
			log.Println("Request authenticated using basic auth token")
			c.Next()
			return
		}
	}

	// If the request is not authenticated using a Bearer token,
	// fall back to using the JWT cookie

	// Get the JWT from cookie
	tokenStringSigned, err := c.Cookie("Authorization")
	if err != nil || tokenStringSigned == "" {
		log.Println("Request not authenticated")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Parse and validate the JWT
	token, err := jwt.Parse(tokenStringSigned, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_AUTH_SECRET_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check JWT expiry
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			log.Println("JWT token expired")
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find JWT subject in the database
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			log.Println("JWT token invalid")
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Attach this user's info to request context
		c.Set("user", user)

		// Continue request
		c.Next()
		return
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
