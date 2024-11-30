package main

import (
	"discord-admin-bot/pkg/secret"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type DiscordOAuth2Resolver struct {
	callbackUrl, clientID, clientSecret string
}

func NewDiscordOAuth2Resolver() *DiscordOAuth2Resolver {
	sec := secret.GetSecret()
	return &DiscordOAuth2Resolver{
		callbackUrl:  sec.JoinForm.Callback,
		clientID:     sec.DiscordSecret.ClientID,
		clientSecret: sec.DiscordSecret.ClientSecret,
	}
}

type AuthInfo struct {
	*UserInfo
	*GuildMemberStatusResponse
}

type MeResponse struct {
	User UserInfo `json:"user"`
}

type GuildMemberStatusResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

type UserInfo struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	GlobalName    string `json:"global_name"`
	PublicFlags   int    `json:"public_flags"`
}

func (r DiscordOAuth2Resolver) Resolve(code string) (authInfo *AuthInfo, err error) {
	tokens, err := r.getTokens(code)
	if err != nil {
		return nil, err
	}
	meResp, err := r.GetMe(tokens.AccessToken)
	if err != nil {
		return nil, err
	}
	if meResp == nil {
		return nil, fmt.Errorf("failed to get me")
	}
	authInfo = new(AuthInfo)
	authInfo.UserInfo = &meResp.User
	mem, err := r.GetGuildMemberStatus(tokens.AccessToken, secret.GetSecret().DiscordSecret.GuildID, meResp.User.Id)
	if err != nil {
		return nil, err
	}
	authInfo.GuildMemberStatusResponse = mem
	return authInfo, nil
}

func (r DiscordOAuth2Resolver) GetMe(accessToken string) (*MeResponse, error) {
	targetUrl := "https://discord.com/api/users/@me"
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request, err: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, err: %v", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response, err: %v", err)
	}
	var data MeResponse
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (r DiscordOAuth2Resolver) GetGuildMemberStatus(accessToken, guildID, userID string) (*GuildMemberStatusResponse, error) {
	targetUrl := fmt.Sprintf("https://discord.com/api/guilds/%s/members/%s", guildID, userID)
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request, err: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, err: %v", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response, err: %v", err)
	}
	var data GuildMemberStatusResponse
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

type oAuth2ResponseData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (r DiscordOAuth2Resolver) getTokens(code string) (*oAuth2ResponseData, error) {
	targetUrl := fmt.Sprintf("%s?code=%s", "https://discord.com/api/oauth2/token", code)
	//basic auth
	req, err := http.NewRequest("POST", targetUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request, err: %v", err)
	}
	//send request
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", r.callbackUrl)
	form.Add("client_id", r.clientID)
	form.Add("client_secret", r.clientSecret)
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, err: %v", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response, err: %v", err)
	}
	var data oAuth2ResponseData
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
