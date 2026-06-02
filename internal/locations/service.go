package locations

import (
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateLocationInput struct {
	UserID  uuid.UUID
	Country string
	State   string
	City    string
	Lat     float64
	Lon     float64
}

type LocationOutput struct {
	ID      uuid.UUID
	Country string
	State   string
	City    string
	Lat     float64
	Lon     float64
}

func (s *Service) CreateLocation(in CreateLocationInput) (LocationOutput, error) {
	location := &Location{
		UserID:  in.UserID,
		Country: in.Country,
		State:   in.State,
		City:    in.City,
		Lat:     in.Lat,
		Lon:     in.Lon,
	}
	if err := s.repo.CreateLocation(location); err != nil {
		return LocationOutput{}, err
	}
	return LocationOutput{
		ID:      location.ID,
		Country: location.Country,
		State:   location.State,
		City:    location.City,
		Lat:     location.Lat,
		Lon:     location.Lon,
	}, nil
}

func (s *Service) GetLocationByUserID(userID uuid.UUID) ([]Location, error) {
	return s.repo.GetLocationByUserID(userID)
}

func (s *Service) DeleteLocation(locationID, userID uuid.UUID) error {
	return s.repo.DeleteLocation(locationID, userID)
}
