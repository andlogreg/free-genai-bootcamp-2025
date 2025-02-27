package handlers

import (
	"net/http"
	"strconv"

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
