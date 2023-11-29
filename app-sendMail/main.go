package main

import (
	"fmt"
	"net/http"
)

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析表单数据
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// 提取图片和字符串
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	message := r.FormValue("message")

	// TODO: 实现发送邮件的逻辑

	fmt.Fprintf(w, "File and message received. Message: %s", message)
}

func main() {
	http.HandleFunc("/upload", handleUpload)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
