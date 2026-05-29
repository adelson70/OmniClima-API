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

type CreateSensorOutput struct {
	ID    uuid.UUID
	Token string
	Nome  string
	Lat   *float64
	Lon   *float64
}

func (s *Service) CreateSensor(in CreateSensorInput) (CreateSensorOutput, error) {
	token, err := generateSensorToken()

	if err != nil {
		return CreateSensorOutput{}, err
	}

	sensor := &Sensor{
		UsuarioID: in.UsuarioID,
		Lat:       in.Lat,
		Lon:       in.Lon,
		Token:     token,
		Nome:      in.Nome,
	}

	if err := s.repo.CreateSensor(sensor); err != nil {
		return CreateSensorOutput{}, err
	}

	return CreateSensorOutput{
		ID:    sensor.ID,
		Token: sensor.Token,
		Nome:  sensor.Nome,
		Lat:   sensor.Lat,
		Lon:   sensor.Lon,
	}, nil
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
