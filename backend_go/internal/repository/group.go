package repository

import (
	"database/sql"
	"time"

	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) GetGroup(id int64) (*models.Group, error) {
	group := &models.Group{}
	err := r.db.QueryRow(`
		SELECT id, name, created_at
		FROM groups
		WHERE id = ?
	`, id).Scan(&group.ID, &group.Name, &group.CreatedAt)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *GroupRepository) GetGroupWithStats(id int64) (*models.GroupWithStats, error) {
	group := &models.GroupWithStats{}
	err := r.db.QueryRow(`
		SELECT 
			g.id, g.name, g.created_at,
			COUNT(DISTINCT wg.word_id) as word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		WHERE g.id = ?
		GROUP BY g.id
	`, id).Scan(
		&group.ID, &group.Name, &group.CreatedAt,
		&group.WordCount,
	)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *GroupRepository) ListGroups() ([]*models.Group, error) {
	rows, err := r.db.Query(`SELECT id, name, created_at FROM groups`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		group := &models.Group{}
		if err := rows.Scan(&group.ID, &group.Name, &group.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, rows.Err()
}

func (r *GroupRepository) CountGroups() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&count)
	return count, err
}

func (r *GroupRepository) CreateGroup(group *models.Group) (*models.Group, error) {
	result, err := r.db.Exec(`
		INSERT INTO groups (name, created_at)
		VALUES (?, ?)
	`, group.Name, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetGroup(id)
}

func (r *GroupRepository) GetGroupWords(groupID int64) ([]*models.Word, error) {
	rows, err := r.db.Query(`
		SELECT w.id, w.portuguese, w.english, w.created_at
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY w.id
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*models.Word
	for rows.Next() {
		word := &models.Word{}
		if err := rows.Scan(&word.ID, &word.Portuguese, &word.English, &word.CreatedAt); err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, rows.Err()
}

func (r *GroupRepository) CountGroupWords(groupID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM words_groups
		WHERE group_id = ?
	`, groupID).Scan(&count)
	return count, err
}

func (r *GroupRepository) UpdateGroup(group *models.Group) (*models.Group, error) {
	_, err := r.db.Exec(`
		UPDATE groups 
		SET name = ? 
		WHERE id = ?
	`, group.Name, group.ID)
	if err != nil {
		return nil, err
	}

	return r.GetGroup(group.ID)
}

func (r *GroupRepository) DeleteGroup(id int64) error {
	// Start a transaction since we need to delete from multiple tables
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Delete from words_groups first (due to foreign key constraint)
	_, err = tx.Exec(`DELETE FROM words_groups WHERE group_id = ?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Then delete from groups
	_, err = tx.Exec(`DELETE FROM groups WHERE id = ?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *GroupRepository) AddWordsToGroup(groupID int64, wordIDs []int64) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Insert each word-group association
	for _, wordID := range wordIDs {
		_, err = tx.Exec(`
			INSERT INTO words_groups (word_id, group_id)
			VALUES (?, ?)
		`, wordID, groupID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *GroupRepository) RemoveWordFromGroup(groupID, wordID int64) error {
	_, err := r.db.Exec(`
		DELETE FROM words_groups 
		WHERE group_id = ? AND word_id = ?
	`, groupID, wordID)
	return err
}
