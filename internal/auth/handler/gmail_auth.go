package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type GmailAuthHandler struct {
	config *oauth2.Config
}

func NewGmailAuthHandler() (*GmailAuthHandler, error) {
	var config *oauth2.Config

	// 1. ลองอ่านจากไฟล์ก่อน
	credsPath := os.Getenv("GMAIL_CREDENTIALS_PATH")
	if credsPath != "" {
		b, err := os.ReadFile(credsPath)
		if err == nil {
			config, err = google.ConfigFromJSON(b, gmail.GmailSendScope)
		}
	}

	// 2. ถ้าไม่มีไฟล์ หรืออ่านไม่ได้ ให้ใช้จาก Env
	if config == nil {
		clientID := os.Getenv("GMAIL_CLIENT_ID")
		clientSecret := os.Getenv("GMAIL_CLIENT_SECRET")
		if clientID == "" || clientSecret == "" {
			return nil, fmt.Errorf("GMAIL_CREDENTIALS_PATH or GMAIL_CLIENT_ID/SECRET not set")
		}

		config = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			Scopes:       []string{gmail.GmailSendScope},
			// ต้องใส่ RedirectURL ให้ตรงกับที่ตั้งใน Google Console ด้วย
			RedirectURL: os.Getenv("GMAIL_REDIRECT_URL"),
		}
		if config.RedirectURL == "" {
			config.RedirectURL = "http://localhost:8080/api/auth/gmail/callback"
		}
	}

	return &GmailAuthHandler{config: config}, nil
}

// GET /api/auth/gmail
func (h *GmailAuthHandler) Login(c *gin.Context) {
	url := h.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GET /api/auth/gmail/callback
func (h *GmailAuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code is missing"})
		return
	}

	tok, err := h.config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve token from web"})
		return
	}

	// บันทึก Token ลงไฟล์
	tokenPath := os.Getenv("GMAIL_TOKEN_PATH")
	if tokenPath == "" {
		tokenPath = "token.json"
	}

	f, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to cache oauth token"})
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(tok)

	c.JSON(http.StatusOK, gin.H{"message": "Gmail authentication successful! token.json has been generated."})
}
