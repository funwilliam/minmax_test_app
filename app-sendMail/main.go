package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 创建一个结构体来匹配前端发送的JSON数据
type UploadData struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var data UploadData
	// 解析JSON数据
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 处理数据（例如打印到控制台）
	fmt.Printf("Received text: %s\n", data.Text)
	fmt.Printf("Received image: %s\n", data.Image)

	// 发送响应回前端
	fmt.Fprintf(w, "Received text and image")
}

func main() {
	http.HandleFunc("/", uploadHandler)
	fmt.Println("Server started at :80")
	http.ListenAndServe(":80", nil)
}
