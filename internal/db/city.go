package db

import (
	"context"
	"fmt"
)

const (
	queryFindCityPostalCodes = `
	SELECT postal_code
	FROM city_postal_codes
	WHERE upper(city) = upper(?)`

	maxPostalCodeLength = 8
)

// FindCity returns postal codes for the given city (case-insensitive).
// Returns an empty slice if nothing is found.
func (m *Manager) FindCity(ctx context.Context, city string) ([]int, error) {
	rows, err := m.client.QueryContext(ctx, queryFindCityPostalCodes, city)
	if err != nil {
		return nil, fmt.Errorf("select city_postal_codes failed: %w", err)
	}
	defer rows.Close()

	postalCodes := make([]int, 0, maxPostalCodeLength)
	for rows.Next() {
		var postalCode int
		if err = rows.Scan(&postalCode); err != nil {
			return nil, fmt.Errorf("scan city_postal_codes row failed: %w", err)
		}
		postalCodes = append(postalCodes, postalCode)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate city_postal_codes rows failed: %w", err)
	}

	return postalCodes, nil
}
