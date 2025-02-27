package repository

import (
	"database/sql"
	"time"

	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
)

type StudyActivityRepository struct {
	db *sql.DB
}

func NewStudyActivityRepository(db *sql.DB) *StudyActivityRepository {
	return &StudyActivityRepository{db: db}
}

func (r *StudyActivityRepository) GetStudyActivity(id int64) (*models.StudyActivity, error) {
	activity := &models.StudyActivity{}
	err := r.db.QueryRow(`
		SELECT id, name, thumbnail_url, description, created_at
		FROM study_activities
		WHERE id = ?
	`, id).Scan(&activity.ID, &activity.Name, &activity.ThumbnailURL, &activity.Description, &activity.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *StudyActivityRepository) ListStudyActivities() ([]models.StudyActivity, error) {
	rows, err := r.db.Query(`
		SELECT id, name, thumbnail_url, description, created_at
		FROM study_activities
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.StudyActivity
	for rows.Next() {
		var activity models.StudyActivity
		err := rows.Scan(
			&activity.ID, &activity.Name, &activity.ThumbnailURL,
			&activity.Description, &activity.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func (r *StudyActivityRepository) CreateStudyActivity(activity *models.StudyActivity) error {
	result, err := r.db.Exec(`
		INSERT INTO study_activities (name, thumbnail_url, description, created_at)
		VALUES (?, ?, ?, ?)
	`, activity.Name, activity.ThumbnailURL, activity.Description, time.Now())
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	activity.ID = id
	return nil
}

func (r *StudyActivityRepository) GetStudyActivitySessions(activityID int64, offset, limit int) ([]models.StudySessionDetail, error) {
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
		WHERE ss.study_activity_id = ?
		GROUP BY ss.id
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`, activityID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.StudySessionDetail
	for rows.Next() {
		var session models.StudySessionDetail
		var endTime sql.NullTime
		err := rows.Scan(
			&session.ID, &session.GroupID, &session.StudyActivityID,
			&session.CreatedAt, &session.ActivityName, &session.GroupName,
			&session.ReviewItemsCount, &endTime,
		)
		if err != nil {
			return nil, err
		}
		if endTime.Valid {
			session.EndTime = &endTime.Time
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *StudyActivityRepository) CountStudyActivitySessions(activityID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM study_sessions
		WHERE study_activity_id = ?
	`, activityID).Scan(&count)
	return count, err
}
