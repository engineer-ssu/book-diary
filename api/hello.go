package handler

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Vercel이 호출하는 엔트리포인트 함수 (반드시 Handler여야 함)
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := Response{
		Message: "Hello from Go on Vercel!",
		Status:  http.StatusOK,
	}

	json.NewEncoder(w).Encode(res)
}