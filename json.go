package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithJSON: 响应错误信息
func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	type errRespose struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errRespose{
		Error: msg,
	})
}

// respondWithJSON: 响应结果JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON respose : %v ", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
