package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	queryAllMedicalDirections = `SELECT name FROM medical_direction`
	queryFindMedicalDirection = `SELECT name FROM medical_direction WHERE upper(reference_name) = upper(?)`
)

// ErrMedicalDirectionNotFound возвращается, когда направление не найдено в базе.
var ErrMedicalDirectionNotFound = errors.New("medical direction not found")

func (m *Manager) GetAllMedicalDirections(ctx context.Context) ([]string, error) {
	rows, err := m.client.QueryContext(ctx, queryAllMedicalDirections)
	if err != nil {
		return nil, fmt.Errorf("query medical_direction: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var res []string
	for rows.Next() {
		var dir string
		if err := rows.Scan(&dir); err != nil {
			return nil, fmt.Errorf("scan medical_direction row: %w", err)
		}
		res = append(res, dir)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate medical_direction rows: %w", err)
	}
	return res, nil
}

func (m *Manager) FindMedicalDirection(ctx context.Context, direction string) (string, error) {
	direction = strings.ToUpper(direction)

	var name string
	err := m.client.QueryRowContext(ctx, queryFindMedicalDirection, direction).Scan(&name)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", ErrMedicalDirectionNotFound
	case err != nil:
		return "", fmt.Errorf("query medical_direction: %w", err)
	}
	return name, nil
}
