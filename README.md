## How to use this project
Clone the repo

### Setting up environment
- Copy `.env.example` file and paste it as `.env` in the same path
- In `.env` file Add your email and password and leave other values as it is
- Add your email address in `EMAIL_TO` in which you want to receive your test email
- For using OAUTH2 to send email and filling up environment variables `CLIENT_ID`, `CLIENT_SECRET`, `ACCESS_TOKEN` and `REFRESH_TOKEN` follow my medium article [Sending Emails with GO (Golang) Using SMTP, Gmail, and OAuth2](https://medium.com/wesionary-team/sending-emails-with-go-golang-using-smtp-gmail-and-oauth2-185ee12ab306)

### Sending email
- For sending test email using SMTP, execute following command
```
go run main.go SMTP
``` 

- For sending test email using Gmail API and OAUTH2, execute following command
```
go run main.go OAUTH
```

NOTE: Environment variables must be filled properly to be able to execute above commands