package service

import (
	"log"

	repository "github.com/Informasjonsforvaltning/fdk-user-feedback-service/repository"
)

type ThreadIdService interface {
	GetThreadId(id string) (*string, error)
	CreateThreadId(id string, threadId string) error
}

type ThreadIdServiceImpl struct {
	ThreadIdRepository repository.ThreadIdRepository
}

func (threadIdService *ThreadIdServiceImpl) GetThreadId(id string) (*string, error) {
	threadId, err := threadIdService.ThreadIdRepository.GetThreadId(id)
	if err != nil {
		log.Println("GetThreadIdByEntityId error.\n[ERROR] -", err)
		return nil, err
	}

	return threadId, err
}

func (threadIdService *ThreadIdServiceImpl) CreateThreadId(id string, threadId string) error {
	err := threadIdService.ThreadIdRepository.CreateThreadId(id, threadId)
	if err != nil {
		log.Println("CreateThreadId error.\n[ERROR] -", err)
		return err
	}
	return nil
}

var CurrentThreadIdService ThreadIdService
