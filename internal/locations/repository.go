package locations

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateLocation(location *Location) error {
	return r.db.Create(location).Error
}

func (r *Repository) GetLocationByUserID(userID uuid.UUID) ([]Location, error) {
	var locations []Location
	if err := r.db.Where("user_id = ?", userID).Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []Location{}, err
		}
		return []Location{}, err
	}
	return locations, nil
}

func (r *Repository) DeleteLocation(locationID, userID uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", locationID, userID).Delete(&Location{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrLocationNotFound
	}
	return nil
}
