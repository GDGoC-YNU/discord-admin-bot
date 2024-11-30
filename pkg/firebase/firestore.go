package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"discord-admin-bot/pkg/secret"
	"google.golang.org/api/option"
	"log"
	"os"
)

var (
	client *firestore.Client
)

func New() *firestore.Client {
	if client == nil {
		var err error
		ctx := context.Background()
		config := secret.GetSecret().Firestore
		var opt option.ClientOption
		if config.CredentialPath != "" {
			log.Printf("config.CredentialPath: %v", config.CredentialPath)
			data, err := os.ReadFile(config.CredentialPath)
			if err != nil {
				log.Fatalf("os.ReadFile err: %v", err)
			}
			opt = option.WithCredentialsJSON(data)
		}
		client, err = firestore.NewClient(ctx, config.ProjectID, opt)
		if err != nil {
			log.Fatalf("firebase.NewClient err: %v", err)
		}
	}
	return client
}
