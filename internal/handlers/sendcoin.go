package handlers

import (
	"encoding/json"
	"net/http"

	"Avito_Merch_project/internal/middleware"
	"Avito_Merch_project/internal/models"
	"Avito_Merch_project/internal/services"
)

var coinService *services.CoinService

func SetCoinService(service *services.CoinService) {
	coinService = service
}

func SendCoinHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req models.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := coinService.SendCoins(username, req.ToUser, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Монетки успешно отправлены"}`))
}
