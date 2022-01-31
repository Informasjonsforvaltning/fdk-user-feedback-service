package unit_tests

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/service"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/tests"
)

func setUpAuthServiceMocks() (*MockUserRepository, service.AuthService) {
	mockUserRepository := MockUserRepository{}

	authService := service.AuthServiceImpl{
		UserRepository: &mockUserRepository,
		KeycloakHost:   "",
	}

	return &mockUserRepository, &authService
}

func setUpAuthServiceMocksWithJwt() (*MockUserRepository, service.AuthService, *httptest.Server) {
	mockJwkStore := tests.MockJwkStore()

	mockUserRepository := MockUserRepository{}

	authService := service.AuthServiceImpl{
		UserRepository: &mockUserRepository,
		KeycloakHost:   mockJwkStore.URL,
	}

	return &mockUserRepository, &authService, mockJwkStore
}

func TestGetUserId(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("Handles UserId repository error", func(t *testing.T) {
		mockUserRepository, authService := setUpAuthServiceMocks()

		var expectedUser *model.User
		expectedError := errors.New("test error")

		mockUserRepository.MockError = expectedError

		actualUser, actualError := authService.GetUser("")

		if actualUser != expectedUser || actualError != expectedError {
			t.Fatalf("Expected %v, %s. Got %v, %s", *expectedUser, expectedError, *actualUser, actualError)
		}
	})

	t.Run("Successfully get user", func(t *testing.T) {
		mockUserRepository, authService := setUpAuthServiceMocks()

		expectedUserId := "1"
		var expectedError error

		mockUserRepository.MockUser = &model.User{UserId: &expectedUserId}

		actualUser, actualError := authService.GetUser("")

		if actualUser.UserId != &expectedUserId || actualError != expectedError {
			t.Fatalf("Expected %s, %s. Got %s, %s", expectedUserId, expectedError, *actualUser.UserId, actualError)
		}
	})
}

func TestAuthenticateJwt(t *testing.T) {
	{
		_, authService, jwtStore := setUpAuthServiceMocksWithJwt()
		defer jwtStore.Close()

		testMail := "test@test.com"
		testInvalidAud := []string{"fdk-test-service"}
		testValidAud := []string{"fdk-feedback-service"}

		var tests = []struct {
			testName       string
			expectedStatus int
			jwt            string
		}{
			{"Empty JWT", http.StatusUnauthorized, ""},
			{"Expired, no mail no audience", http.StatusUnauthorized, *tests.CreateMockJwt(time.Now().Add(-time.Hour).Unix(), nil, nil)},
			{"Valid, no mail no audience", http.StatusForbidden, *tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), nil, nil)},
			{"Valid, valid mail no audience", http.StatusForbidden, *tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &testMail, nil)},
			{"Valid, valid mail invalid audience", http.StatusForbidden, *tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &testMail, &testInvalidAud)},
			{"Valid, valid mail valid audience", http.StatusOK, *tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &testMail, &testValidAud)},
		}

		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				_, actualStatus := authService.AuthenticateJwt(test.jwt)
				if actualStatus != test.expectedStatus {
					t.Fatalf("Expected %d. Got %d", test.expectedStatus, actualStatus)
				}
			})
		}
	}

	t.Run("No keycloak host", func(t *testing.T) {
		_, authService := setUpAuthServiceMocks()

		expectedStatus := http.StatusUnauthorized
		testMail := "test@test.com"
		testValidAud := []string{"fdk-feedback-service"}

		jwt := tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &testMail, &testValidAud)

		_, actualStatus := authService.AuthenticateJwt(*jwt)
		if actualStatus != expectedStatus {
			t.Fatalf("Expected %d. Got %d", expectedStatus, actualStatus)
		}

	})
}
