package api

import (
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/api/handlers"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	dashboardHandler *handlers.DashboardHandler,
	studyActivityHandler *handlers.StudyActivityHandler,
	wordHandler *handlers.WordHandler,
	groupHandler *handlers.GroupHandler,
) *gin.Engine {
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORS())

	// API routes
	api := router.Group("/api")
	{
		// Dashboard routes
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/last_study_session", dashboardHandler.GetLastStudySession)
			dashboard.GET("/study_progress", dashboardHandler.GetStudyProgress)
			dashboard.GET("/quick-stats", dashboardHandler.GetQuickStats)
		}

		// Study activities routes
		activities := api.Group("/study_activities")
		{
			activities.GET("", studyActivityHandler.ListStudyActivities)
			activities.GET("/:id", studyActivityHandler.GetStudyActivity)
			activities.GET("/:id/study_sessions", studyActivityHandler.GetStudyActivitySessions)
			activities.POST("", studyActivityHandler.CreateStudySession)
		}

		// Study sessions routes
		studySessions := api.Group("/study_sessions")
		{
			studySessions.GET("", studyActivityHandler.ListStudySessions)
		}

		// Words routes
		words := api.Group("/words")
		{
			words.GET("", wordHandler.ListWords)
			words.GET("/:id", wordHandler.GetWord)
			words.POST("", wordHandler.CreateWord)
			words.PUT("/:id", wordHandler.UpdateWord)
			words.DELETE("/:id", wordHandler.DeleteWord)
		}

		// Groups routes
		groups := api.Group("/groups")
		{
			groups.GET("", groupHandler.ListGroups)
			groups.GET("/:id", groupHandler.GetGroup)
			groups.GET("/:id/words", groupHandler.GetGroupWords)
			groups.GET("/:id/study_sessions", groupHandler.GetGroupStudySessions)
			groups.POST("", groupHandler.CreateGroup)
			groups.PUT("/:id", groupHandler.UpdateGroup)
			groups.DELETE("/:id", groupHandler.DeleteGroup)
			groups.POST("/:id/words", groupHandler.AddWordsToGroup)
			groups.DELETE("/:id/words/:word_id", groupHandler.RemoveWordFromGroup)
		}
	}

	return router
}
