package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

const (
	smtpHost = "smtp-pulse.com"
	smtpPort = "2525"
	sender   = "William@minmax.com.tw"
)

type UploadData struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "https://funwilliam.github.io")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func sendMail(id, password, to, subject, body, base64File, encodeType, mediaType string) error {
	// smtp伺服器身分驗證
	auth := smtp.PlainAuth("", id, password, smtpHost)

	// 邊界標示符
	boundary := "my_unique_boundary"

	// 根據 MIME 類型確定檔案副檔名
	fileExtension := strings.SplitN(mediaType, "/", 2)[1]

	// 構建郵件內容
	from := sender
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: multipart/mixed; boundary=" + boundary + "\n\n" +
		"--" + boundary + "\n" +
		"Content-Type: text/plain; charset=utf-8\n\n" +
		body + "\n\n"

	// 圖片附檔
	if base64File != "" && encodeType != "" && mediaType != "" {
		msg += "--" + boundary + "\n" +
			"Content-Type: " + mediaType + "\n" +
			"Content-Transfer-Encoding: " + encodeType + "\n" +
			"Content-Disposition: attachment; filename=\"attachment." + fileExtension + "\"\n\n" +
			base64File + "\n\n"
	}

	// 郵件結束標示
	msg += "--" + boundary + "--"

	// 發送
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	return err
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var data UploadData
	var stmt string
	var mediaType string
	var encodeType string
	var base64File string

	enableCors(&w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 檢查是否同時沒有文字和圖片
	if data.Text == "" && data.Image == "" {
		http.Error(w, "Text and image are both empty", http.StatusBadRequest)
		return
	}

	log.Print(data.Text)
	log.Print(data.Image)
	// 如果有文字
	if data.Text != "" {
		stmt = data.Text
	}

	// 如果有圖片，解碼圖片數據(base64編碼 -> 二進制數據) + 取得圖片格式
	if data.Image != "" {
		dataParts := strings.FieldsFunc(data.Image, func(r rune) bool {
			return r == ';' || r == ','
		})

		if len(dataParts) != 3 {
			http.Error(w, "Invalid image data", http.StatusBadRequest)
			return
		}
		mediaType = dataParts[0]
		encodeType = dataParts[1]
		base64File = dataParts[2]
		log.Print(mediaType + " " + encodeType + " " + base64File)
	}

	// SMTP 配置
	login := "minmax.ed.notification.sys@gmail.com" // 更改为您的 Gmail 地址
	password := os.Getenv("SMTP_PASS")              // Gmail 密码或应用专用密码
	to := "minmax.ed.notification.sys@gmail.com"    // 更改为收件人地址
	subject := "測試"                                 // 邮件主题

	if err := sendMail(login, password, to, subject, stmt, base64File, encodeType, mediaType); err != nil {
		log.Printf("Failed to send email: %v\n", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}
	log.Print("receive event: " + stmt)
	log.Printf("Email sent successfully")
	fmt.Fprintf(w, "Email sent successfully")
}

func main() {
	http.HandleFunc("/", uploadHandler)
	fmt.Println("Server started at :443")
	log.Fatal(http.ListenAndServe(":443", nil))
}
