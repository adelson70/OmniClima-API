package middleware

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if verifyRoutePass(c) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Header sem Auth"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato do header inválido"})
			c.Abort()
			return
		}

		tokenJWT := parts[1]

		token, err := jwt.Parse(tokenJWT, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inesperado")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["userID"])
			c.Set("admin", claims["admin"])
		}

		c.Next()
	}
}

func verifyRoutePass(c *gin.Context) bool {
	routePass := []string{"user/create", "auth/login"}
	fullPath := c.FullPath()
	path := strings.Split(fullPath, "/api/")[1]

	if slices.Contains(routePass, path) {
		fmt.Printf("Rota %s permitida passar sem auth\n", path)
		return true
	}

	fmt.Printf("Rota %s não permitida passar sem auth\n", path)
	return false
}
