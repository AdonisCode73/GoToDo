package internal

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func NewFirestore(ctx context.Context) (*firebase.App, *firestore.Client, func(), error) {
	// Set an environment variable with your firestore credentials
	credsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	creds := option.WithCredentialsFile(credsPath)
	projectID := os.Getenv("PROJECT_ID")

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, creds)
	if err != nil {
		return nil, nil, nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	cleanUp := func() {
		_ = client.Close()
	}

	return app, client, cleanUp, nil
}
