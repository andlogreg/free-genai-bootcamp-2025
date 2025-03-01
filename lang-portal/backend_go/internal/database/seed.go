package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type seedWord struct {
	Portuguese string `json:"portuguese"`
	English    string `json:"english"`
}

type seedGroup struct {
	Name  string     `json:"name"`
	Words []seedWord `json:"words"`
}

type wordsAndGroupsSeed struct {
	Groups []seedGroup `json:"groups"`
}

type studyActivity struct {
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
}

type studyActivitiesSeed struct {
	StudyActivities []studyActivity `json:"study_activities"`
}

func RunSeed() error {
	db, err := InitDB()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	// Read and parse words and groups
	wordsData, err := ioutil.ReadFile("seeds/words_and_groups.json")
	if err != nil {
		return fmt.Errorf("failed to read words and groups seed file: %v", err)
	}

	var wordsAndGroups wordsAndGroupsSeed
	if err := json.Unmarshal(wordsData, &wordsAndGroups); err != nil {
		return fmt.Errorf("failed to parse words and groups seed data: %v", err)
	}

	// Read and parse study activities
	activitiesData, err := ioutil.ReadFile("seeds/study_activities.json")
	if err != nil {
		return fmt.Errorf("failed to read study activities seed file: %v", err)
	}

	var activities studyActivitiesSeed
	if err := json.Unmarshal(activitiesData, &activities); err != nil {
		return fmt.Errorf("failed to parse study activities seed data: %v", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert study activities
	log.Println("Seeding study activities...")
	for _, activity := range activities.StudyActivities {
		_, err := tx.Exec(`
			INSERT INTO study_activities (name, thumbnail_url, description, created_at)
			VALUES (?, ?, ?, ?)
		`, activity.Name, activity.ThumbnailURL, activity.Description, time.Now())
		if err != nil {
			return fmt.Errorf("failed to insert study activity: %v", err)
		}
	}

	// Insert groups and words
	log.Println("Seeding groups and words...")
	for _, group := range wordsAndGroups.Groups {
		// Insert group
		groupResult, err := tx.Exec(`
			INSERT INTO groups (name, created_at)
			VALUES (?, ?)
		`, group.Name, time.Now())
		if err != nil {
			return fmt.Errorf("failed to insert group: %v", err)
		}

		groupID, err := groupResult.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get group ID: %v", err)
		}

		// Insert words and create word-group associations
		for _, word := range group.Words {
			// Insert word
			wordResult, err := tx.Exec(`
				INSERT INTO words (portuguese, english, created_at)
				VALUES (?, ?, ?)
			`, word.Portuguese, word.English, time.Now())
			if err != nil {
				return fmt.Errorf("failed to insert word: %v", err)
			}

			wordID, err := wordResult.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get word ID: %v", err)
			}

			// Create word-group association
			_, err = tx.Exec(`
				INSERT INTO words_groups (word_id, group_id)
				VALUES (?, ?)
			`, wordID, groupID)
			if err != nil {
				return fmt.Errorf("failed to insert word-group association: %v", err)
			}
		}
	}

	// Create some sample study sessions and word reviews
	log.Println("Creating sample study sessions and reviews...")
	// Get first group and activity IDs
	var firstGroupID, firstActivityID int64
	err = tx.QueryRow("SELECT id FROM groups LIMIT 1").Scan(&firstGroupID)
	if err != nil {
		return fmt.Errorf("failed to get first group ID: %v", err)
	}
	err = tx.QueryRow("SELECT id FROM study_activities LIMIT 1").Scan(&firstActivityID)
	if err != nil {
		return fmt.Errorf("failed to get first activity ID: %v", err)
	}

	// Create a study session
	sessionResult, err := tx.Exec(`
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, ?)
	`, firstGroupID, firstActivityID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create study session: %v", err)
	}

	sessionID, err := sessionResult.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get session ID: %v", err)
	}

	// Add some word reviews
	rows, err := tx.Query(`
		SELECT id FROM words
		WHERE id IN (
			SELECT word_id FROM words_groups
			WHERE group_id = ?
		)
		LIMIT 5
	`, firstGroupID)
	if err != nil {
		return fmt.Errorf("failed to get words for reviews: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var wordID int64
		if err := rows.Scan(&wordID); err != nil {
			return fmt.Errorf("failed to scan word ID: %v", err)
		}

		_, err = tx.Exec(`
			INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
			VALUES (?, ?, ?, ?)
		`, wordID, sessionID, true, time.Now())
		if err != nil {
			return fmt.Errorf("failed to create word review: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Println("Seeding completed successfully")
	return nil
}
