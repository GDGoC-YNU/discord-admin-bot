package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"discord-admin-bot/pkg/firebase"
	"log"
)

type FirestoreUserInfoRepo struct {
	fs *firestore.Client
}

func NewFirestoreUserInfoRepo() *FirestoreUserInfoRepo {
	return &FirestoreUserInfoRepo{
		fs: firebase.New(),
	}
}

func (r FirestoreUserInfoRepo) SaveUserInfo(ctx context.Context, d *AuthInfo) (string, error) {
	newDoc := r.fs.Collection("users").NewDoc()
	_, err := newDoc.Set(ctx, d)
	if err != nil {
		log.Printf("failed to save user info, err: %v", err)
		return "", err
	}
	log.Printf("saved user info, doc: %v", newDoc.ID)
	return newDoc.ID, nil
}
