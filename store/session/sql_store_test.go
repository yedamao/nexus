package session

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("error closing db: %v", closeErr)
		}
	}()

	store := NewSQLStore(db)
	ctx := context.Background()

	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	sess := &Session{
		UserID:    "user-123",
		Token:     "token-abc",
		CreatedAt: fixedTime,
		ExpiresAt: fixedTime.Add(time.Hour),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sessions (user_id, token, created_at, expires_at) VALUES ($1, $2, $3, $4)`)).
		WithArgs(sess.UserID, sess.Token, sess.CreatedAt, sess.ExpiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.Create(ctx, sess)
	if err != nil {
		t.Errorf("error was not expected while creating session: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("error closing db: %v", closeErr)
		}
	}()

	store := NewSQLStore(db)
	ctx := context.Background()

	token := "token-abc"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "token", "created_at", "expires_at"}).
		AddRow("session-1", "user-123", token, now, now.Add(time.Hour))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, created_at, expires_at FROM sessions WHERE token = $1`)).
		WithArgs(token).
		WillReturnRows(rows)

	sess, err := store.GetByToken(ctx, token)
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if sess == nil {
		t.Errorf("expected session, got nil")
	} else if sess.Token != token {
		t.Errorf("expected token %s, got %s", token, sess.Token)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByTokenExpired(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("error closing db: %v", closeErr)
		}
	}()

	store := NewSQLStore(db)
	ctx := context.Background()

	token := "token-expired"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "token", "created_at", "expires_at"}).
		AddRow("session-2", "user-123", token, now, now.Add(-time.Hour))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, created_at, expires_at FROM sessions WHERE token = $1`)).
		WithArgs(token).
		WillReturnRows(rows)

	_, err = store.GetByToken(ctx, token)
	if err != ErrSessionExpired {
		t.Errorf("expected ErrSessionExpired, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
