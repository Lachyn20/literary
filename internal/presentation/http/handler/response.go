package handler

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Message string     `json:"message,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := JSONResponse{Status: "ok", Data: data}
	_ = json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := JSONResponse{Status: "error", Message: message}
	_ = json.NewEncoder(w).Encode(resp)
}
