package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	queryAllMedicalDirections = `SELECT name FROM medical_direction`
	queryFindMedicalDirection = `SELECT name FROM medical_direction WHERE upper(reference_name) = upper($1)`
)

// ErrMedicalDirectionNotFound is returned when a medical direction
// cannot be found by the given reference name.
var ErrMedicalDirectionNotFound = errors.New("medical direction not found")

// GetAllMedicalDirections returns the list of all medical direction names.
func (db *Manager) GetAllMedicalDirections() ([]string, error) {
	rows, err := db.client.QueryContext(db.ctx, queryAllMedicalDirections)
	if err != nil {
		return nil, fmt.Errorf("query medical_direction failed: %w", err)
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var dir string
		if err := rows.Scan(&dir); err != nil {
			return nil, fmt.Errorf("scan medical_direction row failed: %w", err)
		}
		res = append(res, dir)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate medical_direction rows failed: %w", err)
	}

	return res, nil
}

// FindMedicalDirection looks up a medical direction by its reference name
// (case-insensitive). Returns ErrMedicalDirectionNotFound if nothing matches.
func (db *Manager) FindMedicalDirection(direction string) (string, error) {
	direction = strings.ToUpper(direction)

	var name string
	err := db.client.
		QueryRowContext(db.ctx, queryFindMedicalDirection, direction).
		Scan(&name)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", ErrMedicalDirectionNotFound
	case err != nil:
		return "", fmt.Errorf("query medical_direction failed: %w", err)
	}

	return name, nil
}
