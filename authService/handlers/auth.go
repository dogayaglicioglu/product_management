package handler

import (
	"auth-service/database"
	"auth-service/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

var jwtKey = []byte("my_secret")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Init(database database.DbInstance) {
	db = database.DB
	db.AutoMigrate(&models.AuthUser{})

}

func Verify(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token verified"))

}

func Register(w http.ResponseWriter, r *http.Request) {
	var authUser models.AuthUser
	json.NewDecoder(r.Body).Decode(&authUser)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	authUser.Password = string(hashedPassword)
	if err := db.Create(&authUser).Error; err != nil {
		http.Error(w, "Could not create user", http.StatusBadRequest)
		return
	}
	var foundedUser models.AuthUser
	if err := db.Where("username = ?", authUser.Username).First(&foundedUser).Error; err != nil {
		fmt.Print("Error in here %v", err)
		http.Error(w, "User not found after creation", http.StatusNotFound)
		return
	}

	user := models.User{
		AuthUserID: foundedUser.ID,
		Username:   foundedUser.Username,
	}

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
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.Write([]byte(tokenString))

}
