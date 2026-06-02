package locations

import (
	"OmniClima/internal/platform/apperror"
	"OmniClima/internal/platform/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	locationSvc *Service
}

func NewHandler(locationSvc *Service) *Handler {
	return &Handler{locationSvc: locationSvc}
}

type createLocationReq struct {
	Country string  `json:"country" binding:"required"`
	State   string  `json:"state" binding:"required"`
	City    string  `json:"city" binding:"required"`
	Lat     float64 `json:"lat" binding:"required"`
	Lon     float64 `json:"lon" binding:"required"`
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/", h.Get)
	rg.POST("/", h.Create)
	rg.DELETE("/:locationID", h.Delete)
}

func (h *Handler) Get(c *gin.Context) {
	userID := h.userIDFromContext(c)

	out, err := h.locationSvc.GetLocationByUserID(userID)

	if err != nil {
		apperror.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, out)

}

func (h *Handler) Create(c *gin.Context) {
	var req createLocationReq
	userID := h.userIDFromContext(c)

	if !validator.ValidateRequest(c, &req) {
		return
	}

	out, err := h.locationSvc.CreateLocation(CreateLocationInput{
		UserID:  userID,
		Country: req.Country,
		State:   req.State,
		City:    req.City,
		Lat:     req.Lat,
		Lon:     req.Lon,
	})

	if err != nil {
		apperror.Abort(c, err)
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := h.userIDFromContext(c)
	locationIDRaw, _ := c.Params.Get("locationID")
	locationID, _ := uuid.Parse(locationIDRaw)

	err := h.locationSvc.DeleteLocation(locationID, userID)

	if err != nil {
		apperror.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Localização deletada com sucesso"})
}
func (h *Handler) userIDFromContext(c *gin.Context) uuid.UUID {
	userIDRaw, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDRaw.(string))
	return userID
}
