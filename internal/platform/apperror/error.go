package apperror

import (
	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(status int, message string) *HTTPError {
	return &HTTPError{Status: status, Message: message}
}

func Abort(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
