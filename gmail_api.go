package gomail

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GmailService : Gmail client for sending email
var GmailService *gmail.Service

var authConfig *oauth2.Config

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (oauth2.TokenSource, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	var tokenSource oauth2.TokenSource

	tokFile, err := filepath.Abs("./token.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		getTokenFromWeb(config)
		return tokenSource, err
	}
	tokenSource = config.TokenSource(context.Background(), tok)
	return tokenSource, nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

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

// InitializeGmailAPI initialize gmail api
func InitializeGmailAPI() {
	credentialFilePath, err := filepath.Abs("./oauth_credentials.json")
	if err != nil {
		panic("Unable to load oauth_credentials.json file")
	}

	b, err := ioutil.ReadFile(credentialFilePath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	authConfig = config

	tokenSource, err := getClient(config)
	if err != nil {
		log.Printf("Unable to get token: %v", err)
	}

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		fmt.Println("Email service is initialized \n")
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
	// InitializeGmailAPI()
}

func SendEmail(to string, emailSubject string, emailSubjectData interface{}, emailSubjectTemplateName string, emailBodyData interface{}, emailBodyTemplateBodyName string) (bool, error) {

	isEmailSubjectEmpty := isInterfaceEmpty(emailSubjectData)
	var err error
	if !isEmailSubjectEmpty {
		emailSubject, err = parseTemplate(emailSubjectTemplateName, emailSubjectData)
		if err != nil {
			return false, errors.New("unable to parse email subject template")
		}
	}

	emailBody, err := parseTemplate(emailBodyTemplateBodyName, emailBodyData)
	if err != nil {
		return false, errors.New("unable to parse email body template")
	}

	var message gmail.Message

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	emailTo := "To: " + to + "\r\n"
	subject := "Subject: " + emailSubject + "!\n"
	msg := []byte(emailTo + subject + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err = GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Message sent!")
	}
	return true, nil
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("api/utils/email_templates/%s", templateFileName))
	if err != nil {
		return "", errors.New("invalid template name")
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}

func isInterfaceEmpty(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
