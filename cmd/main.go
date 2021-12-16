package main

import (
	"context"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	fdk_user_feedback_service "github.com/Informasjonsforvaltning/fdk-user-feedback-service"
)

func main() {
	ctx := context.Background()
	if err := funcframework.RegisterHTTPFunctionContext(ctx, "/", fdk_user_feedback_service.EntryPoint); err != nil {
		log.Fatalf("funcframework.RegisterHTTPFunctionContext: %v\n", err)
	}
	// Use PORT environment variable, or default to 8000.
	port := "8000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	} else {
		log.Printf("Bringing user-feedback-service up on localhost:%s\n", port)
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
