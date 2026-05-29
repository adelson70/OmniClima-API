package sensor

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateSensorInput struct {
	UsuarioID uuid.UUID
	Nome      string
	Lat       *float64
	Lon       *float64
}

type ListSensorsInput struct {
	UsuarioID uuid.UUID
}

type SensorInput struct {
	SensorID  uuid.UUID
	UsuarioID uuid.UUID
}

type SensorOutput struct {
	ID    uuid.UUID
	Token string
	Nome  string
	Lat   *float64
	Lon   *float64
}

func (s *Service) CreateSensor(in CreateSensorInput) (SensorOutput, error) {
	token, err := generateSensorToken()

	if err != nil {
		return SensorOutput{}, err
	}

	sensor := &Sensor{
		UsuarioID: in.UsuarioID,
		Lat:       in.Lat,
		Lon:       in.Lon,
		Token:     token,
		Nome:      in.Nome,
	}

	if err := s.repo.CreateSensor(sensor); err != nil {
		return SensorOutput{}, err
	}

	return SensorOutput{
		ID:    sensor.ID,
		Token: sensor.Token,
		Nome:  sensor.Nome,
		Lat:   sensor.Lat,
		Lon:   sensor.Lon,
	}, nil
}

func (s *Service) ListSensors(in ListSensorsInput) ([]SensorOutput, error) {
	sensors, err := s.repo.ListSensors(in.UsuarioID)

	if err != nil {
		return nil, err
	}

	out := make([]SensorOutput, 0, len(sensors))
	for _, sensor := range sensors {
		out = append(out, SensorOutput{
			ID:    sensor.ID,
			Token: sensor.Token,
			Nome:  sensor.Nome,
			Lat:   sensor.Lat,
			Lon:   sensor.Lon,
		})
	}
	return out, nil
}

func (s *Service) DeleteSensor(in SensorInput) error {
	return s.repo.DeleteSensor(in.SensorID, in.UsuarioID)
}

func (s *Service) RenewTokenSensor(in SensorInput) (string, error) {
	new_token, _ := generateSensorToken()
	err := s.repo.RenewTokenSensor(in.SensorID, in.UsuarioID, new_token)

	return new_token, err
}

type ReadingInput struct {
	SensorID uuid.UUID
	Temp     *float64
	Umid     *float64
	Rain     bool
}

func (s *Service) SaveReading(in ReadingInput) error {
	dado := &SensorDados{
		SensorID: in.SensorID,
		Temp:     in.Temp,
		Umid:     in.Umid,
		Rain:     in.Rain,
	}

	return s.repo.CreateRecord(dado)
}

func generateSensorToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("omni_%d_%s", time.Now().Unix(), hex.EncodeToString(b)), nil
}
