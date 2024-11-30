package discord

import (
	"discord-admin-bot/pkg/secret"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
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
	UserInfo
}

type GuildMemberStatusResponse struct {
	Nick     string   `json:"nick"`
	UserInfo UserInfo `json:"user"`
}

type UserInfo struct {
	Id                   string      `json:"id"`
	Username             string      `json:"username"`
	Avatar               string      `json:"avatar"`
	Discriminator        string      `json:"discriminator"`
	PublicFlags          int         `json:"public_flags"`
	Flags                int         `json:"flags"`
	Banner               string      `json:"banner"`
	AccentColor          interface{} `json:"accent_color"`
	GlobalName           string      `json:"global_name"`
	AvatarDecorationData struct {
		Asset     string      `json:"asset"`
		SkuId     string      `json:"sku_id"`
		ExpiresAt interface{} `json:"expires_at"`
	} `json:"avatar_decoration_data"`
	BannerColor  interface{} `json:"banner_color"`
	Clan         interface{} `json:"clan"`
	PrimaryGuild interface{} `json:"primary_guild"`
}

func (r DiscordOAuth2Resolver) Resolve(code string) (authInfo *AuthInfo, err error) {
	tokens, err := r.getTokens(code)
	if err != nil {
		log.Printf("failed to get tokens, err: %v", err)
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
	authInfo.UserInfo = &meResp.UserInfo
	mem, err := r.GetGuildMemberStatus(tokens.AccessToken, secret.GetSecret().DiscordSecret.GuildID, meResp.UserInfo.Id)
	if err != nil {
		return nil, err
	}
	authInfo.GuildMemberStatusResponse = mem
	return authInfo, nil
}

func (r DiscordOAuth2Resolver) GetMe(accessToken string) (*MeResponse, error) {
	targetUrl := "https://discord.com/api/users/@me"
	rt := resty.New()
	resp, err := rt.R().
		SetBasicAuth(r.clientID, r.clientSecret).
		SetAuthToken(accessToken).
		Get(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, err: %v", err)
	}
	b := resp.Body()
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get me, status: %d, body: %s", resp.StatusCode(), string(b))
	}
	var data MeResponse
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (r DiscordOAuth2Resolver) GetGuildMemberStatus(accessToken, guildID, userID string) (*GuildMemberStatusResponse, error) {
	targetUrl := fmt.Sprintf("https://discord.com/api/users/@me/guilds/%s/member", guildID)
	rt := resty.New()
	resp, err := rt.
		R().
		SetAuthToken(accessToken).
		Get(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, err: %v", err)
	}
	b := resp.Body()
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get guild member status, targetUrl: %s, status: %d, body: %s",
			targetUrl, resp.StatusCode(), string(b))
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
	targetUrl := "https://discord.com/api/oauth2/token"
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", r.callbackUrl)
	rt := resty.New()
	resp, err := rt.R().
		SetBasicAuth(r.clientID, r.clientSecret).
		SetFormData(map[string]string{
			"grant_type":   "authorization_code",
			"code":         code,
			"redirect_uri": r.callbackUrl,
		}).
		Post(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, err: %v", err)
	}
	b := resp.Body()
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get tokens, status: %d, body: %s", resp.StatusCode(), string(b))
	}
	var data oAuth2ResponseData
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
