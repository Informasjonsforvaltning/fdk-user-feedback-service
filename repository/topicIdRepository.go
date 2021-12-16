package repository

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ThreadIdRepository interface {
	GetThreadId(id string) (*string, error)
	CreateThreadId(id string, threadId string) error
}

type ThreadIdRepositoryImpl struct {
	FirestoreProjectId    string
	FirestoreCollectionId string
}

func (threadRepository *ThreadIdRepositoryImpl) GetThreadId(id string) (*string, error) {
	ctx := context.Background()

	firestoreClient, err := firestore.NewClient(ctx, threadRepository.FirestoreProjectId)
	if err != nil {
		return nil, err
	}
	defer firestoreClient.Close()

	dataSnapshot, err := firestoreClient.Collection(threadRepository.FirestoreCollectionId).Doc(id).Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	threadIdResponse, err := dataSnapshot.DataAt("topicId")
	if err != nil {
		return nil, err
	}

	threadIdString := fmt.Sprint(threadIdResponse)

	return &threadIdString, err
}

func (threadRepository *ThreadIdRepositoryImpl) CreateThreadId(id string, threadId string) error {
	ctx := context.Background()

	firestoreClient, err := firestore.NewClient(ctx, threadRepository.FirestoreProjectId)
	if err != nil {
		return err
	}
	defer firestoreClient.Close()

	_, err = firestoreClient.Collection(threadRepository.FirestoreCollectionId).Doc(id).Set(ctx, map[string]interface{}{
		"topicId": threadId,
	})

	return err
}

var CurrentThreadIdRepository ThreadIdRepository
