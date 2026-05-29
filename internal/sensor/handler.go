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
	Name string   `json:"nome"`
	Lat  *float64 `json:"lat"`
	Lon  *float64 `json:"lon"`
}

func NewHandler(sensorSvc *Service) *Handler {
	return &Handler{sensorSvc: sensorSvc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.Create)
	rg.GET("/", h.List)
	rg.DELETE("/:sensorID", h.Delete)
	rg.POST("/renew/:sensorID", h.Renew)
}

func (h *Handler) Create(c *gin.Context) {
	var req createSensorReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "JSON Inválido"})
		return
	}

	userID := h.userIDFromContext(c)

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome é Obrigatório"})
		return
	}

	out, err := h.sensorSvc.CreateSensor(CreateSensorInput{
		UserID: userID,
		Name:   name,
		Lat:    req.Lat,
		Lon:    req.Lon,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar sensor"})
		return
	}

	c.JSON(http.StatusCreated, out)

}

func (h *Handler) List(c *gin.Context) {
	userID := h.userIDFromContext(c)

	out, err := h.sensorSvc.ListSensors(ListSensorsInput{
		UserID: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar sensores do usuario"})
		return
	}

	c.JSON(http.StatusOK, out)

}

func (h *Handler) Delete(c *gin.Context) {
	userID := h.userIDFromContext(c)
	sensorID := h.sensorIDExtract(c)

	err := h.sensorSvc.DeleteSensor(SensorInput{
		SensorID: sensorID,
		UserID:   userID,
	})

	if err != nil {
		apperror.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sensor deletado com sucesso"})

}

func (h *Handler) Renew(c *gin.Context) {
	userID := h.userIDFromContext(c)
	sensorID := h.sensorIDExtract(c)

	newToken, err := h.sensorSvc.RenewTokenSensor(SensorInput{
		SensorID: sensorID,
		UserID:   userID,
	})

	if err != nil {
		apperror.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token do sensor renovado", "newToken": newToken})
}

func (h *Handler) userIDFromContext(c *gin.Context) uuid.UUID {
	useIDRaw, _ := c.Get("userID")
	userID, _ := uuid.Parse(useIDRaw.(string))
	return userID
}

func (h *Handler) sensorIDExtract(c *gin.Context) uuid.UUID {
	sensorIDRaw, _ := c.Params.Get("sensorID")
	sensorID, _ := uuid.Parse(sensorIDRaw)

	return sensorID

}
