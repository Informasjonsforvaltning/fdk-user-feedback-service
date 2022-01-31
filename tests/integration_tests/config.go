package integration_tests

import (
	"net/http/httptest"

	controller "github.com/Informasjonsforvaltning/fdk-user-feedback-service/controller"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	repository "github.com/Informasjonsforvaltning/fdk-user-feedback-service/repository"
	service "github.com/Informasjonsforvaltning/fdk-user-feedback-service/service"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/tests"
)

func ConfigureIntegrationTests() ([]string, []string, map[string]*model.Thread, *httptest.Server) {
	mockJwkStore := tests.MockJwkStore()

	entityIds := []string{"entity-with-corresponding-thread", "entity-without-corresponding-thread", "non-existant-entity"}
	threadIds := []string{"1"}
	emails := []string{"a@test.com", "b@test.com", "c@test.com"}
	userIds := []string{"1", "2", "3"}
	postIds := []string{"1", "2"}

	entityMap := map[string]model.Entity{
		entityIds[0]: {
			Title: "Stort testdatasett",
		},
		entityIds[1]: {
			Title: "Åpne Data fra Enhetsregisteret - API Dokumentasjon",
		},
	}

	threadIdMap := map[string]string{
		entityIds[0]: threadIds[0],
	}

	thread_0 := createThread("1", "Stort testdatasett", []*model.Post{
		createPost(postIds[0], threadIds[0], "22", "First post best post", nil),
		createPost(postIds[1], threadIds[0], userIds[0], "Interesting comment.", &postIds[0]),
	})

	threadMap := map[string]*model.Thread{
		threadIds[0]: thread_0,
	}

	userIdMap := map[string]string{
		emails[0]: userIds[0],
		emails[1]: userIds[1],
		emails[2]: userIds[2],
	}

	repository.CurrentEntityRepository = &MockEntityRepository{
		EntityMap: entityMap,
	}
	repository.CurrentThreadIdRepository = &MockThreadIdRepository{
		ThreadIdMap: threadIdMap,
	}
	repository.CurrentThreadRepository = &MockThreadRepository{
		ThreadMap: threadMap,
	}
	repository.CurrentUserRepository = &MockUserRepository{
		UserIdMap: userIdMap,
	}

	service.CurrentAuthService = &service.AuthServiceImpl{
		UserRepository: repository.CurrentUserRepository,
		KeycloakHost:   mockJwkStore.URL,
	}

	service.CurrentEntityService = &service.EntityServiceImpl{
		EntityRepository: repository.CurrentEntityRepository,
	}
	service.CurrentThreadIdService = &service.ThreadIdServiceImpl{
		ThreadIdRepository: repository.CurrentThreadIdRepository,
	}
	service.CurrentThreadService = &service.ThreadServiceImpl{
		ThreadRepository: repository.CurrentThreadRepository,
		ThreadIdService:  service.CurrentThreadIdService,
		EntityService:    service.CurrentEntityService,
	}

	controller.CurrentController = &controller.ControllerImpl{
		AuthService:     service.CurrentAuthService,
		ThreadIdService: service.CurrentThreadIdService,
		ThreadService:   service.CurrentThreadService,
	}

	return entityIds, emails, threadMap, mockJwkStore
}

func createThread(threadId string, title string, posts []*model.Post) *model.Thread {
	return &model.Thread{
		ThreadId: &threadId,
		Title:    &title,
		Posts:    posts,
	}
}

func createPost(postId string, threadId string, userId string, content string, toPostId *string) *model.Post {
	return &model.Post{
		PostId:   &postId,
		ThreadId: &threadId,
		UserId:   &userId,
		Content:  &content,
		ToPostId: toPostId,
	}
}
