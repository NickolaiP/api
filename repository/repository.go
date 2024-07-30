package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Message struct {
	ID        int
	Content   string
	Processed bool
	CreatedAt time.Time
}

type MessageRepository interface {
	Save(ctx context.Context, msg Message) error
	GetProcessedCount(ctx context.Context) (int, error)
}

type PostgresMessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &PostgresMessageRepository{db: db}
}

func (r *PostgresMessageRepository) Save(ctx context.Context, msg Message) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO messages (content, processed) VALUES ($1, $2)", msg.Content, msg.Processed)
	return err
}

func (r *PostgresMessageRepository) GetProcessedCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM messages WHERE processed = TRUE").Scan(&count)
	return count, err
}
