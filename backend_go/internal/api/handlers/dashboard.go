package handlers

import (
	"net/http"

	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/service"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetLastStudySession(c *gin.Context) {
	session, err := h.dashboardService.GetLastStudySession()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if session == nil {
		utils.RespondWithJSON(c, http.StatusOK, gin.H{
			"message": "No study sessions found",
		})
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, session)
}

func (h *DashboardHandler) GetStudyProgress(c *gin.Context) {
	progress, err := h.dashboardService.GetStudyProgress()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, progress)
}

func (h *DashboardHandler) GetQuickStats(c *gin.Context) {
	stats, err := h.dashboardService.GetQuickStats()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, stats)
}
