package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"Avito_Merch_project/internal/middleware"
	"Avito_Merch_project/internal/services"
)

var merchService *services.MerchService

func SetMerchService(service *services.MerchService) {
	merchService = service
}

func BuyHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	item := vars["item"]
	err := merchService.BuyMerch(username, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Покупка прошла успешно"}`))
}
