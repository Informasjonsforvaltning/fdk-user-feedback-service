package repository

import (
	"log"
	"net/http"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/util"
)

type UserRepository interface {
	GetByEmail(email string) (*model.User, error)
}

type UserRepositoryImpl struct {
	ReadApiToken     string
	CommunityBaseUrl string
	UserByEmailPath  string
}

func (userRepository *UserRepositoryImpl) GetByEmail(email string) (*model.User, error) {
	bearerToken := userRepository.ReadApiToken
	method := http.MethodGet
	endpointUrl := userRepository.CommunityBaseUrl + userRepository.UserByEmailPath + email

	response, err := util.Request(util.RequestOptions{
		Method:      method,
		EndpointUrl: endpointUrl,
		AccessToken: &bearerToken,
	})
	if err != nil {
		log.Println("Error on request.\n[ERROR] -", err, endpointUrl)
		return nil, err
	}
	user, err := util.UnmarshalUser(response)

	return user, err
}

var CurrentUserRepository UserRepository
