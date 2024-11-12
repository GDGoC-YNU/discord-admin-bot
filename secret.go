package main

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

func GetSecret() Secret {
	return secret
}

type Secret struct {
	DiscordSecret `yaml:"discord"`
	JoinForm      `yaml:"join_form"`
}

type DiscordSecret struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Token        string `yaml:"token"`
	GuildID      string `yaml:"guild_id"`
	MemberRoleID string `yaml:"member_role_id"`
}

type JoinForm struct {
	Callback           string `yaml:"callback_url"`
	FormRedirectFormat string `yaml:"form_redirect_format"`
}

var secret Secret

func init() {
	secretLocation := "./secret.yaml"
	if os.Getenv("SECRET_LOCATION") != "" {
		secretLocation = os.Getenv("SECRET_LOCATION")
	}
	secret = loadSecret(secretLocation)
}

func loadSecret(location string) Secret {
	var secret Secret
	reader, err := os.Open(location)
	if err != nil {
		log.Fatalf("failed to open config file, err: %v", err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("failed to read config file, err: %v", err)
	}
	err = yaml.Unmarshal(data, &secret)
	if err != nil {
		log.Fatalf("failed to unmarshal config file, err: %v", err)
	}
	return secret
}
