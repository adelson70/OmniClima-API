package middleware

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"OmniClima/internal/devices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func isSuggestionsRoute(c *gin.Context) bool {
	return strings.HasPrefix(c.Request.URL.Path, "/api/suggestions/")
}

func AuthMiddleware(deviceSvc *devices.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		if isSuggestionsRoute(c) {
			ok, reason := verifyDevicePass(c, deviceSvc)
			if ok {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": reason})
			c.Abort()
			return
		}

		if verifyRoutePass(c) {
			c.Next()
			return
		}

		// authHeader := c.GetHeader("Authorization")

		// if authHeader == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Header sem Auth"})
		// 	c.Abort()
		// 	return
		// }

		// parts := strings.Split(authHeader, " ")
		// if len(parts) != 2 || parts[0] != "Bearer" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato do header inválido"})
		// 	c.Abort()
		// 	return
		// }

		// tokenJWT := parts[1]

		// token, err := jwt.Parse(tokenJWT, func(token *jwt.Token) (interface{}, error) {
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("método de assinatura inesperado")
		// 	}
		// 	return []byte(os.Getenv("JWT_SECRET")), nil
		// })

		// if err != nil || !token.Valid {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		// 	c.Abort()
		// 	return
		// }

		// if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 	c.Set("userID", claims["userID"])
		// 	c.Set("admin", claims["admin"])

		verifyAuth(c)
	}
}

func verifyAuth(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Header sem Auth"})
		c.Abort()
		return false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato do header inválido"})
		c.Abort()
		return false
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
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Set("userID", claims["userID"])
		c.Set("admin", claims["admin"])
	}

	c.Next()
	return true
}

func verifyRoutePass(c *gin.Context) bool {
	routePass := []string{"user/create", "user/login", "devices/"}
	fullPath := c.FullPath()
	path := strings.Split(fullPath, "/api/")[1]

	if strings.HasPrefix(path, "suggestions/") {
		return false
	}

	if slices.Contains(routePass, path) {
		fmt.Printf("Rota %s permitida passar sem auth\n", path)
		return true
	}

	fmt.Printf("Rota %s não permitida passar sem auth\n", path)
	return false
}

func verifyDevicePass(
	c *gin.Context,
	deviceSvc *devices.Service,
) (bool, string) {
	authHeader := c.GetHeader("Authorization")

	latStr := c.Param("lat")
	lonStr := c.Param("lon")

	if latStr == "" || lonStr == "" {
		return false, "coordenadas inválidas"
	}

	deviceIDRaw := c.GetHeader("x-device-id")
	signatureHex := c.GetHeader("x-signature")

	if deviceIDRaw == "" || signatureHex == "" {
		return false, "headers obrigatórios"
	}

	deviceID, err := uuid.Parse(deviceIDRaw)
	if err != nil {
		return false, "device id inválido"
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return false, "latitude inválida"
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return false, "longitude inválida"
	}

	if authHeader != "" {
		fmt.Println("Não precisa verificar assinatura pois ja est autenticado")
		verifyAuth(c)
		c.Set("lat", lat)
		c.Set("lon", lon)
		c.Next()
		return true, ""
	}

	message := []byte(
		fmt.Sprintf(
			`{"lat":%.4f,"lon":%.4f}`,
			lat,
			lon,
		),
	)

	publicKey, err := deviceSvc.GetDeviceByID(deviceID)
	if err != nil {
		return false, "device não encontrado"
	}

	if !verifySignature(
		publicKey,
		message,
		signatureHex,
	) {
		return false, "assinatura inválida"
	}

	fmt.Println("assinatura valida")
	c.Set("lat", lat)
	c.Set("lon", lon)

	return true, ""
}

func verifySignature(
	publicKeyHex string,
	message []byte,
	signatureHex string,
) bool {

	pubBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		fmt.Println("Erro public key:", err)
		return false
	}

	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		fmt.Println("Erro signature:", err)
		return false
	}

	if len(pubBytes) != ed25519.PublicKeySize {
		fmt.Println("Public key inválida")
		return false
	}

	if len(sigBytes) != ed25519.SignatureSize {
		fmt.Println("Signature inválida")
		return false
	}

	return ed25519.Verify(
		ed25519.PublicKey(pubBytes),
		message,
		sigBytes,
	)
}
