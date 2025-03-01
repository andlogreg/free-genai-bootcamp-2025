package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{Error: message})
}

func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

func RespondWithPaginatedJSON(c *gin.Context, items interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Items:      items,
		Pagination: pagination,
	})
}
