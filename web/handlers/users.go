package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"product_management/database"
	models "product_management/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

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
