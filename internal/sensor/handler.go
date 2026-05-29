package sensor

import (
	"OmniClima/internal/platform/apperror"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	sensorSvc *Service
}

type createSensorReq struct {
	Nome string   `json:"nome"`
	Lat  *float64 `json:"lat"`
	Lon  *float64 `json:"lon"`
}

func NewHandler(sensorSvc *Service) *Handler {
	return &Handler{sensorSvc: sensorSvc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.Create)
	rg.GET("/", h.List)
	rg.DELETE("/:sensor_id", h.Delete)
}

func (h *Handler) Create(c *gin.Context) {
	var req createSensorReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "JSON Inválido"})
		return
	}

	user_id := h.userIDFromContext(c)

	nome := strings.TrimSpace(req.Nome)
	if nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome é Obrigatório"})
		return
	}

	out, err := h.sensorSvc.CreateSensor(CreateSensorInput{
		UsuarioID: user_id,
		Nome:      nome,
		Lat:       req.Lat,
		Lon:       req.Lon,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar sensor"})
		return
	}

	c.JSON(http.StatusCreated, out)

}

func (h *Handler) List(c *gin.Context) {
	user_id := h.userIDFromContext(c)

	out, err := h.sensorSvc.ListSensors(ListSensorsInput{
		UsuarioID: user_id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar sensores do usuario"})
		return
	}

	c.JSON(http.StatusOK, out)

}

func (h *Handler) Delete(c *gin.Context) {
	user_id := h.userIDFromContext(c)
	sensor_id_raw, _ := c.Params.Get("sensor_id")
	sensor_id, _ := uuid.Parse(sensor_id_raw)

	err := h.sensorSvc.DeleteSensor(DeleteSensorInput{
		SensorID:  sensor_id,
		UsuarioID: user_id,
	})

	if err != nil {
		apperror.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sensor deletado com sucesso"})

}

func (h *Handler) userIDFromContext(c *gin.Context) uuid.UUID {
	user_id_raw, _ := c.Get("user_id")
	user_id, _ := uuid.Parse(user_id_raw.(string))
	return user_id
}
