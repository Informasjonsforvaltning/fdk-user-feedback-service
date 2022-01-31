package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
)

type RequestOptions struct {
	Method          string
	EndpointUrl     string
	AccessToken     *string
	RequestBody     *map[string]string
	QueryParameters *map[string]string
}

func Request(options RequestOptions) (*[]byte, error) {
	parsedUrl, err := buildUrl(options.EndpointUrl, options.QueryParameters)
	if err != nil {
		return nil, err
	}

	processedBody, err := ProcessRequestBody(options.RequestBody)
	if err != nil {
		return nil, err
	}

	request, err := buildRequest(options.Method, parsedUrl.String(), processedBody)
	if err != nil {
		log.Println("Error on request build.\n[ERROR] -", err)
		return nil, err
	}

	buildHeader(request, options.AccessToken, processedBody)

	response, err := doRequest(request)

	return response, err
}

func SuccsessfulStatus(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

func buildUrl(urlString string, queryparameters *map[string]string) (*url.URL, error) {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		log.Println("Error on Request url parse.\n[ERROR] -", err)
		return nil, err
	}

	queryParams := url.Values{}

	if queryparameters != nil {
		for key, value := range *queryparameters {
			queryParams.Add(key, value)

		}
		parsedUrl.RawQuery = queryParams.Encode()
	}

	return parsedUrl, err
}

func ProcessRequestBody(reqBody *map[string]string) (*bytes.Buffer, error) {
	var bodyBuffer *bytes.Buffer
	var err error
	if reqBody != nil {
		processedBody, err := json.Marshal(*reqBody)
		if err != nil {
			log.Println("Error on body marshal.\n[ERROR] -", err)
			return nil, err
		}

		bodyBuffer = bytes.NewBuffer(processedBody)
	}

	return bodyBuffer, err
}

func buildRequest(method string, endpointUrl string, body *bytes.Buffer) (*http.Request, error) {
	if body != nil {
		return http.NewRequest(method, endpointUrl, body)
	}
	return http.NewRequest(method, endpointUrl, nil)
}

func buildHeader(request *http.Request, accessToken *string, body *bytes.Buffer) {
	if accessToken != nil {
		request.Header.Add("Authorization", "Bearer "+*accessToken)
	}
	if body != nil {
		request.Header.Add("Content-Type", "application/json")
	}
	request.Header.Add("Accept", "*/*")
}

func doRequest(request *http.Request) (*[]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Error on request send.\n[ERROR] -", err)
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return nil, err
	}

	if !SuccsessfulStatus(resp.StatusCode) {
		log.Print(string(resBody))
		return nil, errors.New(resp.Status)
	}

	return &resBody, err
}

func ParseRequestUrlPath(requestPath string) (*string, *string, *string) {
	var route *string
	var entityId *string
	var threadId *string

	if requestPath == "" || requestPath == "/" {
		return nil, nil, nil
	}

	splitPath := strings.Split(requestPath, "/")

	if len(splitPath) >= 2 && splitPath[1] != "" {
		route = &splitPath[1]
	}

	if len(splitPath) >= 3 && splitPath[2] != "" {
		entityId = &splitPath[2]
	}

	if len(splitPath) >= 4 && splitPath[3] != "" {
		threadId = &splitPath[3]
	}

	return route, entityId, threadId
}

func DecodePost(body io.ReadCloser) (*model.Post, error) {
	var post model.Post
	err := json.NewDecoder(body).Decode(&post)
	return &post, err
}

func UnmarshalPost(bytes *[]byte) (*model.Post, error) {
	var dbPost model.PostDTO
	if bytes == nil {
		return nil, model.ErrNoBytes
	}
	err := json.Unmarshal(*bytes, &dbPost)
	if err != nil {
		log.Println("Error on Post unmarshal.\n[ERROR] -", err)
	}

	return dbPost.ToPost(), err
}

func UnmarshalThread(bytes *[]byte) (*model.Thread, error) {
	var dbThread model.ThreadDTO
	if bytes == nil {
		return nil, model.ErrNoBytes
	}
	err := json.Unmarshal(*bytes, &dbThread)
	if err != nil {
		log.Println("Error on Thread unmarshal.\n[ERROR] -", err)
	}

	return dbThread.ToThread(), err
}

func UnmarshalPostResponse(bytes *[]byte) (*model.Post, error) {
	var response model.PostResponseDAO
	if bytes == nil {
		return nil, model.ErrNoBytes
	}

	err := json.Unmarshal(*bytes, &response)
	if err != nil {
		log.Println("Error on Post unmarshal.\n[ERROR] -", err)
	}

	if *response.Code != "ok" {
		return nil, model.ErrBadResponse
	}

	return response.Payload.ToPost(), err
}

func UnmarshalThreadResponse(bytes *[]byte) (*model.Thread, error) {
	var response model.ThreadResponseDTO
	if bytes == nil {
		return nil, model.ErrNoBytes
	}

	err := json.Unmarshal(*bytes, &response)
	if err != nil {
		log.Println("Error on Post unmarshal.\n[ERROR] -", err)
	}

	if *response.Code != "ok" {
		return nil, model.ErrBadResponse
	}

	return response.Payload.ThreadData.ToThread(), err
}

func UnmarshalUser(bytes *[]byte) (*model.User, error) {
	var dbUser model.UserDTO
	if bytes == nil {
		return nil, model.ErrNoBytes
	}
	err := json.Unmarshal(*bytes, &dbUser)
	if err != nil {
		log.Println("Error on User unmarshal.\n[ERROR] -", err)
	}

	return dbUser.ToUser(), err
}
