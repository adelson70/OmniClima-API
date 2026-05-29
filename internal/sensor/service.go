package sensor

import (
	"OmniClima/internal/platform/apperror"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
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
	UserID uuid.UUID
	Name   string
	Lat    *float64
	Lon    *float64
}

type UpdateSensorInput struct {
	UserID   uuid.UUID
	SensorID uuid.UUID
	Name     *string
	Lat      *float64
	Lon      *float64
}

type ListSensorsInput struct {
	UserID uuid.UUID
}

type SensorInput struct {
	SensorID uuid.UUID
	UserID   uuid.UUID
}

type SensorOutput struct {
	ID    uuid.UUID
	Token string
	Name  string
	Lat   *float64
	Lon   *float64
}

func (s *Service) CreateSensor(in CreateSensorInput) (SensorOutput, error) {
	token, err := generateSensorToken()

	if err != nil {
		return SensorOutput{}, err
	}

	sensor := &Sensor{
		UserID: in.UserID,
		Lat:    in.Lat,
		Lon:    in.Lon,
		Token:  token,
		Name:   in.Name,
	}

	if err := s.repo.CreateSensor(sensor); err != nil {
		return SensorOutput{}, err
	}

	return SensorOutput{
		ID:    sensor.ID,
		Token: sensor.Token,
		Name:  sensor.Name,
		Lat:   sensor.Lat,
		Lon:   sensor.Lon,
	}, nil
}

func (s *Service) UpdateSensor(in UpdateSensorInput) error {

	if !in.hasUpdates() {
		return apperror.New(http.StatusBadRequest, "Nenhum campo para atualizar")
	}

	updates := map[string]interface{}{}

	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return apperror.New(http.StatusBadRequest, "Nome não pode ser vazio")
		}
		updates["name"] = name
	}
	if in.Lat != nil {
		updates["lat"] = *in.Lat
	}
	if in.Lon != nil {
		updates["lon"] = *in.Lon
	}

	if err := s.repo.UpdateSensor(in.SensorID, in.UserID, updates); err != nil {
		return err
	}

	return nil
}

func (s *Service) ListSensors(in ListSensorsInput) ([]SensorOutput, error) {
	sensors, err := s.repo.ListSensors(in.UserID)

	if err != nil {
		return nil, err
	}

	out := make([]SensorOutput, 0, len(sensors))
	for _, sensor := range sensors {
		out = append(out, SensorOutput{
			ID:    sensor.ID,
			Token: sensor.Token,
			Name:  sensor.Name,
			Lat:   sensor.Lat,
			Lon:   sensor.Lon,
		})
	}
	return out, nil
}

func (s *Service) DeleteSensor(in SensorInput) error {
	return s.repo.DeleteSensor(in.SensorID, in.UserID)
}

func (s *Service) RenewTokenSensor(in SensorInput) (string, error) {
	newToken, _ := generateSensorToken()
	err := s.repo.RenewTokenSensor(in.SensorID, in.UserID, newToken)

	return newToken, err
}

type ReadingInput struct {
	SensorID uuid.UUID
	Temp     *float64
	Umid     *float64
	Rain     bool
}

func (s *Service) SaveReading(in ReadingInput) error {
	dado := &SensorData{
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

func (in UpdateSensorInput) hasUpdates() bool {
	return in.Name != nil || in.Lat != nil || in.Lon != nil
}
