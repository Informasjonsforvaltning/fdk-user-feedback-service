package unit_tests

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/service"
)

func threadServiceMocks() (*MockEntityService, *MockThreadIdService, *MockThreadRepository, service.ThreadService) {

	mockThreadRepository := MockThreadRepository{}

	mockThreadIdService := MockThreadIdService{}

	mockEntityService := MockEntityService{}

	threadService := service.ThreadServiceImpl{
		ThreadRepository: &mockThreadRepository,
		ThreadIdService:  &mockThreadIdService,
		EntityService:    &mockEntityService,
	}

	return &mockEntityService, &mockThreadIdService, &mockThreadRepository, &threadService
}

func TestCreateThreadPost(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Empty post request", func(t *testing.T) {
		_, _, _, threadService := threadServiceMocks()

		var expectedPost *model.Post
		expectedStatusCode := http.StatusBadRequest

		actualPost, actualStatusCode := threadService.CreateThreadPost(model.Post{})

		if actualPost != expectedPost || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedPost, expectedStatusCode, actualPost, actualStatusCode)
		}
	})

	t.Run("Could not create thread post", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadId, postId, userId, content := "1", "1", "1", "content"

		var expectedPost *model.Post
		expectedStatusCode := http.StatusInternalServerError
		mockThreadRepository.MockError = errors.New("testerror")

		actualPost, actualStatusCode := threadService.CreateThreadPost(model.Post{
			ThreadId: &threadId,
			PostId:   &postId,
			UserId:   &userId,
			Content:  &content,
		})

		if actualPost != expectedPost || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedPost, expectedStatusCode, actualPost, actualStatusCode)
		}
	})

	t.Run("Valid create thread post", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadId, postId, userId, content := "1", "1", "1", "content"

		var expectedPost = model.Post{
			ThreadId: &threadId,
			PostId:   &postId,
			UserId:   &userId,
			Content:  &content,
		}
		expectedStatusCode := http.StatusCreated
		mockThreadRepository.MockPost = &expectedPost
		mockThreadRepository.MockError = nil
		actualPost, actualStatusCode := threadService.CreateThreadPost(expectedPost)
		if actualPost != &expectedPost || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedPost, expectedStatusCode, actualPost, actualStatusCode)
		}
	})
}

func TestCreateThread(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Entity not found", func(t *testing.T) {
		mockEntityService, _, _, threadService := threadServiceMocks()

		var expectedThread *model.Thread
		mockEntityService.MockError = errors.New("testerror")
		expectedStatusCode := http.StatusNotFound

		actualThread, actualStatusCode := threadService.CreateThread("")

		if actualThread != expectedThread || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedThread, expectedStatusCode, actualThread, actualStatusCode)
		}
	})

	t.Run("No entity title", func(t *testing.T) {
		mockEntityService, _, _, threadService := threadServiceMocks()

		var expectedThread *model.Thread
		mockEntityService.MockEntity = &model.Entity{}
		expectedStatusCode := http.StatusBadRequest

		actualThread, actualStatusCode := threadService.CreateThread("testid")

		if actualThread != expectedThread || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedThread, expectedStatusCode, actualThread, actualStatusCode)
		}
	})

	t.Run("Could not create thread", func(t *testing.T) {
		mockEntityService, _, mockThreadRepository, threadService := threadServiceMocks()

		var expectedThread *model.Thread
		mockEntityService.MockEntity = &model.Entity{
			Title: "test title",
		}
		mockThreadRepository.MockError = errors.New("testerror")
		expectedStatusCode := http.StatusInternalServerError

		actualThread, actualStatusCode := threadService.CreateThread("testid")

		if actualThread != expectedThread || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedThread, expectedStatusCode, actualThread, actualStatusCode)
		}
	})

	t.Run("Could not create thread ID", func(t *testing.T) {
		mockEntityService, mockThreadIdService, mockThreadRepository, threadService := threadServiceMocks()

		var expectedThread *model.Thread
		mockEntityService.MockEntity = &model.Entity{
			Title: "test title",
		}
		mockThreadRepository.MockThread = &model.Thread{}
		mockThreadIdService.MockError = errors.New("testerror")
		expectedStatusCode := http.StatusInternalServerError

		actualThread, actualStatusCode := threadService.CreateThread("testid")

		if actualThread != expectedThread || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedThread, expectedStatusCode, actualThread, actualStatusCode)
		}
	})

	t.Run("Successfully creates thread", func(t *testing.T) {
		mockEntityService, _, mockThreadRepository, threadService := threadServiceMocks()
		threadTitle := "test title"
		threadId := "1"

		expectedThread := model.Thread{
			Title:    &threadTitle,
			ThreadId: &threadId,
		}
		mockEntityService.MockEntity = &model.Entity{
			Title: threadTitle,
		}
		mockThreadRepository.MockThread = &expectedThread
		expectedStatusCode := http.StatusCreated

		actualThread, actualStatusCode := threadService.CreateThread("testid")
		if actualThread != &expectedThread || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedThread, expectedStatusCode, actualThread, actualStatusCode)
		}
	})
}

