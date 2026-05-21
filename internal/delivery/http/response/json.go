package response

import (
	"encoding/json"
	"net/http"
)

func ResponseEntity(w http.ResponseWriter, code int, payload any) {
	if payload == nil {
		w.WriteHeader(code)
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":500,"message":"internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
