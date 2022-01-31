package unit_tests

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/controller"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/tests"
)

func setUpControllerMocks() (*tests.MockResponseWriter, *MockAuthService, *MockThreadIdService, *MockThreadService, controller.Controller) {
	mockResponseWriter := tests.MockResponseWriter{}
	mockAuthService := MockAuthService{}
	mockThreadIdService := MockThreadIdService{}
	mockThreadService := MockThreadService{}
	controller := controller.ControllerImpl{
		AuthService:     &mockAuthService,
		ThreadIdService: &mockThreadIdService,
		ThreadService:   &mockThreadService,
	}

	return &mockResponseWriter, &mockAuthService, &mockThreadIdService, &mockThreadService, &controller
}

func TestCreateComment(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Test unauthorized call", func(t *testing.T) {
		mockResponseWriter, mockAuthService, _, _, controller := setUpControllerMocks()
		expectedStatusCode := http.StatusUnauthorized
		mockAuthService.MockStatusCode = expectedStatusCode

		controller.CreateComment(mockResponseWriter, &http.Request{})

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Test empty request path", func(t *testing.T) {
		mockResponseWriter, mockAuthService, _, _, controller := setUpControllerMocks()
		mockAuthService.MockStatusCode = http.StatusOK
		userId := "1"
		mockAuthService.MockUser = &model.User{UserId: &userId}

		request, _ := http.NewRequest(
			http.MethodPost,
			"/",
			nil,
		)

		controller.CreateComment(mockResponseWriter, request)

		expectedStatusCode := http.StatusNotFound

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Handles threadservice get thread id error", func(t *testing.T) {
		mockResponseWriter, mockAuthService, _, mockThreadService, controller := setUpControllerMocks()
		mockAuthService.MockStatusCode = http.StatusOK
		userId := "1"
		mockAuthService.MockUser = &model.User{UserId: &userId}
		mockThreadService.MockStatusCode = http.StatusInternalServerError

		requestBody := bytes.NewBuffer([]byte(`{"test": "true"}`))
		request, _ := http.NewRequest(
			http.MethodPost,
			"/route/entityId",
			requestBody,
		)

		controller.CreateComment(mockResponseWriter, request)

		expectedStatusCode := http.StatusInternalServerError

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Handles threadservice create thread failure", func(t *testing.T) {
		mockResponseWriter, mockAuthService, _, mockThreadService, controller := setUpControllerMocks()

		expectedStatusCode := http.StatusInternalServerError

		mockAuthService.MockStatusCode = http.StatusOK
		userId := "1"
		mockAuthService.MockUser = &model.User{UserId: &userId}
		mockThreadService.MockStatusCode = expectedStatusCode

		requestBody := bytes.NewBuffer([]byte(`{"test": "true"}`))
		request, _ := http.NewRequest(
			http.MethodPost,
			"/route/entityId",
			requestBody,
		)

		controller.CreateComment(mockResponseWriter, request)

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Successfully creates thread and post", func(t *testing.T) {
		mockResponseWriter, mockAuthService, mockThreadIdService, mockThreadService, controller := setUpControllerMocks()

		userId := "1"
		threadId := "1"
		content := "test post"

		expectedPost := model.Post{UserId: &userId, ThreadId: &threadId, Content: &content}
		expectedStatusCode := http.StatusCreated

		mockAuthService.MockStatusCode = http.StatusOK
		mockAuthService.MockUser = &model.User{UserId: &userId}
		mockThreadIdService.MockThreadId = &threadId
		mockThreadService.MockStatusCode = http.StatusOK
		mockThreadService.MockPost = &expectedPost

		requestBody := bytes.NewBuffer([]byte(`{"test": "true"}`))
		request, _ := http.NewRequest(
			http.MethodPost,
			"/route/entityId",
			requestBody,
		)

		controller.CreateComment(mockResponseWriter, request)

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})
}

