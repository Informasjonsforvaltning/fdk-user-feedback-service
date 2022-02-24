package integration_tests

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
)

type MockEntityRepository struct {
	EntityMap map[string]model.Entity
}

func (m *MockEntityRepository) GetEntityById(entityID string) (*model.Entity, error) {
	entity, present := m.EntityMap[entityID]
	if !present {
		return nil, errors.New("no entity found")
	}
	return &entity, nil
}

type MockThreadIdRepository struct {
	ThreadIdMap map[string]string
}

func (m *MockThreadIdRepository) GetThreadId(id string) (*string, error) {
	threadId, present := m.ThreadIdMap[id]
	if !present {
		return nil, nil
	}
	return &threadId, nil
}

func (m *MockThreadIdRepository) CreateThreadId(id string, threadId string) error {
	m.ThreadIdMap[id] = threadId
	return nil
}

type MockThreadRepository struct {
	ThreadMap map[string]*model.Thread
}

func (m *MockThreadRepository) GetThread(threadId string, page string) (*model.Thread, error) {
	var pagedId string
	if page == "1" || page == "2" {
		pagedId = threadId
	} else {
		pagedId = "99"
	}

	thread, present := m.ThreadMap[pagedId]
	if !present {
		return nil, errors.New("no thread found")
	}
	return thread, nil
}
func (m *MockThreadRepository) CreateThread(thread model.Thread) (*model.Thread, error) {
	randId := rand.Int()
	threadId := strconv.Itoa(randId)
	savedThread := model.Thread{
		ThreadId: &threadId,
		Title:    thread.Title,
	}
	m.ThreadMap[*savedThread.ThreadId] = &savedThread
	return &savedThread, nil
}
func (m *MockThreadRepository) CreateThreadPost(post model.Post) (*model.Post, error) {
	thread, err := m.GetThread(*post.ThreadId, "1")
	if err != nil {
		return nil, err
	}
	thread.Posts = append(thread.Posts, &post)
	return &post, nil
}
func (m *MockThreadRepository) UpdateThreadPost(post model.Post) error {
	return nil
}
func (m *MockThreadRepository) DeleteThreadPost(post model.Post) error {
	return nil
}

type MockUserRepository struct {
	UserIdMap map[string]string
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	userId, present := m.UserIdMap[email]
	if !present {
		return nil, errors.New("user not found")
	}
	return &model.User{UserId: &userId}, nil
}
