package main

import (
	"discord-admin-bot/pkg/secret"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type DiscordServerlessClient struct {
	secret secret.DiscordSecret
}

func NewDiscordServerlessClient() *DiscordServerlessClient {
	sec := secret.GetSecret()
	return &DiscordServerlessClient{
		secret: sec.DiscordSecret,
	}
}

func (d DiscordServerlessClient) GrantUserRole(userID string) error {
	return d.GrantRole(userID, d.secret.MemberRoleID)
}

func (d DiscordServerlessClient) GrantRole(userID, roleID string) error {
	c, err := d.getClient()
	if err != nil {
		return err
	}
	defer c.Close()
	if err := c.GuildMemberRoleAdd(d.secret.GuildID, userID, roleID); err != nil {
		return fmt.Errorf("failed to grant role, err: %v", err)
	}
	return nil
}

func (d DiscordServerlessClient) getClient() (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + d.secret.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create discord client, err: %v", err)
	}
	return discord, nil
}
