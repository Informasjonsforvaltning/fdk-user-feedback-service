package env

import "os"

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type Environment struct {
	CommunityApiUrl     string
	CommunityCategoryId string
	ThreadBotUid        string
	ReadApiToken        string
	WriteApiToken       string
	SparqlServiceUrl    string
	KeycloakHost        string
	FdkBaseUri          string
	FirestoreCollection string
}

type Constants struct {
	PingPath           string
	CurrentUserPath    string
	ThreadPath         string
	UserByEmailPath    string
	TopicPath          string
	TopicsPath         string
	ThreadSlugPath     string
	PostsPath          string
	FirestoreProjectId string
}

var EnvironmentVariables = Environment{
	CommunityApiUrl:     getEnv("COMMUNITY_API_URL", "https://community.staging.fellesdatakatalog.digdir.no/api/"),
	CommunityCategoryId: getEnv("COMMUNITY_CATEGORY_ID", "25"),
	ThreadBotUid:        getEnv("TOPIC_BOT_UID", "1"),
	ReadApiToken:        getEnv("READ_API_TOKEN", ""),
	WriteApiToken:       getEnv("WRITE_API_TOKEN", ""),
	SparqlServiceUrl:    getEnv("SPARQL_SERVICE_URL", "https://sparql.staging.fellesdatakatalog.digdir.no"),
	KeycloakHost:        getEnv("KEYCLOAK_HOST", "https://sso.staging.fellesdatakatalog.digdir.no/"),
	FdkBaseUri:          getEnv("FDK_BASE_URI", "https://www.staging.fellesdatakatalog.digdir.no/"),
	FirestoreCollection: getEnv("FIRESTORE_COLLECTION", "threadIds_staging"),
}

var ConstantValues = Constants{
	PingPath:           "ping",
	CurrentUserPath:    "current-user",
	ThreadPath:         "thread",
	UserByEmailPath:    "/user/email/",
	TopicPath:          "/topic/",
	TopicsPath:         "/v3/topics/",
	ThreadSlugPath:     "/thread-slug/",
	PostsPath:          "/v3/posts/",
	FirestoreProjectId: "digdir-cloud-functions",
}
