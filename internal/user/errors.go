package user

import (
	"net/http"

	"OmniClima/internal/platform/apperror"
)

var (
	ErrUserNotFound = apperror.New(http.StatusNotFound, "Usuário não encontrado")
	ErrUserExists   = apperror.New(http.StatusConflict, "Usuário já existe")
)
