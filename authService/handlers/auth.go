package handler

import (
	"auth-service/database"
	"auth-service/models"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(database database.DbInstance) {
	db = database.DB
	db.AutoMigrate(&models.AuthUser{})

}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.AuthUser

	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Could not create user", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.AuthUser
	var input models.AuthUser
	json.NewDecoder(r.Body).Decode(&input)

	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token := "example-token"
	w.Write([]byte(token))

}
