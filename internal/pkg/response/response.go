package response

import (
	"net/http"

	"warehouse/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	statusCode := http.StatusBadRequest
	switch code {
	case 401:
		statusCode = http.StatusUnauthorized
	case 403:
		statusCode = http.StatusForbidden
	case 404:
		statusCode = http.StatusNotFound
	case 500:
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
	})
}

func SuccessWithPage(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.CodeSuccess,
		Message: "success",
		Data: PageData{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}
