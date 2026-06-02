package devices

import "github.com/google/uuid"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateDeviceInput struct {
	PublicKey string
}

func (s *Service) CreateDevice(in CreateDeviceInput) (uuid.UUID, error) {
	device := &Device{
		PublicKey: in.PublicKey,
	}
	id, err := s.repo.CreateDevice(device)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *Service) GetDeviceByID(deviceID uuid.UUID) (string, error) {
	device, err := s.repo.GetDeviceByID(deviceID)
	if err != nil {
		return "", err
	}
	return device, nil
}
