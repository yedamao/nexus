package session

import (
	"context"
	"database/sql"
	"time"
)

// SQLStore implements Store using a database/sql connection.
type SQLStore struct {
	db *sql.DB
}

// NewSQLStore creates a new SQLStore.
func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) Create(ctx context.Context, sess *Session) error {
	query := `
		INSERT INTO sessions (user_id, token, created_at, expires_at)
		VALUES ($1, $2, $3, $4)
	`

	if sess.CreatedAt.IsZero() {
		sess.CreatedAt = time.Now()
	}

	_, err := s.db.ExecContext(ctx, query,
		sess.UserID,
		sess.Token,
		sess.CreatedAt,
		sess.ExpiresAt,
	)

	return err
}

func (s *SQLStore) GetByToken(ctx context.Context, token string) (*Session, error) {
	query := `SELECT id, user_id, token, created_at, expires_at FROM sessions WHERE token = $1`

	row := s.db.QueryRowContext(ctx, query, token)

	var sess Session
	err := row.Scan(
		&sess.ID,
		&sess.UserID,
		&sess.Token,
		&sess.CreatedAt,
		&sess.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	if time.Now().After(sess.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return &sess, nil
}
