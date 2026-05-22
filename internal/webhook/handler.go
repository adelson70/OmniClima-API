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
	rg.POST("/:sensor_id", h.HandleWeboHook)
}

func (h *Handler) HandleWeboHook(c *gin.Context) {
	sensor_id_str, _ := c.Get("sensor_id")
	rawPayload, _ := c.Get("payload")

	sensor_id, _ := uuid.Parse(sensor_id_str.(string))
	body := rawPayload.(middleware.WebohhokBody)

	h.sensorSvc.SaveReading(sensor.ReadingInput{
		SensorID: sensor_id,
		Temp:     body.Temp,
		Umid:     body.Umid,
		Rain:     *body.Rain,
	})

	c.JSON(200, gin.H{
		"message": "Sucesso",
	})
}
