package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	queryGetUser = `
		SELECT id, username, name, lastname, is_bot, city, language_code
		FROM user
		WHERE id = ?`

	queryInsertOrUpdateUser = `
		INSERT INTO user(id, username, name, lastname, is_bot, city, language_code)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			username      = excluded.username,
			name          = excluded.name,
			lastname      = excluded.lastname,
			is_bot        = excluded.is_bot,
			city          = excluded.city,
			language_code = excluded.language_code`
)

// ErrUserNotFound is returned when a user cannot be found by the given ID.
var ErrUserNotFound = errors.New("user not found")

// GetUser returns the user by ID. Returns ErrUserNotFound if no user exists.
func (db *Manager) GetUser(ctx context.Context, userID int64) (*User, error) {
	const op = "get_user"

	var user User
	err := db.client.QueryRowContext(ctx, queryGetUser, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Lastname,
		&user.IsBot,
		&user.City,
		&user.LanguageCode,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrUserNotFound
	case err != nil:
		return nil, fmt.Errorf("%s: scan user row failed: %w", op, err)
	}

	return &user, nil
}

// InsertOrUpdateUser inserts a new user or updates an existing one (upsert by id).
func (db *Manager) InsertOrUpdateUser(ctx context.Context, user *User) error {
	const op = "insert_or_update_user"

	_, err := db.client.ExecContext(
		ctx,
		queryInsertOrUpdateUser,
		user.ID,
		user.Username,
		user.Name,
		user.Lastname,
		user.IsBot,
		user.City,
		user.LanguageCode,
	)
	if err != nil {
		return fmt.Errorf("%s: exec failed: %w", op, err)
	}

	return nil
}

type User struct {
	ID           int64
	Username     string
	Name         string
	Lastname     string
	IsBot        bool
	City         string
	LanguageCode string
}
