package middleware

import (
	"OmniClima/internal/sensor"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebohhokBody struct {
	Temp *float64 `json:"temp,omitempty"`
	Umid *float64 `json:"umid,omitempty"`
	Rain *bool    `json:"rain,omitempty"`
}

func bodyVazio(p WebohhokBody) bool {
	return p.Temp != nil || p.Umid != nil || p.Rain != nil
}

func AuthWebhook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload WebohhokBody
		token := c.GetHeader("x-sensor-token")
		sensor_id := c.Param("sensor_id")
		body, _ := c.GetRawData()

		if err := json.Unmarshal(body, &payload); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "json invalido"})
		}

		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "token invalido"})
			return
		}
		if sensor_id == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "sensor invalido"})
			return
		}

		if !bodyVazio(payload) {
			c.AbortWithStatusJSON(400, gin.H{"error": "body vazio"})
			return
		}

		ok, _ := sensor.ValidateToken(db, sensor_id, token)

		if payload.Rain == nil {
			v := false
			payload.Rain = &v
		}

		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "sensor não encontrado"})
			return
		}

		c.Set("sensor_id", sensor_id)
		c.Set("payload", payload)

		c.Next()
	}
}
