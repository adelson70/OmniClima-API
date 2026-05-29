package webhook

import (
	"OmniClima/internal/platform/middleware"
	"OmniClima/internal/sensor"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	sensorSvc *sensor.Service
}

func NewHandler(sensorSvc *sensor.Service) *Handler {
	return &Handler{sensorSvc: sensorSvc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/:sensorID", h.HandleWeboHook)
}

func (h *Handler) HandleWeboHook(c *gin.Context) {
	sensorIDStr, _ := c.Get("sensorID")
	rawPayload, _ := c.Get("payload")

	sensorID, _ := uuid.Parse(sensorIDStr.(string))
	body := rawPayload.(middleware.WebohhokBody)

	h.sensorSvc.SaveReading(sensor.ReadingInput{
		SensorID: sensorID,
		Temp:     body.Temp,
		Umid:     body.Umid,
		Rain:     *body.Rain,
	})

	c.JSON(200, gin.H{
		"message": "Sucesso",
	})
}
