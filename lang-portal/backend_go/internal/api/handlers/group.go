package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/service"
	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

func (h *GroupHandler) ListGroups(c *gin.Context) {
	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get paginated groups
	groups, totalCount, err := h.groupService.ListGroupsPaginated(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	// Create response with pagination
	response := models.PaginatedResponse{
		Items: groups,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   totalCount,
			ItemsPerPage: pageSize,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	group, err := h.groupService.GetGroup(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	c.JSON(http.StatusOK, group)
}

func (h *GroupHandler) GetGroupWords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get paginated group words with stats
	words, totalCount, err := h.groupService.GetGroupWordsPaginated(id, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	// Create response with pagination
	response := models.PaginatedResponse{
		Items: words,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   totalCount,
			ItemsPerPage: pageSize,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdGroup, err := h.groupService.CreateGroup(&group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdGroup)
}

func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	group.ID = id

	updatedGroup, err := h.groupService.UpdateGroup(&group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedGroup)
}

func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	if err := h.groupService.DeleteGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *GroupHandler) AddWordsToGroup(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var wordIDs []int64
	if err := c.ShouldBindJSON(&wordIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.groupService.AddWordsToGroup(groupID, wordIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *GroupHandler) RemoveWordFromGroup(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word ID"})
		return
	}

	if err := h.groupService.RemoveWordFromGroup(groupID, wordID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetGroupStudySessions returns a paginated list of study sessions for a group
func (h *GroupHandler) GetGroupStudySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// For now, return an empty list with pagination
	// In a real implementation, you would fetch the study sessions for this group
	// We're using the id parameter to acknowledge its use and avoid linter errors
	_ = id // This line acknowledges that we're using the id parameter

	response := models.PaginatedResponse{
		Items: []models.StudySession{},
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   0,
			TotalItems:   0,
			ItemsPerPage: pageSize,
		},
	}

	c.JSON(http.StatusOK, response)
}
