// This file is now a wrapper around the apiClient
// It re-exports all the functions from the apiClient for backward compatibility

import apiClient from './apiClient';

// Re-export all functions from the apiClient
export const {
  // Dashboard APIs
  getLastStudySession,
  getStudyProgress,
  getQuickStats,
  
  // Study Activities APIs
  getStudyActivities,
  getStudyActivity,
  getStudyActivitySessions,
  launchStudyActivity,
  
  // Words APIs
  getWords,
  getWord,
  
  // Groups APIs
  getGroups,
  getGroup,
  getGroupWords,
  getGroupStudySessions,
  
  // Study Sessions APIs
  getStudySessions,
  getStudySession,
  getStudySessionWords,
  
  // Word Review APIs
  submitWordReview,
  
  // Settings APIs
  resetHistory,
  resetFull,
} = apiClient;

// For backward compatibility, keep the API_BASE_URL constant
export const API_BASE_URL = "http://localhost:3000/api";
