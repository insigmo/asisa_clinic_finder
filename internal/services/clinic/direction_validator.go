package clinic

import (
	"context"
	"fmt"
	"time"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
)

type DirectionValidator struct {
	dbManager  *db.Manager
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewDirectionValidator(dbManager *db.Manager) *DirectionValidator {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	return &DirectionValidator{
		dbManager:  dbManager,
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}

func (d *DirectionValidator) TakeCity(userID int64) (string, error) {
	user, err := d.dbManager.GetUser(d.ctx, userID)
	if err != nil || user == nil {
		return "", fmt.Errorf("user %d not found: %v", userID, err)
	}

	return user.City, nil
}

func (d *DirectionValidator) ValidateDirection(direction string) error {
	_, err := d.dbManager.FindMedicalDirection(direction)
	if err != nil {
		return fmt.Errorf("direction not found: %v", err)
	}
	return nil
}

func (d *DirectionValidator) FindSimilarDirections(direction string) ([]string, error) {
	allDirections, err := d.dbManager.GetAllMedicalDirections()
	if err != nil {
		return []string{}, err
	}
	similarDirections := FindSimilar(direction, allDirections)

	return similarDirections, nil
}
