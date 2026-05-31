package validator

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	validatorpkg "github.com/go-playground/validator/v10"
)

func ValidateRequest(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		var validationErrors validatorpkg.ValidationErrors

		if errors.As(err, &validationErrors) {
			fields := make([]string, 0, len(validationErrors))

			t := reflect.TypeOf(req)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}

			for _, fieldErr := range validationErrors {
				field, ok := t.FieldByName(fieldErr.Field())
				if !ok {
					fields = append(fields, fieldErr.Field())
					continue
				}

				jsonTag := field.Tag.Get("json")
				jsonName := strings.Split(jsonTag, ",")[0]

				if jsonName == "" {
					jsonName = fieldErr.Field()
				}

				fields = append(fields, jsonName)
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "campo necessario",
				"fields": fields,
			})

			return false
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_json",
		})

		return false
	}

	return true
}
