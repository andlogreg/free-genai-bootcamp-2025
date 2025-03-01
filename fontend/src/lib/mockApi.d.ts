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

// Dashboard APIs
export function getLastStudySession(): Promise<LastStudySession | null>;
export function getStudyProgress(): Promise<StudyProgress>;
export function getQuickStats(): Promise<QuickStats>;

// Study Activities APIs
export function getStudyActivities(): Promise<StudyActivity[]>;
export function getStudyActivity(id: number): Promise<StudyActivity>;
export function getStudyActivitySessions(
  id: number,
  page?: number
): Promise<PaginatedResponse<StudySession>>;
export function launchStudyActivity(
  study_activity_id: number,
  group_id: number
): Promise<StudyLaunchResponse>;

// Words APIs
export function getWords(
  page?: number
): Promise<PaginatedResponse<Word>>;
export function getWord(id: number): Promise<WordDetail>;

// Groups APIs
export function getGroups(
  page?: number
): Promise<PaginatedResponse<Group>>;
export function getGroup(id: number): Promise<GroupDetail>;
export function getGroupWords(
  id: number,
  page?: number
): Promise<PaginatedResponse<Word>>;
export function getGroupStudySessions(
  id: number,
  page?: number
): Promise<PaginatedResponse<StudySession>>;

// Study Sessions APIs
export function getStudySessions(
  page?: number
): Promise<PaginatedResponse<StudySession>>;
export function getStudySession(
  id: number
): Promise<StudySessionDetail>;
export function getStudySessionWords(
  id: number,
  page?: number
): Promise<PaginatedResponse<Word>>;

// Word Review APIs
export function submitWordReview(
  study_session_id: number,
  word_id: number,
  correct: boolean
): Promise<{ success: boolean }>;

// Settings APIs
export function resetHistory(): Promise<ResetResponse>;
export function resetFull(): Promise<ResetResponse>; 