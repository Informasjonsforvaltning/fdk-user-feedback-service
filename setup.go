package fdk_user_feedback_service

import (
	"fmt"
	"net/http"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/controller"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/env"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/util"
)

func EntryPoint(w http.ResponseWriter, r *http.Request) {
	Configure()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "Application/JSON")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, access-control-allow-origin")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	route, _, _ := util.ParseRequestUrlPath(r.URL.Path)
	if route == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch *route {
	case env.ConstantValues.PingPath:
		ping(w, r)
	case env.ConstantValues.ThreadPath:
		thread(w, r)
	case env.ConstantValues.CurrentUserPath:
		currentUser(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func thread(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		controller.CurrentController.CreateComment(w, r)
	case http.MethodGet:
		controller.CurrentController.GetComments(w, r)
	case http.MethodPut:
		controller.CurrentController.UpdateComment(w, r)
	case http.MethodDelete:
		controller.CurrentController.DeleteComment(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func currentUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		controller.CurrentController.CurrentUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
