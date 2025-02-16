package handlers

import (
	"encoding/json"
	"net/http"

	"Avito_Merch_project/internal/middleware"
	"Avito_Merch_project/internal/models"
	"Avito_Merch_project/internal/repository"
)

var repoInstance *repository.Repository

func SetRepository(repo *repository.Repository) {
	repoInstance = repo
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	coins, inventory, coinHistory, err := repoInstance.GetUserInfo(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := models.InfoResponse{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
