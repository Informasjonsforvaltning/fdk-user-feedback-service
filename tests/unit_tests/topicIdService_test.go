package unit_tests

import (
	"errors"
	"io/ioutil"
	"log"
	"testing"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/service"
)

func threadIdServiceMocks() (*MockThreadIdRepository, service.ThreadIdService) {

	mockThreadIdRepository := MockThreadIdRepository{}

	threadIdService := service.ThreadIdServiceImpl{
		ThreadIdRepository: &mockThreadIdRepository,
	}

	return &mockThreadIdRepository, &threadIdService
}

func TestGetThreadId(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Handles repository error", func(t *testing.T) {
		mockThreadIdRepository, threadIdService := threadIdServiceMocks()

		var expectedThreadId *string
		expectedError := errors.New("test error")

		mockThreadIdRepository.MockError = expectedError

		actualThreadId, actualError := threadIdService.GetThreadId("testId")

		if actualThreadId != expectedThreadId || actualError != expectedError {
			t.Fatalf("expected: %s, %s. Got: %s, %s", *expectedThreadId, expectedError, *actualThreadId, actualError)
		}
	})

	t.Run("Successful threadId get", func(t *testing.T) {
		mockThreadIdRepository, threadIdService := threadIdServiceMocks()

		expectedThreadId := "test"
		var expectedError error

		mockThreadIdRepository.MockThreadId = &expectedThreadId

		actualThreadId, actualError := threadIdService.GetThreadId("testId")

		if actualThreadId != &expectedThreadId || actualError != expectedError {
			t.Fatalf("expected: %s, %s. Got: %s, %s", expectedThreadId, expectedError, *actualThreadId, actualError)
		}
	})
}

func TestCreateThreadId(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Handles repository error", func(t *testing.T) {
		mockThreadIdRepository, threadIdService := threadIdServiceMocks()

		expectedError := errors.New("test error")

		mockThreadIdRepository.MockError = expectedError

		actualError := threadIdService.CreateThreadId("testId", "testId")

		if actualError != expectedError {
			t.Fatalf("expected: %s. Got: %s", expectedError, actualError)
		}
	})

	t.Run("Successful create", func(t *testing.T) {
		_, threadIdService := threadIdServiceMocks()

		var expectedError error

		actualError := threadIdService.CreateThreadId("testId", "testId")

		if actualError != expectedError {
			t.Fatalf("expected: %s. Got: %s", expectedError, actualError)
		}
	})
}
