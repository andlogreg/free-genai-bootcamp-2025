package repository

import (
	"database/sql"
	"time"

	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
)

type StudySessionRepository struct {
	db *sql.DB
}

func (r *StudySessionRepository) GetTotalDistinctWordsStudied() (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(DISTINCT word_id)
		FROM word_review_items
	`).Scan(&count)
	return count, err
}

func (r *StudySessionRepository) GetWordReviewStats() (correct int, total int, err error) {
	err = r.db.QueryRow(`
		SELECT 
			COUNT(CASE WHEN correct = 1 THEN 1 END) as correct,
			COUNT(*) as total
		FROM word_review_items
	`).Scan(&correct, &total)
	return
}

func (r *StudySessionRepository) GetTotalActiveGroups() (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(DISTINCT group_id)
		FROM study_sessions
	`).Scan(&count)
	return count, err
}

func (r *StudySessionRepository) GetStudyStreak() (int, error) {
	var streak int
	err := r.db.QueryRow(`
		WITH RECURSIVE dates AS (
			SELECT date(created_at) as study_date
			FROM study_sessions
			GROUP BY date(created_at)
			ORDER BY study_date DESC
		),
		streak_calc AS (
			SELECT study_date, 1 as streak
			FROM dates
			WHERE study_date = date('now')
			UNION ALL
			SELECT d.study_date, s.streak + 1
			FROM dates d
			JOIN streak_calc s ON date(d.study_date, '+1 day') = date(s.study_date)
		)
		SELECT COALESCE(MAX(streak), 0)
		FROM streak_calc
	`).Scan(&streak)
	return streak, err
}

func NewStudySessionRepository(db *sql.DB) *StudySessionRepository {
	return &StudySessionRepository{db: db}
}

func (r *StudySessionRepository) GetStudySession(id int64) (*models.StudySessionDetail, error) {
	session := &models.StudySessionDetail{}
	var endTime sql.NullTime
	err := r.db.QueryRow(`
		SELECT 
			ss.id, ss.group_id, ss.study_activity_id, ss.created_at,
			sa.name as activity_name,
			g.name as group_name,
			COUNT(wri.id) as review_items_count,
			MAX(wri.created_at) as end_time
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		WHERE ss.id = ?
		GROUP BY ss.id
	`, id).Scan(
		&session.ID, &session.GroupID, &session.StudyActivityID,
		&session.CreatedAt, &session.ActivityName, &session.GroupName,
		&session.ReviewItemsCount, &endTime,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if endTime.Valid {
		session.EndTime = &endTime.Time
	}
	return session, nil
}

func (r *StudySessionRepository) ListStudySessions(offset, limit int) ([]models.StudySessionDetail, error) {
	rows, err := r.db.Query(`
		SELECT 
			ss.id, ss.group_id, ss.study_activity_id, ss.created_at,
			sa.name as activity_name,
			g.name as group_name,
			COUNT(wri.id) as review_items_count,
			MAX(wri.created_at) as end_time
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		GROUP BY ss.id
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.StudySessionDetail
	for rows.Next() {
		var session models.StudySessionDetail
		var endTimeStr sql.NullString
		var groupID, studyActivityID int64
		err := rows.Scan(
			&session.ID, &groupID, &studyActivityID,
			&session.CreatedAt, &session.ActivityName, &session.GroupName,
			&session.ReviewItemsCount, &endTimeStr,
		)
		if err != nil {
			return nil, err
		}
		// Store the IDs in the hidden fields
		session.GroupID = groupID
		session.StudyActivityID = studyActivityID

		// Handle the end time string
		if endTimeStr.Valid {
			endTime, err := time.Parse(time.RFC3339, endTimeStr.String)
			if err == nil {
				session.EndTime = &endTime
			}
		}

		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *StudySessionRepository) CountStudySessions() (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM study_sessions
	`).Scan(&count)
	return count, err
}

func (r *StudySessionRepository) CreateStudySession(session *models.StudySession) error {
	result, err := r.db.Exec(`
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, ?)
	`, session.GroupID, session.StudyActivityID, time.Now())
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	session.ID = id
	return nil
}

func (r *StudySessionRepository) GetLastStudySession() (*models.StudySessionDetail, error) {
	session := &models.StudySessionDetail{}
	// Use NullString instead of NullTime for end_time because:
	// 1. SQLite stores timestamps as strings, not native time types
	// 2. When scanning directly into sql.NullTime, it causes a type mismatch error:
	//    "sql: Scan error on column index 7, name "end_time": unsupported Scan, storing driver.Value type string into type *time.Time"
	// 3. We need to first scan into a string, then parse it into a time.Time
	var endTimeStr sql.NullString
	err := r.db.QueryRow(`
		SELECT 
			ss.id, ss.group_id, ss.study_activity_id, ss.created_at,
			sa.name as activity_name,
			g.name as group_name,
			COUNT(wri.id) as review_items_count,
			MAX(wri.created_at) as end_time
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		GROUP BY ss.id
		ORDER BY ss.created_at DESC
		LIMIT 1
	`).Scan(
		&session.ID, &session.GroupID, &session.StudyActivityID,
		&session.CreatedAt, &session.ActivityName, &session.GroupName,
		&session.ReviewItemsCount, &endTimeStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if endTimeStr.Valid {
		// Parse the string timestamp into a time.Time object
		endTime, err := time.Parse(time.RFC3339, endTimeStr.String)
		if err == nil {
			session.EndTime = &endTime
		}
	}
	return session, nil
}

func (r *StudySessionRepository) GetStudySessionWords(sessionID int64, offset, limit int) ([]models.WordWithStats, error) {
	rows, err := r.db.Query(`
		SELECT 
			w.id, w.portuguese, w.english, w.created_at,
			COUNT(CASE WHEN wri.correct = 1 THEN 1 END) as correct_count,
			COUNT(CASE WHEN wri.correct = 0 THEN 1 END) as wrong_count
		FROM words w
		JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
		GROUP BY w.id
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`, sessionID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.WordWithStats
	for rows.Next() {
		var word models.WordWithStats
		err := rows.Scan(
			&word.ID, &word.Portuguese, &word.English, &word.CreatedAt,
			&word.CorrectCount, &word.WrongCount,
		)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}

func (r *StudySessionRepository) CountStudySessionWords(sessionID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(DISTINCT word_id)
		FROM word_review_items
		WHERE study_session_id = ?
	`, sessionID).Scan(&count)
	return count, err
}
