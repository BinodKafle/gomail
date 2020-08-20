package gomail

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("Unable to encode token")
	}
}

func GetAuthorizationCodeCallback(c *gin.Context) {

	code := c.Request.FormValue("code")
	if code == "" {
		log.Println("Unable to receive code")
		return
	}

	log.Println("checking code", code)
	log.Printf("chekcing authCOnfig %v", authConfig)

	tok, err := authConfig.Exchange(context.Background(), code)

	if err != nil {
		log.Printf("Unable to retrieve token from web: %v", err)
	}

	tokFile, err := filepath.Abs("./token.json")
	saveToken(tokFile, tok)
}
