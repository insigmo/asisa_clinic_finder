package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	ID           int64
	Username     string
	Name         string
	Lastname     string
	IsBot        bool
	City         string
	LanguageCode string
	State        string
}

func (d *Manager) GetUser(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT
			id,
			username,
			name,
			lastname,
			is_bot,
			city,
			language_code,
			state
		FROM user
		WHERE id = ?
		LIMIT 1
	`

	row := d.client.QueryRowContext(ctx, query, userID)

	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Lastname,
		&user.IsBot,
		&user.City,
		&user.LanguageCode,
		&user.State,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get user failed: %w", err)
	}

	return user, nil
}

func (d *Manager) InsertOrUpdateUser(ctx context.Context, user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	query := `
		INSERT INTO user (
			id,
			username,
			name,
			lastname,
			is_bot,
			city,
			language_code,
			state
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			username = excluded.username,
			name = excluded.name,
			lastname = excluded.lastname,
			is_bot = excluded.is_bot,
			city = excluded.city,
			language_code = excluded.language_code,
			state = excluded.state
	`

	_, err := d.client.ExecContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Name,
		user.Lastname,
		user.IsBot,
		user.City,
		user.LanguageCode,
		user.State,
	)
	if err != nil {
		return fmt.Errorf("insert or update user failed: %w", err)
	}

	return nil
}
