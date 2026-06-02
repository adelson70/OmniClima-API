package devices

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	deviceSvc *Service
}

func NewHandler(deviceSvc *Service) *Handler {
	return &Handler{deviceSvc: deviceSvc}
}

type createDeviceReq struct {
	PublicKey string `json:"public_key" binding:"required"`
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.Create)
}

func (h *Handler) Create(c *gin.Context) {
	var req createDeviceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "JSON Inválido"})
		return
	}

	id, err := h.deviceSvc.CreateDevice(CreateDeviceInput{
		PublicKey: req.PublicKey,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar dispositivo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"deviceID": id})
}
