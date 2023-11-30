package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UploadData struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // 允许任何源
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// 预检请求的处理
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

	fmt.Fprintf(w, "Received text and image")
}

func main() {
	http.HandleFunc("/", uploadHandler)
	fmt.Println("Server started at :443")
	http.ListenAndServe(":443", nil)
}
