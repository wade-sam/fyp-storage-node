package configuration

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/wade-sam/fypstoragenode/entity"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Config struct {
	ClientName string   `json:"clientname"`
	Policies   []string `json:"policies"`
}
type Service struct {
	repo Repository
}

func NewConfigurationService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}
func (s *Service) SetStorageNode(name string) error {
	return s.repo.SetStorageNode(name)
}

func (s *Service) GetStorageNode() (string, error) {
	name, err := s.repo.GetStorageNode()
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", entity.ErrFieldWasEmpty
	}
	return name, nil
}
