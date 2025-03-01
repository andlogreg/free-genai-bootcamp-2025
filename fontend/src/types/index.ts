
export interface Pagination {
  current_page: number;
  total_pages: number;
  total_items: number;
  items_per_page: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  pagination: Pagination;
}

export interface Word {
  id: number;
  portuguese: string;
  english: string;
  correct_count: number;
  wrong_count: number;
}

export interface WordDetail extends Word {
  stats: {
    correct_count: number;
    wrong_count: number;
  };
  groups: {
    id: number;
    name: string;
  }[];
}

export interface Group {
  id: number;
  name: string;
  word_count: number;
}

export interface GroupDetail {
  id: number;
  name: string;
  stats: {
    total_word_count: number;
  };
}

export interface StudyActivity {
  id: number;
  name: string;
  thumbnail_url: string;
  description: string;
}

export interface StudySession {
  id: number;
  activity_name: string;
  group_name: string;
  start_time: string;
  end_time: string;
  review_items_count: number;
}

export interface LastStudySession {
  id: number;
  group_id: number;
  created_at: string;
  study_activity_id: number;
  group_name: string;
}

export interface StudyProgress {
  total_words_studied: number;
  total_available_words: number;
}

export interface QuickStats {
  success_rate: number;
  total_study_sessions: number;
  total_active_groups: number;
  study_streak_days: number;
}

export interface ResetResponse {
  success: boolean;
  message: string;
}

export interface StudyLaunchResponse {
  id: number;
  group_id: number;
}

export interface StudySessionDetail {
  id: number;
  activity_name: string;
  group_name: string;
  start_time: string;
  end_time: string;
  review_items_count: number;
}
