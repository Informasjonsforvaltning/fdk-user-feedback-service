package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/service"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/util"
)

type Controller interface {
	CreateComment(w http.ResponseWriter, r *http.Request)
	GetComments(w http.ResponseWriter, r *http.Request)
	UpdateComment(w http.ResponseWriter, r *http.Request)
	DeleteComment(w http.ResponseWriter, r *http.Request)
	CurrentUser(w http.ResponseWriter, r *http.Request)
}

type ControllerImpl struct {
	AuthService     service.AuthService
	ThreadIdService service.ThreadIdService
	ThreadService   service.ThreadService
}

func (controller *ControllerImpl) CreateComment(w http.ResponseWriter, r *http.Request) {
	user, statusCode := controller.AuthService.AuthenticateAndGetUser(r.Header.Get("Authorization"))
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, entityId, _ := util.ParseRequestUrlPath(r.URL.Path)
	if entityId == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	post, err := util.DecodePost(r.Body)
	if err != nil || post == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	created, statusCode := controller.ThreadService.CreatePostForEntityId(model.Post{
		PostId:   post.PostId,
		UserId:   user.UserId,
		ThreadId: post.ThreadId,
		Content:  post.Content,
		ToPostId: post.ToPostId,
	}, *entityId)
	if !util.SuccsessfulStatus(statusCode) {
		w.WriteHeader(statusCode)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (controller *ControllerImpl) GetComments(w http.ResponseWriter, r *http.Request) {
	_, entityId, _ := util.ParseRequestUrlPath(r.URL.Path)
	if entityId == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	thread, statusCode := controller.ThreadService.GetThreadByEntityId(*entityId, util.GetPageQueryParam(r.URL.Query()))
	if !util.SuccsessfulStatus(statusCode) || thread == nil {
		w.WriteHeader(statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(thread)
}

func (controller *ControllerImpl) UpdateComment(w http.ResponseWriter, r *http.Request) {
	user, statusCode := controller.AuthService.AuthenticateAndGetUser(r.Header.Get("Authorization"))
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, entityId, postId := util.ParseRequestUrlPath(r.URL.Path)
	if entityId == nil || postId == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	threadId, err := controller.ThreadIdService.GetThreadId(*entityId)
	if err != nil || threadId == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	post, err := util.DecodePost(r.Body)
	if err != nil || post == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	created, statusCode := controller.ThreadService.UpdateThreadPost(model.Post{
		PostId:   postId,
		UserId:   user.UserId,
		ThreadId: threadId,
		Content:  post.Content,
		ToPostId: post.ToPostId,
	})
	if !util.SuccsessfulStatus(statusCode) {
		w.WriteHeader(statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(created)
}

func (controller *ControllerImpl) DeleteComment(w http.ResponseWriter, r *http.Request) {
	user, statusCode := controller.AuthService.AuthenticateAndGetUser(r.Header.Get("Authorization"))
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, entityId, postId := util.ParseRequestUrlPath(r.URL.Path)
	if entityId == nil || postId == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	threadId, err := controller.ThreadIdService.GetThreadId(*entityId)
	if err != nil || threadId == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	statusCode = controller.ThreadService.DeleteThreadPost(model.Post{
		PostId:   postId,
		UserId:   user.UserId,
		ThreadId: threadId,
	})

	w.WriteHeader(statusCode)
}

func (controller *ControllerImpl) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user, statusCode := controller.AuthService.AuthenticateAndGetUser(r.Header.Get("Authorization"))
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

var CurrentController Controller
