package suggestion

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	suggestionSvc *Service
}

func NewHandler(suggestionSvc *Service) *Handler {
	return &Handler{suggestionSvc: suggestionSvc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/:lat/:lon", h.Get)
}

func (h *Handler) Get(c *gin.Context) {
	latitude := c.Param("lat")
	longitude := c.Param("lon")

	fmt.Println("bateu aqui")
	fmt.Println("latitude:", latitude)
	fmt.Println("longitude:", longitude)
}