func TestGetThread(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Handles thread repository error", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()

		var expectedThread *model.Thread
		mockThreadRepository.MockGetError = errors.New("testerror")
		expectedStatusCode := http.StatusNotFound

		actualThread, actualStatusCode := threadService.GetThread("")

		if actualThread != expectedThread || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedThread, expectedStatusCode, actualThread, actualStatusCode)
		}
	})

	t.Run("No thread found", func(t *testing.T) {
		_, _, _, threadService := threadServiceMocks()

		var expectedThread *model.Thread
		expectedStatusCode := http.StatusNotFound

		actualThread, actualStatusCode := threadService.GetThread("")

		if actualThread != expectedThread || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedThread, expectedStatusCode, actualThread, actualStatusCode)
		}
	})

	t.Run("Successfully gets thread", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadTitle := "test thread"

		expectedThread := model.Thread{
			Title: &threadTitle,
		}
		expectedStatusCode := http.StatusOK
		mockThreadRepository.MockGetThread = &expectedThread

		actualThread, actualStatusCode := threadService.GetThread("")

		if !reflect.DeepEqual(*actualThread, expectedThread) || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedThread, expectedStatusCode, actualThread, actualStatusCode)
		}
	})
}

func TestUpdateThreadPost(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Could not get thread", func(t *testing.T) {
		_, _, _, threadService := threadServiceMocks()
		threadId := "1"

		var expectedPost *model.Post
		expectedStatusCode := http.StatusNotFound

		actualPost, actualStatusCode := threadService.UpdateThreadPost(model.Post{
			ThreadId: &threadId,
		})

		if actualPost != expectedPost || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedPost, expectedStatusCode, actualPost, actualStatusCode)
		}
	})

	t.Run("Thread does not contain post", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadId := "1"

		var expectedPost *model.Post
		expectedStatusCode := http.StatusNotFound

		mockThreadRepository.MockGetThread = &model.Thread{}

		actualPost, actualStatusCode := threadService.UpdateThreadPost(model.Post{
			ThreadId: &threadId,
		})

		if actualPost != expectedPost || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedPost, expectedStatusCode, actualPost, actualStatusCode)
		}
	})

	t.Run("Post to update attributed to different user", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadId, postId, userId, otherUserId := "1", "1", "1", "2"
		post := model.Post{PostId: &postId, UserId: &userId}

		var expectedPost *model.Post
		expectedStatusCode := http.StatusUnauthorized

		posts := []*model.Post{&post}
		mockThreadRepository.MockGetThread = &model.Thread{
			Posts: posts,
		}

		actualPost, actualStatusCode := threadService.UpdateThreadPost(model.Post{
			ThreadId: &threadId,
			PostId:   &postId,
			UserId:   &otherUserId,
		})

		if actualPost != expectedPost || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedPost, expectedStatusCode, actualPost, actualStatusCode)
		}
	})

	t.Run("Could not update post", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadId, postId, userId := "1", "1", "1"
		post := model.Post{PostId: &postId, UserId: &userId}

		var expectedPost *model.Post
		expectedStatusCode := http.StatusInternalServerError

		posts := []*model.Post{&post}
		mockThreadRepository.MockGetThread = &model.Thread{
			Posts: posts,
		}
		mockThreadRepository.MockError = errors.New("test error")

		actualPost, actualStatusCode := threadService.UpdateThreadPost(model.Post{
			ThreadId: &threadId,
			PostId:   &postId,
			UserId:   &userId,
		})

		if actualPost != expectedPost || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedPost, expectedStatusCode, actualPost, actualStatusCode)
		}
	})

	t.Run("Successful post update", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadId, postId, userId, content := "1", "1", "1", "content"
		post := model.Post{PostId: &postId, UserId: &userId}

		expectedPost := model.Post{
			ThreadId: &threadId,
			PostId:   &postId,
			UserId:   &userId,
			Content:  &content,
		}
		expectedStatusCode := http.StatusOK

		posts := []*model.Post{&post}
		mockThreadRepository.MockGetThread = &model.Thread{
			Posts: posts,
		}

		actualPost, actualStatusCode := threadService.UpdateThreadPost(expectedPost)

		if *actualPost != expectedPost || actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %#v, %d. Got: %#v, %d", expectedPost, expectedStatusCode, actualPost, actualStatusCode)
		}
	})
}

