import { toast } from "sonner";
import {
  Group,
  GroupDetail,
  LastStudySession,
  PaginatedResponse,
  QuickStats,
  ResetResponse,
  StudyActivity,
  StudyLaunchResponse,
  StudyProgress,
  StudySession,
  StudySessionDetail,
  Word,
  WordDetail,
} from "@/types";

// Import configuration
import config from "./config";

// Import the mock API implementation
// Using dynamic import to avoid TypeScript errors
let mockApi: any = null;
try {
  // This is a workaround for TypeScript not finding the module
  // In a real application, you would properly set up the module resolution
  mockApi = require("./mockApi");
} catch (error) {
  console.error("Error importing mockApi:", error);
}

// Helper function to handle API responses
async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.json().catch(() => ({
      message: "An unknown error occurred",
    }));
    throw new Error(error.message || `API Error: ${response.status}`);
  }
  return response.json();
}

// Generic API fetch function with error handling
async function apiFetch<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  try {
    const url = `${config.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options?.headers,
      },
    });
    return handleResponse<T>(response);
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error occurred";
    toast.error(errorMessage);
    throw error;
  }
}

// API client factory that decides whether to use mock or real implementation
function createApiClient() {
  // If mock API is enabled and mockApi is available, return the mock implementation
  if (config.useMockApi && mockApi) {
    console.log("Using mock API implementation");
    return mockApi;
  }

  // If mock API is enabled but mockApi is not available, log a warning
  if (config.useMockApi && !mockApi) {
    console.warn("Mock API is enabled but the implementation is not available. Using real API instead.");
  }

  // Otherwise, return the real API implementation
  console.log("Using real API implementation");
  return {
    // Dashboard APIs
    getLastStudySession: async (): Promise<LastStudySession | null> => {
      try {
        // Fetch the last study session from the API
        const session = await apiFetch<any>("/dashboard/last_study_session");
        
        // Transform the response to match the expected format
        return {
          id: session.id,
          group_id: session.group_id,
          created_at: session.created_at || new Date().toISOString(),
          study_activity_id: session.study_activity_id,
          group_name: session.group_name || "Unknown Group",
        };
      } catch (error) {
        console.error("Error fetching last study session:", error);
        return null;
      }
    },

    getStudyProgress: async (): Promise<StudyProgress> => {
      try {
        // Fetch the study progress from the API
        const progress = await apiFetch<any>("/dashboard/study_progress");
        
        // Transform the response to match the expected format
        return {
          total_words_studied: progress.total_words_studied || 0,
          total_available_words: progress.total_available_words || 0,
        };
      } catch (error) {
        console.error("Error fetching study progress:", error);
        // Return default values if there's an error
        return {
          total_words_studied: 0,
          total_available_words: 0,
        };
      }
    },

    getQuickStats: async (): Promise<QuickStats> => {
      try {
        // Fetch the quick stats from the API
        const stats = await apiFetch<any>("/dashboard/quick-stats");
        
        // Transform the response to match the expected format
        return {
          success_rate: stats.success_rate || 0,
          total_study_sessions: stats.total_study_sessions || 0,
          total_active_groups: stats.total_active_groups || 0,
          study_streak_days: stats.study_streak_days || 0,
        };
      } catch (error) {
        console.error("Error fetching quick stats:", error);
        // Return default values if there's an error
        return {
          success_rate: 0,
          total_study_sessions: 0,
          total_active_groups: 0,
          study_streak_days: 0,
        };
      }
    },

    // Study Activities APIs
    getStudyActivities: async (): Promise<StudyActivity[]> => {
      // Fetch the array of study activities from the API
      const activities = await apiFetch<any[]>("/study_activities");
      
      // Transform the response to match the expected format
      return activities.map(activity => ({
        id: activity.id,
        name: activity.name || "Unknown Activity",
        thumbnail_url: activity.thumbnail_url || "https://images.unsplash.com/photo-1486312338219-ce68d2c6f44d",
        description: activity.description || "No description available",
      }));
    },

    getStudyActivity: async (id: number): Promise<StudyActivity> => {
      // Fetch the study activity from the API
      const activity = await apiFetch<any>(`/study_activities/${id}`);
      
      // Transform the response to match the expected format
      return {
        id: activity.id,
        name: activity.name || "Unknown Activity",
        thumbnail_url: activity.thumbnail_url || "https://images.unsplash.com/photo-1486312338219-ce68d2c6f44d",
        description: activity.description || "No description available",
      };
    },

    getStudyActivitySessions: async (
      id: number,
      page: number = 1
    ): Promise<PaginatedResponse<StudySession>> => {
      // Fetch the paginated study sessions from the API
      const response = await apiFetch<PaginatedResponse<any>>(`/study_activities/${id}/study_sessions?page=${page}`);
      
      // Transform the items to match the expected format
      const transformedSessions: StudySession[] = response.items.map(session => ({
        id: session.id,
        activity_name: session.activity_name || "Unknown Activity",
        group_name: session.group_name || "Unknown Group",
        start_time: session.start_time || session.created_at || new Date().toISOString(),
        end_time: session.end_time || new Date().toISOString(),
        review_items_count: session.review_items_count || 0,
      }));
      
      // Return the response with transformed items
      return {
        items: transformedSessions,
        pagination: response.pagination,
      };
    },

    launchStudyActivity: async (
      study_activity_id: number,
      group_id: number
    ): Promise<StudyLaunchResponse> => {
      return apiFetch<StudyLaunchResponse>("/study_activities", {
        method: "POST",
        body: JSON.stringify({ study_activity_id, group_id }),
      });
    },

    // Words APIs
    getWords: async (
      page: number = 1
    ): Promise<PaginatedResponse<Word>> => {
      // Fetch the paginated words from the API
      // The API now returns the expected format with pagination
      return apiFetch<PaginatedResponse<Word>>(`/words?page=${page}`);
    },

    getWord: async (id: number): Promise<WordDetail> => {
      // Fetch the word from the API
      // The API now returns the word with stats
      return apiFetch<WordDetail>(`/words/${id}`);
    },

    createWord: async (word: Partial<Word>): Promise<Word> => {
      return apiFetch<Word>("/words", {
        method: "POST",
        body: JSON.stringify(word),
      });
    },

    updateWord: async (id: number, word: Partial<Word>): Promise<Word> => {
      return apiFetch<Word>(`/words/${id}`, {
        method: "PUT",
        body: JSON.stringify(word),
      });
    },

    deleteWord: async (id: number): Promise<void> => {
      await apiFetch<void>(`/words/${id}`, {
        method: "DELETE",
      });
    },

    // Groups APIs
    getGroups: async (
      page: number = 1
    ): Promise<PaginatedResponse<Group>> => {
      // Fetch the paginated groups from the API
      // The API now returns the expected format with pagination
      return apiFetch<PaginatedResponse<Group>>(`/groups?page=${page}`);
    },

    getGroup: async (id: number): Promise<GroupDetail> => {
      return apiFetch<GroupDetail>(`/groups/${id}`);
    },

    getGroupWords: async (
      id: number,
      page: number = 1
    ): Promise<PaginatedResponse<Word>> => {
      // Fetch the paginated group words from the API
      // The API now returns the expected format with pagination
      return apiFetch<PaginatedResponse<Word>>(`/groups/${id}/words?page=${page}`);
    },

    getGroupStudySessions: async (
      id: number,
      page: number = 1
    ): Promise<PaginatedResponse<StudySession>> => {
      // Fetch the paginated study sessions for a group from the API
      return apiFetch<PaginatedResponse<StudySession>>(`/groups/${id}/study_sessions?page=${page}`);
    },

    createGroup: async (group: Partial<Group>): Promise<Group> => {
      return apiFetch<Group>("/groups", {
        method: "POST",
        body: JSON.stringify(group),
      });
    },

    updateGroup: async (id: number, group: Partial<Group>): Promise<Group> => {
      return apiFetch<Group>(`/groups/${id}`, {
        method: "PUT",
        body: JSON.stringify(group),
      });
    },

    deleteGroup: async (id: number): Promise<void> => {
      await apiFetch<void>(`/groups/${id}`, {
        method: "DELETE",
      });
    },

    addWordsToGroup: async (
      groupId: number,
      wordIds: number[]
    ): Promise<void> => {
      await apiFetch<void>(`/groups/${groupId}/words`, {
        method: "POST",
        body: JSON.stringify(wordIds),
      });
    },

    removeWordFromGroup: async (
      groupId: number,
      wordId: number
    ): Promise<void> => {
      await apiFetch<void>(`/groups/${groupId}/words/${wordId}`, {
        method: "DELETE",
      });
    },

    // Study Session APIs
    getStudySession: async (id: number): Promise<StudySessionDetail> => {
      return apiFetch<StudySessionDetail>(`/study_sessions/${id}`);
    },

    getStudySessions: async (
      page: number = 1
    ): Promise<PaginatedResponse<StudySession>> => {
      // Fetch the paginated study sessions from the API
      return apiFetch<PaginatedResponse<StudySession>>(`/study_sessions?page=${page}`);
    },

    resetStudySession: async (id: number): Promise<ResetResponse> => {
      return apiFetch<ResetResponse>(`/study_sessions/${id}/reset`, {
        method: "POST",
      });
    },

    submitAnswer: async (
      sessionId: number,
      wordId: number,
      isCorrect: boolean
    ): Promise<void> => {
      await apiFetch<void>(`/study_sessions/${sessionId}/submit`, {
        method: "POST",
        body: JSON.stringify({
          word_id: wordId,
          correct: isCorrect,
        }),
      });
    },
  };
}

export default createApiClient(); 