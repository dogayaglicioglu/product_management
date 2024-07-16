package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"product_management/database"
	models "product_management/models"

	"github.com/gorilla/mux"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Id is required to find user", http.StatusBadRequest)
		return
	}
	if err := database.DB.DB.First(&user, "username = ?", id); err != nil {

	}
	jsonResp, err := json.Marshal(user)
	if err != nil {
		fmt.Print("Error in marshal operation %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

}
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := database.DB.DB.Find(&users); err != nil {
		fmt.Errorf("Error in fetching users from db %v", err)

	}
	jsonResp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error in marshal operation.", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

}
