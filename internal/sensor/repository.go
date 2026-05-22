package sensor

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func ValidateToken(db *gorm.DB, sensor_id, token string) (bool, error) {
	var count int64
	err := db.Model(&Sensor{}).
		Where("id = ? AND token = ?", sensor_id, token).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateRecord(dado *SensorDados) error {
	return r.db.Create(dado).Error
}
