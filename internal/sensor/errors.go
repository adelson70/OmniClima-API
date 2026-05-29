package sensor

import "net/http"

type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

var ErrSensorNotFound = &HTTPError{
	Status:  http.StatusNotFound,
	Message: "Sensor não encontrado",
}
