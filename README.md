# FDK User Feedback Service

This application provides an API for handling feedback on resources in the data catalog. Designed to be hosted as a
cloud function.

For a broader understanding of the system’s context, refer to
the [architecture documentation](https://github.com/Informasjonsforvaltning/architecture-documentation) wiki. For more
specific context on this application, see the **Portal** subsystem section.

## Getting Started

These instructions will give you a copy of the project up and running on your local machine for development and testing
purposes.

### Prerequisites

Ensure you have the following installed:

- Go
- Firebase CLI (https://firebaseopensource.com/projects/firebase/firebase-tools/)

Clone the repository.

```sh
git clone https://github.com/Informasjonsforvaltning/fdk-user-feedback-service.git
cd fdk-user-feedback-service
```

#### Expose environment variables

```sh
export READ_API_TOKEN="{admin read token for Datalandsbyen}"
export WRITE_API_TOKEN="{admin write token for Datalandsbyen}"
export FIRESTORE_EMULATOR_HOST="localhost:8080"
```

#### Start firebase emulator

```
firebase emulators:start
```

#### Install required dependencies

```shell
go get
```

#### Start application

```sh
go run cmd/main.go
```

### Running tests

```shell
go test ./...
```

To generate a test coverage report, use the following command:

```shell
go test -v -race -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...
```