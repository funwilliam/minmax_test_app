package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type UploadData struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}

const serviceAccountFile = "service-account.json"

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func sendMail(service *gmail.Service, from, to, subject, bodyMessage string) error {
	var message gmail.Message

	emailBody := fmt.Sprintf("Subject: %s\n\n%s\n\nAttached Image:\n%s", subject, bodyMessage, bodyMessage)

	// Base64 encode the email
	message.Raw = base64.URLEncoding.EncodeToString([]byte(emailBody))

	// Send the message
	_, err := service.Users.Messages.Send("me", &message).Do()
	return err
}

func createGmailService() (*gmail.Service, error) {
	ctx := context.Background()

	// Load the service account key
	data, err := os.ReadFile(serviceAccountFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read service account file: %v", err)
	}

	// Authenticate with the service account
	conf, err := google.JWTConfigFromJSON(data, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("unable to acquire JWT config: %v", err)
	}
	client := conf.Client(ctx)

	// Create the Gmail client
	return gmail.New(client)
}

func uploadHandler(w http.ResponseWriter, r *http.Request, gmailService *gmail.Service) {
	enableCors(&w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var data UploadData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received text: %s\n", data.Text)
	fmt.Printf("Received image: %s\n", data.Image)

	// Send email using Gmail API
	if err := sendMail(gmailService, "your-email@gmail.com", "recipient-email@example.com", "測試", data.Text); err != nil {
		log.Printf("Failed to send email: %v\n", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Received text and image, email sent")
}

func main() {
	gmailService, err := createGmailService()
	if err != nil {
		log.Fatalf("Failed to create Gmail service: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uploadHandler(w, r, gmailService)
	})
	fmt.Println("Server started at :443")
	log.Fatal(http.ListenAndServe(":443", nil))
}
