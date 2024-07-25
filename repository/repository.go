package repository

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Message struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Processed bool      `json:"processed"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageRepository struct {
	db *sql.DB
}

func NewPostgresDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(msg *Message) error {
	query := `INSERT INTO messages (content, processed, created_at) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, msg.Content, msg.Processed, time.Now()).Scan(&msg.ID)
}

func (r *MessageRepository) GetStats() (map[string]int, error) {
	stats := make(map[string]int)
	query := `SELECT processed, COUNT(*) FROM messages GROUP BY processed`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var processed bool
		var count int
		if err := rows.Scan(&processed, &count); err != nil {
			return nil, err
		}
		if processed {
			stats["processed"] = count
		} else {
			stats["unprocessed"] = count
		}
	}

	return stats, nil
}
