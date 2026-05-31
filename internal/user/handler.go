package user

import (
	"OmniClima/internal/platform/apperror"
	"OmniClima/internal/platform/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userSvc *Service
}

func NewHandler(userSvc *Service) *Handler {
	return &Handler{userSvc: userSvc}
}

type createUserReq struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("create", h.Create)
}

func (h *Handler) Create(c *gin.Context) {
	var req createUserReq

	if !validator.ValidateRequest(c, &req) {
		return
	}

	out, err := h.userSvc.CreateUser(CreateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})

	if err != nil {
		apperror.Abort(c, err)
		return
	}

	c.JSON(http.StatusCreated, out)

}
