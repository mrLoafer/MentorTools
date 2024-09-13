package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type RequestData struct {
	Number int `json:"number"`
}

type ResponseData struct {
	Result int `json:"result"`
}

func incrementHandler(w http.ResponseWriter, r *http.Request) {
	var reqData RequestData
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result := reqData.Number + 1

	resData := ResponseData{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resData)
}

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.HandleFunc("/increment", incrementHandler)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
