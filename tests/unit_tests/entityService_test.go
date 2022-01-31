package unit_tests

import (
	"errors"
	"io/ioutil"
	"log"
	"testing"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/service"
)

func entityServiceMocks() (*MockEntityRepository, service.EntityService) {

	mockEntityRepository := MockEntityRepository{}

	entityService := service.EntityServiceImpl{
		EntityRepository: &mockEntityRepository,
	}

	return &mockEntityRepository, &entityService
}

func TestGetEntity(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Handles repository error", func(t *testing.T) {
		mockEntityRepository, entityService := entityServiceMocks()

		var expectedEntity *model.Entity
		expectedError := errors.New("test error")

		mockEntityRepository.MockError = expectedError

		actualEntity, actualError := entityService.GetEntity("testId")

		if actualEntity != expectedEntity || actualError != expectedError {
			t.Fatalf("expected: %s, %s. Got: %s, %s", *expectedEntity, expectedError, *actualEntity, actualError)
		}
	})

	t.Run("Successful Entity get", func(t *testing.T) {
		mockEntityRepository, entityService := entityServiceMocks()

		expectedEntity := model.Entity{
			Title: "test",
		}
		var expectedError error

		mockEntityRepository.MockEntity = &expectedEntity

		actualEntity, actualError := entityService.GetEntity("testId")

		if actualEntity != &expectedEntity || actualError != expectedError {
			t.Fatalf("expected: %s, %s. Got: %s, %s", expectedEntity, expectedError, *actualEntity, actualError)
		}
	})
}