func TestDeleteThreadPost(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Could not get thread", func(t *testing.T) {
		_, _, _, threadService := threadServiceMocks()
		threadID := "1"

		expectedStatusCode := http.StatusNotFound

		actualStatusCode := threadService.DeleteThreadPost(model.Post{
			ThreadId: &threadID,
		})

		if actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %d. Got: %d", expectedStatusCode, actualStatusCode)
		}
	})

	t.Run("Thread does not contain post", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadID := "1"

		expectedStatusCode := http.StatusNotFound

		mockThreadRepository.MockGetThread = &model.Thread{}

		actualStatusCode := threadService.DeleteThreadPost(model.Post{
			ThreadId: &threadID,
		})

		if actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %d. Got: %d", expectedStatusCode, actualStatusCode)
		}
	})

	t.Run("Post to update attributed to different user", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadId, postId, userId, otherUserId := "1", "1", "1", "2"
		post := model.Post{PostId: &postId, UserId: &userId}

		expectedStatusCode := http.StatusUnauthorized

		posts := []*model.Post{&post}
		mockThreadRepository.MockGetThread = &model.Thread{
			Posts: posts,
		}

		actualStatusCode := threadService.DeleteThreadPost(model.Post{
			ThreadId: &threadId,
			PostId:   &postId,
			UserId:   &otherUserId,
		})

		if actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %d. Got: %d", expectedStatusCode, actualStatusCode)
		}
	})

	t.Run("Could not delete post", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadId, postId, userId := "1", "1", "1"
		post := model.Post{PostId: &postId, UserId: &userId}

		expectedStatusCode := http.StatusInternalServerError

		posts := []*model.Post{&post}
		mockThreadRepository.MockGetThread = &model.Thread{
			Posts: posts,
		}
		mockThreadRepository.MockError = errors.New("test error")

		actualStatusCode := threadService.DeleteThreadPost(model.Post{
			ThreadId: &threadId,
			PostId:   &postId,
			UserId:   &userId,
		})

		if actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %d. Got: %d", expectedStatusCode, actualStatusCode)
		}
	})

	t.Run("Successful post delete", func(t *testing.T) {
		_, _, mockThreadRepository, threadService := threadServiceMocks()
		threadId, postId, userId, content := "1", "1", "1", "test"

		expectedPost := model.Post{
			ThreadId: &threadId,
			PostId:   &postId,
			UserId:   &userId,
			Content:  &content,
		}
		expectedStatusCode := http.StatusOK

		posts := []*model.Post{&expectedPost}
		mockThreadRepository.MockGetThread = &model.Thread{
			Posts: posts,
		}

		actualStatusCode := threadService.DeleteThreadPost(expectedPost)

		if actualStatusCode != expectedStatusCode {
			t.Fatalf("expected post response and status code: %d. Got: %d", expectedStatusCode, actualStatusCode)
		}
	})
}
