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

// Helper function to handle API responses (kept for consistency)
async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.json().catch(() => ({
      message: "An unknown error occurred",
    }));
    throw new Error(error.message || `API Error: ${response.status}`);
  }
  return response.json();
}

// Generic API fetch function with error handling (kept for consistency)
async function apiFetch<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  try {
    const url = `${endpoint}`;
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

// Dashboard APIs
export async function getLastStudySession(): Promise<LastStudySession | null> {
  // In a real implementation, this would fetch from the API
  // For now, we'll return mock data
  return {
    id: 123,
    group_id: 456,
    created_at: new Date().toISOString(),
    study_activity_id: 789,
    group_name: "Basic Greetings",
  };
}

export async function getStudyProgress(): Promise<StudyProgress> {
  // Mock data
  return {
    total_words_studied: 3,
    total_available_words: 124,
  };
}

export async function getQuickStats(): Promise<QuickStats> {
  // Mock data
  return {
    success_rate: 80.0,
    total_study_sessions: 4,
    total_active_groups: 3,
    study_streak_days: 4,
  };
}

// Study Activities APIs
export async function getStudyActivities(): Promise<StudyActivity[]> {
  // Mock data
  return [
    {
      id: 1,
      name: "Vocabulary Quiz",
      thumbnail_url: "https://images.unsplash.com/photo-1486312338219-ce68d2c6f44d",
      description: "Practice your vocabulary with flashcards",
    },
    {
      id: 2,
      name: "Word Matching",
      thumbnail_url: "https://images.unsplash.com/photo-1488590528505-98d2b5aba04b",
      description: "Match Portuguese words with their English translations",
    },
    {
      id: 3,
      name: "Audio Pronunciation",
      thumbnail_url: "https://images.unsplash.com/photo-1473091534298-04dcbce3278c",
      description: "Listen to pronunciation and choose the correct word",
    },
  ];
}

export async function getStudyActivity(id: number): Promise<StudyActivity> {
  // Mock data based on id
  const activities = await getStudyActivities();
  const activity = activities.find((a) => a.id === id);
  
  if (!activity) {
    throw new Error("Study activity not found");
  }
  
  return activity;
}

export async function getStudyActivitySessions(
  id: number,
  page: number = 1
): Promise<PaginatedResponse<StudySession>> {
  // Mock data
  return {
    items: [
      {
        id: 123,
        activity_name: "Vocabulary Quiz",
        group_name: "Basic Greetings",
        start_time: new Date(Date.now() - 3600000).toISOString(),
        end_time: new Date().toISOString(),
        review_items_count: 20,
      },
      {
        id: 124,
        activity_name: "Vocabulary Quiz",
        group_name: "Common Phrases",
        start_time: new Date(Date.now() - 86400000).toISOString(),
        end_time: new Date(Date.now() - 82800000).toISOString(),
        review_items_count: 15,
      },
    ],
    pagination: {
      current_page: page,
      total_pages: 1,
      total_items: 2,
      items_per_page: 100,
    },
  };
}

export async function launchStudyActivity(
  study_activity_id: number,
  group_id: number
): Promise<StudyLaunchResponse> {
  // Mock successful response
  return {
    id: Date.now(),
    group_id,
  };
}

// Words APIs
export async function getWords(
  page: number = 1
): Promise<PaginatedResponse<Word>> {
  // Mock data
  const mockWords: Word[] = [
    { id: 1, portuguese: "olá", english: "hello", correct_count: 5, wrong_count: 2 },
    { id: 2, portuguese: "adeus", english: "goodbye", correct_count: 3, wrong_count: 1 },
    { id: 3, portuguese: "obrigado", english: "thank you", correct_count: 7, wrong_count: 0 },
    { id: 4, portuguese: "desculpe", english: "sorry", correct_count: 2, wrong_count: 3 },
    { id: 5, portuguese: "por favor", english: "please", correct_count: 4, wrong_count: 1 },
  ];
  
  return {
    items: mockWords,
    pagination: {
      current_page: page,
      total_pages: 1,
      total_items: mockWords.length,
      items_per_page: 100,
    },
  };
}

export async function getWord(id: number): Promise<WordDetail> {
  // Mock data based on id
  return {
    id,
    portuguese: "olá",
    english: "hello",
    correct_count: 5,
    wrong_count: 2,
    stats: {
      correct_count: 5,
      wrong_count: 2,
    },
    groups: [
      {
        id: 1,
        name: "Basic Greetings",
      },
      {
        id: 2,
        name: "Common Phrases",
      },
    ],
  };
}

// Groups APIs
export async function getGroups(
  page: number = 1
): Promise<PaginatedResponse<Group>> {
  // Mock data
  const mockGroups: Group[] = [
    { id: 1, name: "Basic Greetings", word_count: 20 },
    { id: 2, name: "Common Phrases", word_count: 30 },
    { id: 3, name: "Food Items", word_count: 45 },
    { id: 4, name: "Travel Vocabulary", word_count: 35 },
    { id: 5, name: "Business Terms", word_count: 25 },
  ];
  
  return {
    items: mockGroups,
    pagination: {
      current_page: page,
      total_pages: 1,
      total_items: mockGroups.length,
      items_per_page: 100,
    },
  };
}

export async function getGroup(id: number): Promise<GroupDetail> {
  // Mock data based on id
  return {
    id,
    name: "Basic Greetings",
    stats: {
      total_word_count: 20,
    },
  };
}

export async function getGroupWords(
  id: number,
  page: number = 1
): Promise<PaginatedResponse<Word>> {
  // Mock data
  return await getWords(page);
}

export async function getGroupStudySessions(
  id: number,
  page: number = 1
): Promise<PaginatedResponse<StudySession>> {
  // Re-use the same mock data
  return await getStudySessions(page);
}

// Study Sessions APIs
export async function getStudySessions(
  page: number = 1
): Promise<PaginatedResponse<StudySession>> {
  // Mock data
  const mockSessions: StudySession[] = [
    {
      id: 123,
      activity_name: "Vocabulary Quiz",
      group_name: "Basic Greetings",
      start_time: new Date(Date.now() - 3600000).toISOString(),
      end_time: new Date().toISOString(),
      review_items_count: 20,
    },
    {
      id: 124,
      activity_name: "Word Matching",
      group_name: "Common Phrases",
      start_time: new Date(Date.now() - 86400000).toISOString(),
      end_time: new Date(Date.now() - 82800000).toISOString(),
      review_items_count: 15,
    },
    {
      id: 125,
      activity_name: "Audio Pronunciation",
      group_name: "Food Items",
      start_time: new Date(Date.now() - 172800000).toISOString(),
      end_time: new Date(Date.now() - 169200000).toISOString(),
      review_items_count: 25,
    },
  ];
  
  return {
    items: mockSessions,
    pagination: {
      current_page: page,
      total_pages: 1,
      total_items: mockSessions.length,
      items_per_page: 100,
    },
  };
}

export async function getStudySession(
  id: number
): Promise<StudySessionDetail> {
  // Mock data based on id
  return {
    id,
    activity_name: "Vocabulary Quiz",
    group_name: "Basic Greetings",
    start_time: new Date(Date.now() - 3600000).toISOString(),
    end_time: new Date().toISOString(),
    review_items_count: 20,
  };
}

export async function getStudySessionWords(
  id: number,
  page: number = 1
): Promise<PaginatedResponse<Word>> {
  // Re-use the same mock data
  return await getWords(page);
}

// Word Review APIs
export async function submitWordReview(
  study_session_id: number,
  word_id: number,
  correct: boolean
): Promise<{ success: boolean }> {
  // Mock successful response
  return {
    success: true,
  };
}

// Settings APIs
export async function resetHistory(): Promise<ResetResponse> {
  // Mock successful response
  return {
    success: true,
    message: "Study history has been reset",
  };
}

export async function resetFull(): Promise<ResetResponse> {
  // Mock successful response
  return {
    success: true,
    message: "System has been fully reset",
  };
} 