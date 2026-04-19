package db

import (
	"context"
	"database/sql"
	"fmt"
)

func (db *Manager) GetUser(ctx context.Context, userID int) (*User, error) {
	var err error
	stmt, err := db.client.PrepareContext(ctx, `
		SELECT username, name, lastname, is_bot, city, language_code FROM user where id = ?
	`)

	defer func(stmt *sql.Stmt) {
		stmtErr := stmt.Close()
		if stmtErr != nil {
			err = stmtErr
		}
	}(stmt)

	if err != nil {
		return nil, fmt.Errorf("cannot prepare select request")
	}
	res, err := stmt.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("cannot select data, %v", err)
	}

	var user User
	err = res.Scan(&user.Username, &user.Name, &user.Lastname, &user.IsBot, &user.City, &user.LanguageCode)
	if err != nil {
		return nil, fmt.Errorf("failed when scan fields: %v", err)
	}
	return &user, nil
}

func (db *Manager) InsertOrUpdateUser(ctx context.Context, user *User) error {
	const op = "insert_or_update"
	query := `
		INSERT INTO user(id, username, name, lastname, is_bot, city, language_code) 
		VALUES (?, ?, ?, ?, ?, ?, ? )
		ON CONFLICT (id) DO 
		UPDATE SET 
		   username = excluded.username, 
		   name = excluded.name, 
		   lastname = excluded.lastname, 
		   is_bot = excluded.is_bot, 
		   city = excluded.city, 
		   language_code = excluded.language_code
	`
	stmt, err := db.client.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: failed on prepare insert or update query: %v", op, err)
	}
	_, err = stmt.ExecContext(
		ctx,
		user.ID,
		user.Username,
		user.Name,
		user.Lastname,
		user.IsBot,
		user.City,
		user.LanguageCode,
	)
	if err != nil {
		return fmt.Errorf("%s: failed %v", op, err)
	}
	return nil
}
