package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}

func GetPaginationFromContext(c *gin.Context, defaultItemsPerPage int) (int, int) {
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	itemsPerPage := defaultItemsPerPage
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
			itemsPerPage = pp
		}
	}

	return page, itemsPerPage
}

func CalculatePagination(currentPage, itemsPerPage, totalItems int) Pagination {
	totalPages := int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))

	return Pagination{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: itemsPerPage,
	}
}

func CalculateOffset(page, itemsPerPage int) int {
	return (page - 1) * itemsPerPage
}
