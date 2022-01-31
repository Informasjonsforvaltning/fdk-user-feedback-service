package tests

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type MockResponseWriter struct {
	MockHeader         map[string][]string
	MockStatusCode     int
	MockError          error
	CurrentStatusCode  int
	CurrentWriteOutput []byte
}

func (m *MockResponseWriter) Header() http.Header {
	return m.MockHeader
}
func (m *MockResponseWriter) Write(bytes []byte) (int, error) {
	m.CurrentWriteOutput = bytes
	return m.MockStatusCode, m.MockError
}
func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.CurrentStatusCode = statusCode
}

var MockRsaKey, _ = rsa.GenerateKey(rand.Reader, 2048)

func MockJwkStore() *httptest.Server {
	key, err := jwk.New(MockRsaKey)
	if err != nil {
		fmt.Println(err)
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Add("Content-Type", "application/json")

				key.Set(jwk.KeyIDKey, "testkid")

				buf, err := json.MarshalIndent(key, "", "  ")
				if err != nil {
					fmt.Printf("failed to marshal key into JSON: %s\n", err)
					return
				}

				fmt.Fprintf(rw, `{"keys":[%s]}`, buf)
			},
		),
	)

	return server
}

func CreateMockJwt(expiresAt int64, email *string, audience *[]string) *string {
	t := jwt.New()
	t.Set(jwt.SubjectKey, `https://github.com/lestrrat-go/jwx/jwt`)
	t.Set(jwt.IssuedAtKey, time.Now().Unix())
	t.Set(jwt.ExpirationKey, expiresAt)
	t.Set(`email`, email)
	if audience != nil {
		t.Set(jwt.AudienceKey, *audience)
	}

	jwk_key, _ := jwk.New(MockRsaKey)

	jwk_key.Set(jwk.KeyIDKey, "testkid")

	signed, _ := jwt.Sign(t, jwa.RS256, jwk_key)

	signed_string := string(signed)

	return &signed_string

}
