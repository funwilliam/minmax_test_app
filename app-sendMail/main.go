package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type UploadData struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}

const serviceAccountFile = "service-account.json"

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "https://funwilliam.github.io")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func sendMail(service *gmail.Service, from, to, subject, bodyMessage, base64Image string) error {
	var message gmail.Message

	// MIME邮件头部和正文
	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "multipart/mixed; boundary=boundary"

	var msg strings.Builder
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	// 正文部分
	msg.WriteString("--boundary\r\n")
	msg.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
	msg.WriteString(bodyMessage + "\r\n")

	// 图像附件
	msg.WriteString("--boundary\r\n")
	msg.WriteString("Content-Type: image/png\r\n")
	msg.WriteString("Content-Transfer-Encoding: base64\r\n")
	msg.WriteString("Content-Disposition: attachment; filename=\"image.png\"\r\n\r\n")
	msg.WriteString(base64Image + "\r\n")
	msg.WriteString("--boundary--")

	// Base64 encode the email
	message.Raw = base64.URLEncoding.EncodeToString([]byte(msg.String()))

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
	return gmail.NewService(ctx, option.WithHTTPClient(client))
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
	if err := sendMail(gmailService, "id-001@cloud-notify-sys.iam.gserviceaccount.com", "minmax.ed.notification.sys@gmail.com", "測試", data.Text, data.Image); err != nil {
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
