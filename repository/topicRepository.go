package repository

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/util"
)

type ThreadRepository interface {
	GetThread(threadId string) (*model.Thread, error)
	CreateThread(thread model.Thread) (*model.Thread, error)
	CreateThreadPost(post model.Post) (*model.Post, error)
	UpdateThreadPost(post model.Post) error
	DeleteThreadPost(post model.Post) error
}

type ThreadRepositoryImpl struct {
	WriteApiToken       string
	ReadApiToken        string
	CommunityApiUrl     string
	TopicPath           string
	TopicsPath          string
	PostsPath           string
	ThreadBotUid        string
	CommunityCategoryId string
}

func (threadRepository *ThreadRepositoryImpl) GetThread(threadId string) (*model.Thread, error) {
	bearerToken := threadRepository.ReadApiToken
	method := http.MethodGet
	endpointUrl := threadRepository.CommunityApiUrl + threadRepository.TopicPath + threadId + "?sort=newest_to_oldest"

	response, err := util.Request(util.RequestOptions{
		Method:      method,
		EndpointUrl: endpointUrl,
		AccessToken: &bearerToken,
	})
	if err != nil {
		log.Println("Error on request.\n[ERROR] -", err, method, endpointUrl)
		return nil, err
	}
	thread, err := util.UnmarshalThread(response)

	return thread, err
}

func (threadRepository *ThreadRepositoryImpl) CreateThread(thread model.Thread) (*model.Thread, error) {
	if thread.Title == nil || thread.Content == nil {
		return nil, fmt.Errorf("cannot create thread without title and content")
	}

	bearerToken := threadRepository.WriteApiToken
	method := http.MethodPost
	endpointUrl := threadRepository.CommunityApiUrl + threadRepository.TopicsPath

	postBody := map[string]string{
		"_uid":    threadRepository.ThreadBotUid,
		"cid":     threadRepository.CommunityCategoryId,
		"title":   *thread.Title,
		"content": *thread.Content,
	}

	response, err := util.Request(util.RequestOptions{
		Method:      method,
		EndpointUrl: endpointUrl,
		AccessToken: &bearerToken,
		RequestBody: &postBody,
	})
	if err != nil {
		log.Println("Error on request.\n[ERROR] -", err, method, endpointUrl)
		return nil, err
	}

	responseThread, err := util.UnmarshalThreadResponse(response)

	return responseThread, err
}

func (threadRepository *ThreadRepositoryImpl) CreateThreadPost(post model.Post) (*model.Post, error) {
	if post.ThreadId == nil || post.Content == nil {
		return nil, fmt.Errorf("cannot create post without threadid and content")
	}

	toPid := ""
	if post.ToPostId != nil {
		toPid = *post.ToPostId
	}

	bearerToken := threadRepository.WriteApiToken
	method := http.MethodPost
	endpointUrl := threadRepository.CommunityApiUrl + threadRepository.TopicsPath + *post.ThreadId

	postBody := map[string]string{
		"_uid":    *post.UserId,
		"content": *post.Content,
		"toPid":   toPid,
	}
	response, err := util.Request(util.RequestOptions{
		Method:      method,
		EndpointUrl: endpointUrl,
		AccessToken: &bearerToken,
		RequestBody: &postBody,
	})
	if err != nil {
		log.Println("Error on request.\n[ERROR] -", err, method, endpointUrl)
		return nil, err
	}

	responsePost, err := util.UnmarshalPostResponse(response)
	if err != nil {
		log.Println("Error on response unmarshal.\n[ERROR] -", err)
		return nil, err
	}

	return responsePost, err
}

func (threadRepository *ThreadRepositoryImpl) UpdateThreadPost(post model.Post) error {
	if post.ThreadId == nil || post.PostId == nil || post.Content == nil {
		return fmt.Errorf("cannot update post without threadid, postid, and content")
	}

	bearerToken := threadRepository.WriteApiToken
	method := http.MethodPut
	endpointUrl := threadRepository.CommunityApiUrl + threadRepository.PostsPath + *post.PostId

	putBody := map[string]string{
		"_uid":    *post.UserId,
		"content": *post.Content,
	}
	_, err := util.Request(util.RequestOptions{
		Method:      method,
		EndpointUrl: endpointUrl,
		AccessToken: &bearerToken,
		RequestBody: &putBody,
	})
	if err != nil {
		log.Println("Error on request.\n[ERROR] -", err, method, endpointUrl)
	}

	return err
}

func (threadRepository *ThreadRepositoryImpl) DeleteThreadPost(post model.Post) error {
	if post.PostId == nil || post.UserId == nil {
		return fmt.Errorf("cannot update post without postId and userId")
	}

	bearerToken := threadRepository.WriteApiToken
	method := http.MethodDelete
	endpointUrl := threadRepository.CommunityApiUrl + threadRepository.PostsPath + *post.PostId + "/state"

	deleteBody := map[string]string{
		"_uid": *post.UserId,
	}

	_, err := util.Request(util.RequestOptions{
		Method:      method,
		EndpointUrl: endpointUrl,
		AccessToken: &bearerToken,
		RequestBody: &deleteBody,
	})
	if err != nil {
		log.Println("Error on request.\n[ERROR] -", err, method, endpointUrl)
	}

	return err
}

var CurrentThreadRepository ThreadRepository
