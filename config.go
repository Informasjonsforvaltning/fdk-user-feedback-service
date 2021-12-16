package fdk_user_feedback_service

import (
	controller "github.com/Informasjonsforvaltning/fdk-user-feedback-service/controller"
	env "github.com/Informasjonsforvaltning/fdk-user-feedback-service/env"
	repository "github.com/Informasjonsforvaltning/fdk-user-feedback-service/repository"
	service "github.com/Informasjonsforvaltning/fdk-user-feedback-service/service"
)

func Configure() {
	repository.CurrentEntityRepository = &repository.EntityRepositoryImpl{
		SparqlServiceUrl: env.EnvironmentVariables.SparqlServiceUrl,
	}
	repository.CurrentThreadIdRepository = &repository.ThreadIdRepositoryImpl{
		FirestoreProjectId:    env.ConstantValues.FirestoreProjectId,
		FirestoreCollectionId: env.ConstantValues.FirestoreCollection,
	}
	repository.CurrentThreadRepository = &repository.ThreadRepositoryImpl{
		WriteApiToken:       env.EnvironmentVariables.WriteApiToken,
		ReadApiToken:        env.EnvironmentVariables.ReadApiToken,
		CommunityApiUrl:     env.EnvironmentVariables.CommunityApiUrl,
		ThreadBotUid:        env.EnvironmentVariables.ThreadBotUid,
		CommunityCategoryId: env.EnvironmentVariables.CommunityCategoryId,
		TopicPath:           env.ConstantValues.TopicPath,
		TopicsPath:          env.ConstantValues.TopicsPath,
		PostsPath:           env.ConstantValues.PostsPath,
	}
	repository.CurrentUserRepository = &repository.UserRepositoryImpl{
		ReadApiToken:     env.EnvironmentVariables.ReadApiToken,
		CommunityBaseUrl: env.EnvironmentVariables.CommunityApiUrl,
		UserByEmailPath:  env.ConstantValues.UserByEmailPath,
	}

	service.CurrentAuthService = &service.AuthServiceImpl{
		UserRepository: repository.CurrentUserRepository,
		KeycloakHost:   env.EnvironmentVariables.KeycloakHost,
	}
	service.CurrentEntityService = &service.EntityServiceImpl{
		EntityRepository: repository.CurrentEntityRepository,
	}
	service.CurrentThreadIdService = &service.ThreadIdServiceImpl{
		ThreadIdRepository: repository.CurrentThreadIdRepository,
	}
	service.CurrentThreadService = &service.ThreadServiceImpl{
		ThreadRepository: repository.CurrentThreadRepository,
		ThreadIdService:  service.CurrentThreadIdService,
		EntityService:    service.CurrentEntityService,
	}

	controller.CurrentController = &controller.ControllerImpl{
		AuthService:     service.CurrentAuthService,
		ThreadIdService: service.CurrentThreadIdService,
		ThreadService:   service.CurrentThreadService,
	}

}
