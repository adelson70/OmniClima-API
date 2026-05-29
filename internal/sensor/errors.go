package sensor

import (
	"net/http"

	"OmniClima/internal/platform/apperror"
)

var (
	ErrSensorNotFound = apperror.New(http.StatusNotFound, "Sensor não encontrado")
	ErrSensorExists   = apperror.New(http.StatusConflict, "Sensor já existe")
)
