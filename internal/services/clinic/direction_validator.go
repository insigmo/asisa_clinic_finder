package clinic

import (
	"context"
	"fmt"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
)

type DirectionValidator struct {
	dbManager *db.Manager
}

func NewDirectionValidator(dbManager *db.Manager) *DirectionValidator {
	return &DirectionValidator{dbManager: dbManager}
}

func (d *DirectionValidator) TakeCity(ctx context.Context, userID int64) (string, error) {
	user, err := d.dbManager.GetUser(ctx, userID)
	if err != nil || user == nil {
		return "", fmt.Errorf("user %d not found: %w", userID, err)
	}
	return user.City, nil
}

func (d *DirectionValidator) ValidateDirection(ctx context.Context, direction string) (string, error) {
	validName, err := d.dbManager.FindMedicalDirection(ctx, direction)
	if err != nil {
		return "", fmt.Errorf("direction not found: %w", err)
	}
	return validName, nil
}

func (d *DirectionValidator) FindSimilarDirections(ctx context.Context, direction string) ([]string, error) {
	allDirections, err := d.dbManager.GetAllMedicalDirections(ctx)
	if err != nil {
		return nil, err
	}
	return FindSimilar(direction, allDirections), nil
}
