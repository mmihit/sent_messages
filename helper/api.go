package helper

import (
	"encoding/json"
	"net/http"
)

func (api *ApiResponse) Sent(w http.ResponseWriter) {
	err := json.NewEncoder(w).Encode(api)
	if err != nil {
		http.Error(w, "error encoding data to front-end", http.StatusInternalServerError)
	}
}