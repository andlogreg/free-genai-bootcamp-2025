package handlers

import (
	"net/http"
	"strconv"

	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/service"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StudyActivityHandler struct {
	studyActivityService *service.StudyActivityService
}

func NewStudyActivityHandler(studyActivityService *service.StudyActivityService) *StudyActivityHandler {
	return &StudyActivityHandler{
		studyActivityService: studyActivityService,
	}
}

func (h *StudyActivityHandler) GetStudyActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	activity, err := h.studyActivityService.GetStudyActivity(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if activity == nil {
		utils.RespondWithError(c, http.StatusNotFound, "Activity not found")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, activity)
}

func (h *StudyActivityHandler) GetStudyActivitySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	page, perPage := utils.GetPaginationFromContext(c, 100)
	sessions, total, err := h.studyActivityService.GetStudyActivitySessions(id, page, perPage)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	utils.RespondWithPaginatedJSON(c, sessions, pagination)
}

type CreateStudySessionRequest struct {
	GroupID         int64 `json:"group_id" binding:"required"`
	StudyActivityID int64 `json:"study_activity_id" binding:"required"`
}

func (h *StudyActivityHandler) CreateStudySession(c *gin.Context) {
	var req CreateStudySessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.studyActivityService.CreateStudySession(req.GroupID, req.StudyActivityID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, session)
}

// ListStudyActivities returns a list of all study activities
func (h *StudyActivityHandler) ListStudyActivities(c *gin.Context) {
	activities, err := h.studyActivityService.ListStudyActivities()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch study activities")
		return
	}
	c.JSON(http.StatusOK, activities)
}

// ListStudySessions returns a paginated list of study sessions
func (h *StudyActivityHandler) ListStudySessions(c *gin.Context) {
	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Calculate offset for database query
	offset := (page - 1) * pageSize

	// Fetch study sessions
	sessions, err := h.studyActivityService.ListStudySessions(offset, pageSize)
	if err != nil {
		// Log the error for debugging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch study sessions: " + err.Error()})
		return
	}

	// Count total study sessions for pagination
	total, err := h.studyActivityService.CountStudySessions()
	if err != nil {
		// Log the error for debugging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count study sessions: " + err.Error()})
		return
	}

	// Calculate total pages
	totalPages := (total + pageSize - 1) / pageSize

	// Return paginated response
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Items: sessions,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   total,
			ItemsPerPage: pageSize,
		},
	})
}
