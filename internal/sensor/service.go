package sensor

import "github.com/google/uuid"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
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
