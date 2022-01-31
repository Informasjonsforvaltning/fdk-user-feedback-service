# fdk-user-feedback-service
A microservice handling feedback on resources in the data catalog. Designed to be hosted as a cloud function.

## To run locally:
* install Golang >=1.16
* install firebase CLI https://firebaseopensource.com/projects/firebase/firebase-tools/
* expose environment variables:
```
export READ_API_TOKEN="{admin read token for Datalandsbyen}"
export WRITE_API_TOKEN="{admin write token for Datalandsbyen}"
export FIRESTORE_EMULATOR_HOST="localhost:8080"
```
* Tests:
```
// Run all tests
go test ./...
// Run all tests and generate coverage report
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
// Print generated coverage report per file
go tool cover -func profile.cov
```
* Start firebase emulator:
```
firebase emulators:start
```
* Start the service:
```
go get
go run cmd/main.go
```
* Check that service is running:
```
curl -X 'GET' \
  'localhost:8000/ping' \
  -H 'accept: application/json'
```
