package handler

import (
	"discord-admin-bot/pkg/discord"
	"discord-admin-bot/repo"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type FormRedirectHandler struct {
	clientID, callback, protocol, host,
	formRedirectFormat,
	memberRoleID string
	oauth2   *discord.DiscordOAuth2Resolver
	userRepo *repo.FirestoreUserInfoRepo
}

func NewInitialFormHandler(
	clientID, callback, protocol, host,
	formRedirectFormat, memberRoleID string,
	oauth2 *discord.DiscordOAuth2Resolver,
	userRepo *repo.FirestoreUserInfoRepo,
) *FormRedirectHandler {
	return &FormRedirectHandler{
		clientID:           clientID,
		callback:           callback,
		protocol:           protocol,
		host:               host,
		formRedirectFormat: formRedirectFormat,
		memberRoleID:       memberRoleID,
		oauth2:             oauth2,
		userRepo:           userRepo,
	}
}

func (h FormRedirectHandler) Redirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectUrl := fmt.Sprintf(
			"https://discord.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify%%20guilds.members.read",
			h.clientID, h.callback,
		)
		c.HTML(200, "redirect.html.tmpl", gin.H{
			"RedirectURL": redirectUrl,
			"OGPImageURL": fmt.Sprintf("%s%s/static/gdgoc-ynu-ogp-join.webp", h.protocol, h.host),
			"LogoURL":     fmt.Sprintf("%s%s/static/gdgoc_ynu_logo.webp", h.protocol, h.host),
			"Debug":       false,
		})
		return
	}
}

func (h FormRedirectHandler) RedirectDebug() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "redirect.html.tmpl", gin.H{
			"RedirectURL": "https://www.example.com",
			"OGPImageURL": fmt.Sprintf("%s%s/static/gdgoc-ynu-ogp-join.webp", h.protocol, h.host),
			"LogoURL":     fmt.Sprintf("%s%s/static/gdgoc_ynu_logo.webp", h.protocol, h.host),
			"Debug":       true,
		})
		return
	}
}

func (h FormRedirectHandler) Callback() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(400, gin.H{
				"message": "code is required",
			})
		}
		authInfo, err := h.oauth2.Resolve(code)
		if err != nil {
			log.Printf("failed to resolve, err: %v", err)
			c.JSON(500, gin.H{
				"message": "failed to resolve",
			})
			return
		}
		userID, err := h.userRepo.SaveUserInfo(c, *authInfo)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "failed to save user info",
			})
			return
		}
		c.Redirect(302, fmt.Sprintf(h.formRedirectFormat, userID))
	}
}

func (h FormRedirectHandler) AcceptSubmit() gin.HandlerFunc {
	return func(c *gin.Context) {
		type JoinFormSubmit struct {
			UserID string `json:"user_id"`
		}

		var form JoinFormSubmit
		err := c.BindJSON(&form)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "failed to bind json",
			})
			return
		}
		// get user info from firestore
		userInfo, err := h.userRepo.GetUserInfo(c, form.UserID)
		if err != nil {
			log.Printf("failed to get user info, err: %v", err)
			c.JSON(500, gin.H{
				"message": "failed to get user info",
			})
			return
		}
		d := discord.NewDiscordServerlessClient()
		err = d.GrantRole(userInfo.Id, h.memberRoleID)
		if err != nil {
			log.Printf("failed to grant role, err: %v", err)
			c.JSON(500, gin.H{
				"message": "failed to grant role",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "ok",
		})
	}
}
