package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	models "product_management/models"
	"product_management/repository"

	"github.com/gorilla/mux"
)

var dbRepository repository.DbRepository

func SetUserRepo(dbRepo repository.DbRepository) {
	dbRepository = dbRepo

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Id is required to find user", http.StatusBadRequest)
		return
	}
	user, err := dbRepository.GetUserById(id)
	if err != nil {
		fmt.Print("Error getting user from db.. %v", err)

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
	users, err := dbRepository.GetAllUsers()
	if err != nil {
		fmt.Print("Error getting users from db.. %v", err)
	}
	jsonResp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error in marshal operation.", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

}
