package main

import (
	"discord-admin-bot/handler"
	"discord-admin-bot/pkg/discord"
	"discord-admin-bot/pkg/secret"
	"discord-admin-bot/repo"
	"github.com/gin-gonic/gin"
)

func main() {
	sec := secret.GetSecret()
	e := gin.Default()

	e.Static("/static", "./static")

	fsUserInfo := repo.NewFirestoreUserInfoRepo()
	e.LoadHTMLGlob("templates/*")
	oauth2 := discord.NewDiscordOAuth2Resolver()

	formHandler := handler.NewInitialFormHandler(
		sec.ClientID, sec.JoinForm.Callback, sec.System.Protocol, sec.System.Host,
		sec.JoinForm.FormRedirectFormat, sec.DiscordSecret.MemberRoleID,
		oauth2, fsUserInfo)

	e.GET("/api/initial/form/debug", formHandler.RedirectDebug())
	e.GET("/api/initial/form", formHandler.Redirect())
	e.GET("/api/initial/form/callback", formHandler.Callback())
	e.POST("/api/initial/form/submit", formHandler.AcceptSubmit())

	e.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	e.Run(":8080")
}
