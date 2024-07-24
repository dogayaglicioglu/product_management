package handler

import (
	"auth-service/database"
	"auth-service/kafka"
	"auth-service/models"
	"auth-service/verify"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
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

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	if username != "" {
		http.Error(w, "The username must be entered.", http.StatusBadRequest)
		return
	}

	var foundedUser models.AuthUser
	result := db.Where("username = ?", username).First(&foundedUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "The user is not found, so you can't delete it", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching user from the database", http.StatusInternalServerError)
		}
		return
	}
	if err := db.Delete(&foundedUser).Error; err != nil {
		http.Error(w, "Error in deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"] //exists username
	if username != "" {
		http.Error(w, "The username must be entered.", http.StatusBadRequest)
		return
	}

	var updatedUser models.AuthUser
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "The user couldn't found.", http.StatusInternalServerError)
		return
	}
	var foundedUser models.AuthUser
	if err := db.Where("username = ?", username).First(&foundedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching user from the database", http.StatusInternalServerError)
		}
		return
	}
	if updatedUser.Username != "" {
		foundedUser.Username = updatedUser.Username
	}
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		foundedUser.Password = string(hashedPassword)

	}
	foundedUser.Password = updatedUser.Password
	if err := db.Save(&foundedUser).Error; err != nil {
		http.Error(w, "Error updating user in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func Verify(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	verified := verify.VerifyToken(tokenStr)
	if verified != false {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Token verified"))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Token verification failed."))
	}

}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var updatedUser models.AuthUser
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var checkUser models.AuthUser
	result := db.Where("username = ?", updatedUser.Username).First(&checkUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("The user does not exist, you cant change password..")
			return
		} else {
			// another error is occured
			http.Error(w, "Error checking user registration", http.StatusInternalServerError)
			return
		}
	}
	checkUser.Password = updatedUser.Password
	if err := db.Save(&checkUser).Error; err != nil {
		http.Error(w, "Error while updating password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Password updated successfully.")

}

func ChangeUsername(w http.ResponseWriter, r *http.Request) {
	var updatedUser models.AuthUser
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//check the user is exist ?
	var checkUser models.AuthUser
	result := db.Where("username = ?", updatedUser.Username).First(&checkUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("The user does not exist, you cant change username..")
			return
		} else {
			// another error is occured
			http.Error(w, "Error checking user registration", http.StatusInternalServerError)
			return
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error in generating hashedPassword", http.StatusInternalServerError)
		return
	}
	checkUser.Password = string(hashedPassword)
	if err := db.Save(&checkUser).Error; err != nil {
		http.Error(w, "Error while updating password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Password updated successfully.")
}
func Register(w http.ResponseWriter, r *http.Request) {
	var authUser models.AuthUser
	if err := json.NewDecoder(r.Body).Decode(&authUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	//check whether the user is already registered
	var existingUser models.AuthUser
	result := db.Where("username = ?", authUser.Username).First(&existingUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("The user does not exist")
		} else {
			// another error is occured
			http.Error(w, "Error checking user registration", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	//if the user is not registered, register it
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
	fmt.Print("THE USER IS CREATED...")

	//sync web db in here
	kafka.ProduceEvent(authUser.Username)
	fmt.Print("THE MESSAGE SUCCES. SENT..")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User is registered."))
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

	w.Write([]byte("Successfully logged in."))

}
