package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/repository"
	gocloak "github.com/Nerzal/gocloak/v10"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	AuthenticateAndGetUser(jwt string) (*model.User, int)
	AuthenticateJwt(jwt string) (*jwt.MapClaims, int)
	GetUser(email string) (*model.User, error)
}

type AuthServiceImpl struct {
	UserRepository repository.UserRepository
	KeycloakHost   string
}

func (authService *AuthServiceImpl) AuthenticateAndGetUser(jwt string) (*model.User, int) {
	claims, statusCode := authService.AuthenticateJwt(jwt)
	if statusCode != 200 {
		return nil, statusCode
	}

	user, err := authService.GetUser(fmt.Sprint((*claims)["email"]))
	if err != nil && (user == nil || user.UserId == nil) {
		return nil, http.StatusUnauthorized
	}

	return user, http.StatusOK
}

func (authService *AuthServiceImpl) AuthenticateJwt(jwt string) (*jwt.MapClaims, int) {
	client := gocloak.NewClient(authService.KeycloakHost)

	ctx := context.Background()
	_, claims, err := client.DecodeAccessToken(ctx, jwt, "fdk")

	if err != nil {
		log.Println(err.Error())
		return nil, http.StatusUnauthorized
	}

	if claims != nil && !claims.VerifyAudience("fdk-feedback-service", true) {
		return nil, http.StatusForbidden
	}

	if claims == nil || (*claims)["email"] == nil {
		return nil, http.StatusForbidden
	}

	return claims, http.StatusOK
}

func (authService *AuthServiceImpl) GetUser(email string) (*model.User, error) {
	user, err := authService.UserRepository.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

var CurrentAuthService AuthService
