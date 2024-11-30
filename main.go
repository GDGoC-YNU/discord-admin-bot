package main

import (
	"discord-admin-bot/pkg/secret"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type JoinFormSubmit struct {
	UserID string `json:"user_id"`
}

func main() {
	sec := secret.GetSecret()
	e := gin.Default()

	e.Static("/static", "./static")

	fsUserInfo := NewFirestoreUserInfoRepo()
	e.LoadHTMLGlob("templates/*")
	e.GET("/api/initial/form/debug", func(c *gin.Context) {
		c.HTML(200, "redirect.html.tmpl", gin.H{
			"RedirectURL": fmt.Sprintf("%s%s/api/initial/form", sec.System.Protocol, sec.System.Host),
			"OGPImageURL": fmt.Sprintf("%s%s/static/gdgoc-ynu-ogp-join.webp", sec.System.Protocol, sec.System.Host),
			"LogoURL":     fmt.Sprintf("%s%s/static/gdgoc_ynu_logo.webp", sec.System.Protocol, sec.System.Host),
			"Debug":       true,
		})
		return
	})

	e.GET("/api/initial/form", func(c *gin.Context) {
		redirectUrl := fmt.Sprintf(
			"https://discord.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify%%20guilds.members.read",
			sec.DiscordSecret.ClientID,
			sec.JoinForm.Callback,
		)
		c.HTML(200, "redirect.html.tmpl", gin.H{
			"RedirectURL": redirectUrl,
			"OGPImageURL": fmt.Sprintf("%s%s/static/gdgoc-ynu-ogp-join.webp", sec.System.Protocol, sec.System.Host),
			"LogoURL":     fmt.Sprintf("%s%s/static/gdgoc_ynu_logo.webp", sec.System.Protocol, sec.System.Host),
			"Debug":       false,
		})
		return
	})

	e.GET("/api/initial/form/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(400, gin.H{
				"message": "code is required",
			})
		}
		oauth2 := NewDiscordOAuth2Resolver()
		authInfo, err := oauth2.Resolve(code)
		if err != nil {
			log.Printf("failed to resolve, err: %v", err)
			c.JSON(500, gin.H{
				"message": "failed to resolve",
			})
			return
		}
		userID, err := fsUserInfo.SaveUserInfo(c, *authInfo)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "failed to save user info",
			})
			return
		}
		c.Redirect(302, fmt.Sprintf(sec.JoinForm.FormRedirectFormat, userID))
	})

	e.POST("/api/initial/form/submit", func(c *gin.Context) {
		var form JoinFormSubmit
		err := c.BindJSON(&form)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "failed to bind json",
			})
			return
		}
		// get user info from firestore
		userInfo, err := fsUserInfo.GetUserInfo(c, form.UserID)
		if err != nil {
			log.Printf("failed to get user info, err: %v", err)
			c.JSON(500, gin.H{
				"message": "failed to get user info",
			})
			return
		}
		d := NewDiscordServerlessClient()
		err = d.GrantRole(userInfo.Id, sec.DiscordSecret.MemberRoleID)
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
	})

	e.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	e.Run(":8080")
}
