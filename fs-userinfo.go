package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"discord-admin-bot/pkg/firebase"
	"log"
)

type FSUserInfo struct {
	Id            string `json:"id" firestore:"id"`
	Username      string `json:"username" firestore:"username"`
	Avatar        string `json:"avatar" firestore:"avatar"`
	Discriminator string `json:"discriminator" firestore:"discriminator"`
	GlobalName    string `json:"global_name" firestore:"global_name"`
	Nickname      string `json:"nick" firestore:"nick"`
}

type FirestoreUserInfoRepo struct {
	fs *firestore.Client
}

func NewFirestoreUserInfoRepo() *FirestoreUserInfoRepo {
	return &FirestoreUserInfoRepo{
		fs: firebase.New(),
	}
}

func (r FirestoreUserInfoRepo) SaveUserInfo(ctx context.Context, d AuthInfo) (string, error) {
	u := FSUserInfo{
		Id:            d.UserInfo.Id,
		Username:      d.UserInfo.Username,
		Avatar:        d.UserInfo.Avatar,
		Discriminator: d.UserInfo.Discriminator,
		GlobalName:    d.UserInfo.GlobalName,
		Nickname:      d.GuildMemberStatusResponse.Nick,
	}
	newDoc := r.fs.Collection("users").NewDoc()
	_, err := newDoc.Set(ctx, u)
	if err != nil {
		log.Printf("failed to save user info, err: %v", err)
		return "", err
	}
	log.Printf("saved user info, doc: %v", newDoc.ID)
	return newDoc.ID, nil
}

func (r FirestoreUserInfoRepo) GetUserInfo(ctx context.Context, userID string) (*FSUserInfo, error) {
	doc, err := r.fs.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		log.Printf("failed to get user info, err: %v", err)
		return nil, err
	}
	var u FSUserInfo
	if err := doc.DataTo(&u); err != nil {
		log.Printf("failed to convert data, err: %v", err)
		return nil, err
	}
	return &u, nil
}
