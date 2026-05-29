package sensor

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (r *Repository) CreateSensor(s *Sensor) error {
	return r.db.Create(s).Error
}

func (r *Repository) ListSensors(usuarioID uuid.UUID) ([]Sensor, error) {
	var sensors []Sensor
	err := r.db.Where("usuario_id = ?", usuarioID).Find(&sensors).Error
	return sensors, err
}

func (r *Repository) CreateRecord(dado *SensorDados) error {
	return r.db.Create(dado).Error
}

func (r *Repository) DeleteSensor(sensorID, usuarioID uuid.UUID) error {
	result := r.db.Where("id = ? AND usuario_id = ?", sensorID, usuarioID).Delete(&Sensor{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrSensorNotFound
	}

	return nil
}

func (r *Repository) RenewTokenSensor(sensorID, usuarioID uuid.UUID, new_token string) error {
	result := r.db.Model(&Sensor{}).
		Where("id=? AND usuario_id=?", sensorID, usuarioID).
		Update("token", new_token)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSensorNotFound
	}
	return nil

}