func TestGetComments(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Could not parse path", func(t *testing.T) {
		mockResponseWriter, _, _, _, controller := setUpControllerMocks()
		expectedStatusCode := http.StatusNotFound

		request, _ := http.NewRequest(
			http.MethodGet,
			"/",
			nil,
		)

		controller.GetComments(mockResponseWriter, request)

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Handles get thread id error", func(t *testing.T) {
		mockResponseWriter, _, _, mockThreadService, controller := setUpControllerMocks()
		expectedStatusCode := http.StatusNotFound
		mockThreadService.MockStatusCode = http.StatusNotFound

		request, _ := http.NewRequest(
			http.MethodGet,
			"/route/entityId",
			nil,
		)

		controller.GetComments(mockResponseWriter, request)

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Handles get thread error", func(t *testing.T) {
		mockResponseWriter, _, mockThreadIdService, mockThreadService, controller := setUpControllerMocks()
		expectedStatusCode := http.StatusBadRequest
		threadId := "1"
		mockThreadIdService.MockThreadId = &threadId
		mockThreadService.MockStatusCode = http.StatusBadRequest

		request, _ := http.NewRequest(
			http.MethodGet,
			"/route/entityId",
			nil,
		)

		controller.GetComments(mockResponseWriter, request)

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Successfully gets comments", func(t *testing.T) {
		mockResponseWriter, _, mockThreadIdService, mockThreadService, controller := setUpControllerMocks()
		expectedStatusCode := http.StatusOK
		threadId := "1"
		mockThreadIdService.MockThreadId = &threadId
		mockThreadService.MockThread = &model.Thread{}
		mockThreadService.MockStatusCode = http.StatusOK

		request, _ := http.NewRequest(
			http.MethodGet,
			"/route/entityId",
			nil,
		)

		controller.GetComments(mockResponseWriter, request)

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})
}

func TestUpdateComment(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Test unauthorized call", func(t *testing.T) {
		mockResponseWriter, mockAuthService, _, _, controller := setUpControllerMocks()
		expectedStatusCode := http.StatusUnauthorized
		mockAuthService.MockStatusCode = expectedStatusCode

		controller.UpdateComment(mockResponseWriter, &http.Request{})

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Test empty request path", func(t *testing.T) {
		mockResponseWriter, mockAuthService, _, _, controller := setUpControllerMocks()
		mockAuthService.MockStatusCode = http.StatusOK
		userId := "1"
		mockAuthService.MockUser = &model.User{UserId: &userId}

		request, _ := http.NewRequest(
			http.MethodPut,
			"/",
			nil,
		)

		controller.UpdateComment(mockResponseWriter, request)

		expectedStatusCode := http.StatusNotFound

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Handles threadservice update thread failure", func(t *testing.T) {
		mockResponseWriter, mockAuthService, mockThreadIdService, mockThreadService, controller := setUpControllerMocks()

		expectedStatusCode := http.StatusInternalServerError

		mockAuthService.MockStatusCode = http.StatusOK
		userId := "1"
		threadId := "1"
		mockAuthService.MockUser = &model.User{UserId: &userId}
		mockThreadIdService.MockThreadId = &threadId
		mockThreadService.MockStatusCode = expectedStatusCode

		requestBody := bytes.NewBuffer([]byte(`{"test": "true"}`))
		request, _ := http.NewRequest(
			http.MethodPut,
			"/route/entityId/postId",
			requestBody,
		)

		controller.UpdateComment(mockResponseWriter, request)

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Successfully updates post", func(t *testing.T) {
		mockResponseWriter, mockAuthService, mockThreadIdService, mockThreadService, controller := setUpControllerMocks()

		userId := "1"
		threadId := "1"
		content := "test post"

		expectedPost := model.Post{UserId: &userId, ThreadId: &threadId, Content: &content}
		expectedStatusCode := http.StatusOK

		mockAuthService.MockStatusCode = http.StatusOK
		mockAuthService.MockUser = &model.User{UserId: &userId}
		mockThreadIdService.MockThreadId = &threadId
		mockThreadService.MockStatusCode = http.StatusOK
		mockThreadService.MockPost = &expectedPost

		requestBody := bytes.NewBuffer([]byte(`{"test": "true"}`))
		request, _ := http.NewRequest(
			http.MethodPut,
			"/route/entityId/threadId",
			requestBody,
		)

		controller.UpdateComment(mockResponseWriter, request)

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})
}

