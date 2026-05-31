package user

import (
	"OmniClima/internal/platform/apperror"
	"OmniClima/internal/platform/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	rg.DELETE("/", h.Delete)
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

func (h *Handler) Delete(c *gin.Context) {
	useIDRaw, _ := c.Get("userID")
	userID, _ := uuid.Parse(useIDRaw.(string))

	err := h.userSvc.DeleteUser(userID)

	if err != nil {
		apperror.Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário deletado com sucesso"})
}
