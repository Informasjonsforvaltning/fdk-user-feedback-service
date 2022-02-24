package unit_tests

import (
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/golang-jwt/jwt/v4"
)

type MockEntityRepository struct {
	MockEntity *model.Entity
	MockError  error
}

func (m *MockEntityRepository) GetEntityById(entityID string) (*model.Entity, error) {
	return m.MockEntity, m.MockError
}

type MockThreadIdRepository struct {
	MockThreadId *string
	MockError    error
}

func (m *MockThreadIdRepository) GetThreadId(id string) (*string, error) {
	return m.MockThreadId, m.MockError
}

func (m *MockThreadIdRepository) CreateThreadId(id string, threadId string) error {
	return m.MockError
}

type MockThreadRepository struct {
	MockError     error
	MockThread    *model.Thread
	MockPost      *model.Post
	MockGetThread *model.Thread
	MockGetError  error
}

func (m *MockThreadRepository) GetThread(threadId string, page *string) (*model.Thread, error) {
	return m.MockGetThread, m.MockGetError
}
func (m *MockThreadRepository) CreateThread(thread model.Thread) (*model.Thread, error) {
	return m.MockThread, m.MockError
}
func (m *MockThreadRepository) CreateThreadPost(post model.Post) (*model.Post, error) {
	return m.MockPost, m.MockError
}
func (m *MockThreadRepository) UpdateThreadPost(post model.Post) error {
	return m.MockError
}
func (m *MockThreadRepository) DeleteThreadPost(post model.Post) error {
	return m.MockError
}

type MockUserRepository struct {
	MockUser  *model.User
	MockError error
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	return m.MockUser, m.MockError
}

type MockThreadIdService struct {
	MockThreadId *string
	MockError    error
}

func (m *MockThreadIdService) GetThreadId(id string) (*string, error) {
	return m.MockThreadId, m.MockError
}
func (m *MockThreadIdService) CreateThreadId(id string, threadId string) error {
	return m.MockError
}

type MockEntityService struct {
	MockEntity *model.Entity
	MockError  error
}

func (m *MockEntityService) GetEntity(id string) (*model.Entity, error) {
	return m.MockEntity, m.MockError
}

type MockAuthService struct {
	MockUser       *model.User
	MockStatusCode int
	MockClaims     *jwt.MapClaims
	MockError      error
}

func (m *MockAuthService) AuthenticateAndGetUser(jwt string) (*model.User, int) {
	return m.MockUser, m.MockStatusCode
}
func (m *MockAuthService) AuthenticateJwt(jwt string) (*jwt.MapClaims, int) {
	return m.MockClaims, m.MockStatusCode
}
func (m *MockAuthService) GetUser(email string) (*model.User, error) {
	return m.MockUser, m.MockError
}

type MockThreadService struct {
	MockPost       *model.Post
	MockThread     *model.Thread
	MockStatusCode int
}

func (m *MockThreadService) CreateThreadPost(postRequest model.Post) (*model.Post, int) {
	return m.MockPost, m.MockStatusCode
}
func (m *MockThreadService) CreateThread(forEntityId string) (*model.Thread, int) {
	return m.MockThread, m.MockStatusCode
}
func (m *MockThreadService) GetThread(id string, page *string) (*model.Thread, int) {
	return m.MockThread, m.MockStatusCode
}
func (m *MockThreadService) UpdateThreadPost(updatedPost model.Post) (*model.Post, int) {
	return m.MockPost, m.MockStatusCode
}
func (m *MockThreadService) DeleteThreadPost(postToDelete model.Post) int {
	return m.MockStatusCode
}

func (m *MockThreadService) CreatePostForEntityId(postRequest model.Post, entityId string) (*model.Post, int) {
	return m.MockPost, m.MockStatusCode
}

func (m *MockThreadService) GetThreadByEntityId(entityId string, page *string) (*model.Thread, int) {
	return m.MockThread, m.MockStatusCode
}
