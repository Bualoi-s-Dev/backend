package configs

import (
	"context"
	"encoding/json"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func GetFirebaseCloudCredentials() string {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create secret manager client: %v", err)
	}
	defer client.Close()

	secretName := os.Getenv("PROD_PRIVATE_KEY_SECRET_NAME")
	if secretName == "" {
		log.Fatalf("SECRET_NAME environment variable not set")
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatalf("Failed to access secret version: %v", err)
	}

	var credentials map[string]interface{}
	if err := json.Unmarshal(result.Payload.Data, &credentials); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	credentialsPath := os.Getenv("PROD_FIREBASE_CREDENTIALS_PATH")
	if credentialsPath == "" {
		log.Fatalf("CREDENTIALS_PATH environment variable not set")
	}

	if err := os.WriteFile(credentialsPath, result.Payload.Data, 0644); err != nil {
		log.Fatalf("Failed to write credentials file: %v", err)
	}

	return credentialsPath
}

func GetFirebaseLocalCredentials() string {
	return os.Getenv("LOCAL_FIREBASE_CREDENTIALS_PATH")
}

func InitializeFirebaseAuth() *auth.Client {
	var credentialPath string
	if os.Getenv("APP_MODE") == "production" {
		credentialPath = GetFirebaseCloudCredentials()
	} else {
		credentialPath = GetFirebaseLocalCredentials()
	}
	opt := option.WithCredentialsFile(credentialPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err)
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error getting Firebase Auth client: %v", err)
	}
	return authClient
}
