package repository

import (
	"database/sql"
	"time"

	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
)

type WordRepository struct {
	db *sql.DB
}

func NewWordRepository(db *sql.DB) *WordRepository {
	return &WordRepository{db: db}
}

func (r *WordRepository) GetWord(id int64) (*models.Word, error) {
	word := &models.Word{}
	err := r.db.QueryRow(`
		SELECT id, portuguese, english, created_at
		FROM words
		WHERE id = ?
	`, id).Scan(&word.ID, &word.Portuguese, &word.English, &word.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (r *WordRepository) GetWordWithStats(id int64) (*models.WordWithStats, error) {
	word := &models.WordWithStats{}
	err := r.db.QueryRow(`
		SELECT 
			w.id, w.portuguese, w.english, w.created_at,
			COUNT(CASE WHEN wri.correct = 1 THEN 1 END) as correct_count,
			COUNT(CASE WHEN wri.correct = 0 THEN 1 END) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE w.id = ?
		GROUP BY w.id
	`, id).Scan(
		&word.ID, &word.Portuguese, &word.English, &word.CreatedAt,
		&word.CorrectCount, &word.WrongCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (r *WordRepository) ListWords() ([]*models.Word, error) {
	rows, err := r.db.Query(`SELECT id, portuguese, english, created_at FROM words`)
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
	return words, nil
}

func (r *WordRepository) CountWords() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&count)
	return count, err
}

func (r *WordRepository) UpdateWord(word *models.Word) (*models.Word, error) {
	_, err := r.db.Exec(`
		UPDATE words 
		SET portuguese = ?, english = ? 
		WHERE id = ?
	`, word.Portuguese, word.English, word.ID)
	if err != nil {
		return nil, err
	}

	return r.GetWord(word.ID)
}

func (r *WordRepository) DeleteWord(id int64) error {
	_, err := r.db.Exec(`DELETE FROM words WHERE id = ?`, id)
	return err
}

func (r *WordRepository) CreateWord(word *models.Word) (*models.Word, error) {
	result, err := r.db.Exec(`
		INSERT INTO words (portuguese, english, created_at)
		VALUES (?, ?, ?)
	`, word.Portuguese, word.English, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetWord(id)
}
