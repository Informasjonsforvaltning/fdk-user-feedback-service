package service

import (
	"log"
	"net/http"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/repository"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/util"
)

type ThreadService interface {
	CreateThreadPost(postRequest model.Post) (*model.Post, int)
	CreateThread(forEntityId string) (*model.Thread, int)
	GetThread(id string) (*model.Thread, int)
	UpdateThreadPost(updatedPost model.Post) (*model.Post, int)
	DeleteThreadPost(postToDelete model.Post) int
	CreatePostForEntityId(postRequest model.Post, entityId string) (*model.Post, int)
	GetThreadByEntityId(entityId string) (*model.Thread, int)
}

type ThreadServiceImpl struct {
	ThreadRepository repository.ThreadRepository
	ThreadIdService  ThreadIdService
	EntityService    EntityService
}

func (threadService *ThreadServiceImpl) CreatePostForEntityId(postRequest model.Post, entityId string) (*model.Post, int) {
	threadId, err := threadService.ThreadIdService.GetThreadId(entityId)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	if threadId == nil {
		createdThread, statusCode := threadService.CreateThread(entityId)

		if !util.SuccsessfulStatus(statusCode) {
			return nil, statusCode
		}

		if createdThread == nil {
			return nil, http.StatusInternalServerError
		}

		threadId = createdThread.ThreadId

	}

	created, statusCode := threadService.CreateThreadPost(model.Post{
		PostId:   postRequest.PostId,
		UserId:   postRequest.UserId,
		ThreadId: threadId,
		Content:  postRequest.Content,
		ToPostId: postRequest.ToPostId,
	})

	if !util.SuccsessfulStatus(statusCode) {
		return nil, statusCode
	}

	return created, http.StatusCreated
}

func (threadService *ThreadServiceImpl) GetThreadByEntityId(entityId string) (*model.Thread, int) {
	threadId, err := threadService.ThreadIdService.GetThreadId(entityId)
	if err != nil || threadId == nil {
		return nil, http.StatusNotFound
	}

	thread, statusCode := threadService.GetThread(*threadId)
	if !util.SuccsessfulStatus(statusCode) || thread == nil {
		return nil, statusCode
	}

	return thread, http.StatusOK
}

func (threadService *ThreadServiceImpl) CreateThreadPost(postRequest model.Post) (*model.Post, int) {
	if postRequest.Content == nil || postRequest.UserId == nil || postRequest.ThreadId == nil {
		return nil, http.StatusBadRequest
	}

	post, err := threadService.ThreadRepository.CreateThreadPost(postRequest)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return post, http.StatusCreated
}

func (threadService *ThreadServiceImpl) CreateThread(forEntityId string) (*model.Thread, int) {
	entity, err := threadService.EntityService.GetEntity(forEntityId)
	if err != nil {
		return nil, http.StatusNotFound
	}

	thread, err := entity.ToThread()
	if err != nil || thread == nil {
		return nil, http.StatusBadRequest
	}

	createdThread, err := threadService.ThreadRepository.CreateThread(*thread)
	if err != nil || createdThread == nil || createdThread.ThreadId == nil {
		log.Println("CreateThread error.\n[ERROR] -", err)
		return nil, http.StatusInternalServerError
	}

	err = threadService.ThreadIdService.CreateThreadId(forEntityId, *createdThread.ThreadId)
	if err != nil {
		log.Println("CreateThreadId error.\n[ERROR] -", err)
		return nil, http.StatusInternalServerError
	}

	return createdThread, http.StatusCreated
}

func (threadService *ThreadServiceImpl) GetThread(id string) (*model.Thread, int) {
	thread, err := threadService.ThreadRepository.GetThread(id)
	if err != nil || thread == nil {
		log.Println("GetThread error.\n[ERROR] -", err)
		return nil, http.StatusNotFound
	}

	filteredThread := thread.FilterDeletedPosts()

	return filteredThread, http.StatusOK
}

func (threadService *ThreadServiceImpl) UpdateThreadPost(updatedPost model.Post) (*model.Post, int) {
	thread, statusCode := threadService.GetThread(*updatedPost.ThreadId)
	if !util.SuccsessfulStatus(statusCode) {
		return nil, statusCode
	}

	postToUpdate, err := thread.FindThreadPostById(updatedPost.PostId)
	if err != nil {
		log.Println("FindThreadPostById error.\n[ERROR] -", err)
		return postToUpdate, http.StatusNotFound
	}

	if postToUpdate.UserId != nil && updatedPost.UserId != nil && *postToUpdate.UserId != *updatedPost.UserId {
		log.Println("Unauthorized error.\n[ERROR] -", err)
		return nil, http.StatusUnauthorized
	}

	err = threadService.ThreadRepository.UpdateThreadPost(updatedPost)
	if err != nil {
		log.Println("Could not update post.\n[ERROR] -", err)
		return nil, http.StatusInternalServerError
	}

	return &updatedPost, http.StatusOK
}

func (threadService *ThreadServiceImpl) DeleteThreadPost(postToDelete model.Post) int {
	thread, statusCode := threadService.GetThread(*postToDelete.ThreadId)
	if !util.SuccsessfulStatus(statusCode) {
		return statusCode
	}

	currentPost, err := thread.FindThreadPostById(postToDelete.PostId)
	if err != nil {
		log.Println("FindThreadPostById error.\n[ERROR] -", err)
		return http.StatusNotFound
	}

	if currentPost.UserId != nil && postToDelete.UserId != nil && *currentPost.UserId != *postToDelete.UserId {
		log.Println("Unauthorized error.\n[ERROR] -", err)
		return http.StatusUnauthorized
	}

	err = threadService.ThreadRepository.DeleteThreadPost(postToDelete)
	if err != nil {
		log.Println("Could not update post.\n[ERROR] -", err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

var CurrentThreadService ThreadService
