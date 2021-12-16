package service

import (
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/repository"
)

type EntityService interface {
	GetEntity(id string) (*model.Entity, error)
}

type EntityServiceImpl struct {
	EntityRepository repository.EntityRepository
}

func (entityService *EntityServiceImpl) GetEntity(id string) (*model.Entity, error) {
	entity, err := entityService.EntityRepository.GetEntityById(id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

var CurrentEntityService EntityService
