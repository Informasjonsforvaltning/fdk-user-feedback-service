package integration_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/controller"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/tests"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/util"
)

func TestUserFeedbackService(t *testing.T) {
	endpointUrl := "https://example.com"
	routePath := "/route"

	t.Run("Get posts", func(t *testing.T) {
		entityIds, _, threadMap, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		currentEntity := entityIds[0]
		expectedStatusCode := http.StatusOK
		expectedResponse := threadMap[currentEntity]

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity),
			nil,
		)
		controller.CurrentController.GetComments(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}

		var actualResponse model.ThreadDTO
		err := json.Unmarshal(w.CurrentWriteOutput, &actualResponse)

		if err != nil {
			t.Fatal("error decoding response")
		}

		if reflect.DeepEqual(actualResponse.ToThread(), expectedResponse) {
			t.Fatalf("Expected: %#v\nGot: %#v", expectedResponse, actualResponse)
		}
	})

	t.Run("Get posts for entity without thread", func(t *testing.T) {
		entityIds, _, _, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		currentEntity := entityIds[1]
		expectedStatusCode := http.StatusNotFound

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity),
			nil,
		)
		controller.CurrentController.GetComments(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}
	})

	t.Run("Create post", func(t *testing.T) {
		entityIds, emails, _, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		postId := "100"
		userId := "1"
		threadId := "1"
		content := "test post"
		requestBody, _ := util.ProcessRequestBody(&map[string]string{
			"pid":     postId,
			"content": content,
		})

		currentEntity := entityIds[0]
		expectedStatusCode := http.StatusCreated
		expectedResponse := model.Post{
			PostId:   &postId,
			UserId:   &userId,
			ThreadId: &threadId,
			Content:  &content,
		}

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity),
			requestBody,
		)
		audience := []string{"fdk-feedback-service"}
		jwt := tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &emails[0], &audience)
		r.Header.Set("Authorization", *jwt)
		controller.CurrentController.CreateComment(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}

		var actualResponse model.Post
		err := json.Unmarshal(w.CurrentWriteOutput, &actualResponse)

		if err != nil {
			t.Fatal("error decoding response")
		}

		if !tests.DeepEqualsPost(&actualResponse, &expectedResponse) {
			t.Fatalf("Expected: %#v\nGot: %#v", expectedResponse, actualResponse)
		}
	})

	t.Run("Create new thread and post", func(t *testing.T) {
		entityIds, emails, _, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		postId := "100"
		userId := "1"
		content := "test post"
		requestBody, _ := util.ProcessRequestBody(&map[string]string{
			"pid":     postId,
			"content": content,
		})

		currentEntity := entityIds[1]
		expectedStatusCode := http.StatusCreated

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity),
			requestBody,
		)

		audience := []string{"fdk-feedback-service"}
		jwt := tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &emails[0], &audience)
		r.Header.Set("Authorization", *jwt)
		controller.CurrentController.CreateComment(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}

		var actualResponse model.Post
		err := json.Unmarshal(w.CurrentWriteOutput, &actualResponse)

		if err != nil {
			t.Fatal("error decoding response")
		}

		expectedResponse := model.Post{
			PostId:   &postId,
			UserId:   &userId,
			ThreadId: actualResponse.ThreadId,
			Content:  &content,
		}

		if !tests.DeepEqualsPost(&actualResponse, &expectedResponse) {
			t.Fatalf("Expected: %#v\nGot: %#v", expectedResponse, actualResponse)
		}
	})

	t.Run("Edit own post", func(t *testing.T) {
		entityIds, emails, _, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		postId := "2"
		userId := "1"
		threadId := "1"
		content := "test post"
		requestBody, _ := util.ProcessRequestBody(&map[string]string{
			"pid":     postId,
			"content": content,
		})

		currentEntity := entityIds[0]
		expectedStatusCode := http.StatusOK

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity+"/"+postId),
			requestBody,
		)

		audience := []string{"fdk-feedback-service"}
		jwt := tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &emails[0], &audience)
		r.Header.Set("Authorization", *jwt)
		controller.CurrentController.UpdateComment(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}

		var actualResponse model.Post
		err := json.Unmarshal(w.CurrentWriteOutput, &actualResponse)

		if err != nil {
			t.Fatal("error decoding response")
		}

		expectedResponse := model.Post{
			PostId:   &postId,
			UserId:   &userId,
			ThreadId: &threadId,
			Content:  &content,
		}

		if !tests.DeepEqualsPost(&actualResponse, &expectedResponse) {
			t.Fatalf("Expected: %#v\nGot: %#v", expectedResponse, actualResponse)
		}
	})

	t.Run("Edit other users post", func(t *testing.T) {
		entityIds, emails, _, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		postId := "1"
		content := "test post"
		requestBody, _ := util.ProcessRequestBody(&map[string]string{
			"pid":     postId,
			"content": content,
		})

		currentEntity := entityIds[0]
		expectedStatusCode := http.StatusUnauthorized

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity+"/"+postId),
			requestBody,
		)

		audience := []string{"fdk-feedback-service"}
		jwt := tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &emails[0], &audience)
		r.Header.Set("Authorization", *jwt)
		controller.CurrentController.UpdateComment(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}
	})

	t.Run("Delete own post", func(t *testing.T) {
		entityIds, emails, _, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		postId := "2"

		currentEntity := entityIds[0]
		expectedStatusCode := http.StatusOK

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity+"/"+postId),
			nil,
		)

		audience := []string{"fdk-feedback-service"}
		jwt := tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &emails[0], &audience)
		r.Header.Set("Authorization", *jwt)
		controller.CurrentController.DeleteComment(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}
	})

	t.Run("Delete other users post", func(t *testing.T) {
		entityIds, emails, _, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		postId := "1"

		currentEntity := entityIds[0]
		expectedStatusCode := http.StatusUnauthorized

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity+"/"+postId),
			nil,
		)

		audience := []string{"fdk-feedback-service"}
		jwt := tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &emails[0], &audience)
		r.Header.Set("Authorization", *jwt)
		controller.CurrentController.DeleteComment(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}
	})

	t.Run("Create post expired token", func(t *testing.T) {
		entityIds, emails, _, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		postId := "100"
		content := "test post"
		requestBody, _ := util.ProcessRequestBody(&map[string]string{
			"pid":     postId,
			"content": content,
		})

		currentEntity := entityIds[0]
		expectedStatusCode := http.StatusUnauthorized

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity),
			requestBody,
		)
		audience := []string{"fdk-feedback-service"}
		jwt := tests.CreateMockJwt(time.Now().Add(-time.Hour).Unix(), &emails[0], &audience)
		r.Header.Set("Authorization", *jwt)
		controller.CurrentController.CreateComment(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}
	})

	t.Run("Create post wrong audience", func(t *testing.T) {
		entityIds, emails, _, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		postId := "100"
		content := "test post"
		requestBody, _ := util.ProcessRequestBody(&map[string]string{
			"pid":     postId,
			"content": content,
		})

		currentEntity := entityIds[0]
		expectedStatusCode := http.StatusForbidden

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity),
			requestBody,
		)
		audience := []string{"fdk-comments-service"}
		jwt := tests.CreateMockJwt(time.Now().Add(time.Hour).Unix(), &emails[0], &audience)
		r.Header.Set("Authorization", *jwt)
		controller.CurrentController.CreateComment(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}
	})

	t.Run("Create post no auth", func(t *testing.T) {
		entityIds, _, _, mockJwkStore := ConfigureIntegrationTests()
		defer mockJwkStore.Close()

		postId := "100"
		content := "test post"
		requestBody, _ := util.ProcessRequestBody(&map[string]string{
			"pid":     postId,
			"content": content,
		})

		currentEntity := entityIds[0]
		expectedStatusCode := http.StatusUnauthorized

		w := tests.MockResponseWriter{}
		r, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprint(endpointUrl+routePath+"/"+currentEntity),
			requestBody,
		)

		r.Header.Set("Authorization", "")
		controller.CurrentController.CreateComment(&w, r)

		if w.CurrentStatusCode != expectedStatusCode {
			t.Fatalf("expected statuscode %d, got %d", expectedStatusCode, w.CurrentStatusCode)
		}
	})
}
