package main

import (
	"discord-admin-bot/pkg/secret"
	"fmt"
	"github.com/gin-gonic/gin"
)

type JoinFormSubmit struct {
	Email string `json:"email"`
}

func main() {
	sec := secret.GetSecret()
	e := gin.Default()

	fsUserInfo := NewFirestoreUserInfoRepo()

	e.GET("/api/initial/form", func(c *gin.Context) {
		redirectUrl := fmt.Sprintf(
			"https://discord.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify%%20guilds.members.read",
			sec.DiscordSecret.ClientID,
			sec.JoinForm.Callback,
		)
		c.Redirect(302, redirectUrl)
	})

	e.GET("/api/initial/form/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(400, gin.H{
				"message": "code is required",
			})
		}
		oauth2 := NewDiscordOAuth2Resolver()
		tokens, err := oauth2.getTokens(code)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "failed to get tokens",
			})
			return
		}
		authInfo, err := oauth2.Resolve(tokens.AccessToken)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "failed to resolve",
			})
			return
		}
		//TODO: store authenticated info
		userID, err := fsUserInfo.SaveUserInfo(c, authInfo)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "failed to save user info",
			})
			return
		}
		c.Redirect(302, fmt.Sprintf(sec.JoinForm.FormRedirectFormat, userID))
	})
	e.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	e.Run(":8080")
}