func TestDeleteComment(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Test unauthorized call", func(t *testing.T) {
		mockResponseWriter, mockAuthService, _, _, controller := setUpControllerMocks()
		expectedStatusCode := http.StatusUnauthorized
		mockAuthService.MockStatusCode = expectedStatusCode

		controller.DeleteComment(mockResponseWriter, &http.Request{})

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Test empty request path", func(t *testing.T) {
		mockResponseWriter, mockAuthService, mockThreadIdService, _, controller := setUpControllerMocks()
		mockAuthService.MockStatusCode = http.StatusOK
		userId := "1"
		threadId := "1"
		mockAuthService.MockUser = &model.User{UserId: &userId}
		mockThreadIdService.MockThreadId = &threadId

		request, _ := http.NewRequest(
			http.MethodPut,
			"/",
			nil,
		)

		controller.DeleteComment(mockResponseWriter, request)

		expectedStatusCode := http.StatusNotFound

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Handles threadservice delete thread failure", func(t *testing.T) {
		mockResponseWriter, mockAuthService, mockThreadIdService, mockThreadService, controller := setUpControllerMocks()

		expectedStatusCode := http.StatusInternalServerError

		mockAuthService.MockStatusCode = http.StatusOK
		userId := "1"
		threadId := "1"
		mockAuthService.MockUser = &model.User{UserId: &userId}
		mockThreadService.MockStatusCode = expectedStatusCode
		mockThreadIdService.MockThreadId = &threadId

		requestBody := bytes.NewBuffer([]byte(`{"test": "true"}`))
		request, _ := http.NewRequest(
			http.MethodPut,
			"/route/entityId/threadId",
			requestBody,
		)

		controller.DeleteComment(mockResponseWriter, request)

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Successfully deletes post", func(t *testing.T) {
		mockResponseWriter, mockAuthService, mockThreadIdService, mockThreadService, controller := setUpControllerMocks()

		userId := "1"
		threadId := "1"
		content := "test post"

		expectedPost := model.Post{UserId: &userId, ThreadId: &threadId, Content: &content}
		expectedStatusCode := http.StatusOK

		mockAuthService.MockStatusCode = http.StatusOK
		mockAuthService.MockUser = &model.User{UserId: &userId}
		mockThreadIdService.MockThreadId = &threadId
		mockThreadService.MockStatusCode = http.StatusOK
		mockThreadService.MockPost = &expectedPost

		requestBody := bytes.NewBuffer([]byte(`{"test": "true"}`))
		request, _ := http.NewRequest(
			http.MethodPut,
			"/route/entityId/threadId",
			requestBody,
		)

		controller.DeleteComment(mockResponseWriter, request)

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})
}

func TestCurrentUser(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Test unauthorized call", func(t *testing.T) {
		mockResponseWriter, mockAuthService, _, _, controller := setUpControllerMocks()
		expectedStatusCode := http.StatusUnauthorized
		mockAuthService.MockStatusCode = expectedStatusCode

		controller.CurrentUser(mockResponseWriter, &http.Request{})

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Test user not found", func(t *testing.T) {
		mockResponseWriter, mockAuthService, _, _, controller := setUpControllerMocks()
		expectedStatusCode := http.StatusUnauthorized
		mockAuthService.MockStatusCode = http.StatusOK

		controller.CurrentUser(mockResponseWriter, &http.Request{})

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})

	t.Run("Successfully gets current user", func(t *testing.T) {
		mockResponseWriter, mockAuthService, _, _, controller := setUpControllerMocks()
		testUser := "testUser"
		expectedStatusCode := http.StatusOK
		expectedUser := model.User{UserId: &testUser}
		mockAuthService.MockStatusCode = expectedStatusCode
		mockAuthService.MockUser = &expectedUser

		controller.CurrentUser(mockResponseWriter, &http.Request{})

		if mockResponseWriter.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected %d. Got %d", expectedStatusCode, mockResponseWriter.CurrentStatusCode)
		}
	})
}
