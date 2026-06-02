package locations

import (
	"net/http"

	"OmniClima/internal/platform/apperror"
)

var (
	ErrLocationNotFound = apperror.New(http.StatusNotFound, "Nenhuma localização encontrada")
	ErrLocationExists   = apperror.New(http.StatusConflict, "Localização já existe")
)
