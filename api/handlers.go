package api

import (
	"encoding/json"
	"mysub/models"
	"net/http"
)

func HandleGet(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(subsKey).(models.Subscription)
	if !ok {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
