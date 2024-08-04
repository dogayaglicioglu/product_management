package handler

import (
	"auth-service/database"
	"auth-service/kafka"
	"auth-service/logger"
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

var jwtKey = []byte("my_secret")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Delete(w http.ResponseWriter, r *http.Request) {
	loggerInst := r.Context().Value(logger.LoggerKey).(*logger.LogInstance)
	vars := mux.Vars(r)
	username := vars["username"]
	if username != "" {
		loggerInst.Error(r.Context(), "The username must be entered.")
		http.Error(w, "The username must be entered.", http.StatusBadRequest)
		return
	}

	var foundedUser models.AuthUser
	result := database.DB.DB.Where("username = ?", username).First(&foundedUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			loggerInst.Error(r.Context(), "The user is not found, so you can't delete it", result.Error)
			http.Error(w, "The user is not found, so you can't delete it", http.StatusNotFound)
		} else {
			loggerInst.Error(r.Context(), "Error fetching user from the database", result.Error)
			http.Error(w, "Error fetching user from the database", http.StatusInternalServerError)
		}
		return
	}
	if err := database.DB.DB.Delete(&foundedUser).Error; err != nil {
		loggerInst.Error(r.Context(), "Error in deleting user", err)
		http.Error(w, "Error in deleting user", http.StatusInternalServerError)
		return
	}
	req := models.RequestPayload{
		OldUsername: username,
	}
	kafka.ProduceEvent(req, "delete-user")
	loggerInst.Info(r.Context(), "The user succ. deleted.")
	w.WriteHeader(http.StatusOK)

}

func Update(w http.ResponseWriter, r *http.Request) {
	loggerInst := r.Context().Value(logger.LoggerKey).(*logger.LogInstance)
	vars := mux.Vars(r)
	username := vars["username"] //exists username
	if username != "" {
		loggerInst.Error(r.Context(), "The username must be entered.")
		http.Error(w, "The username must be entered.", http.StatusBadRequest)
		return
	}

	var updatedUser models.AuthUser
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		loggerInst.Error(r.Context(), "The user couldn't found.", err)
		http.Error(w, "The user couldn't found.", http.StatusInternalServerError)
		return
	}
	var foundedUser models.AuthUser
	if err := database.DB.DB.Where("username = ?", username).First(&foundedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			loggerInst.Error(r.Context(), "User not found.", err)
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			loggerInst.Error(r.Context(), "Error fetching user from the database.", err)
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
			loggerInst.Error(r.Context(), "Error hashing password", err)
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		foundedUser.Password = string(hashedPassword)

	}
	foundedUser.Password = updatedUser.Password
	if err := database.DB.DB.Save(&foundedUser).Error; err != nil {
		loggerInst.Error(r.Context(), "Error updating user in the database", err)
		http.Error(w, "Error updating user in the database", http.StatusInternalServerError)
		return
	}
	req := models.RequestPayload{
		NewUsername: foundedUser.Username,
		OldUsername: username,
	}
	//sync. web db
	kafka.ProduceEvent(req, "update-user") //only send username field.
	loggerInst.Info(r.Context(), "The user is updated succ.")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The updated successsfully."))

}
func Verify(w http.ResponseWriter, r *http.Request) {
	loggerInst := r.Context().Value(logger.LoggerKey).(*logger.LogInstance)
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		loggerInst.Error(r.Context(), "Token required.")
		return
	}

	verified := verify.VerifyToken(tokenStr)
	if verified != false {
		w.WriteHeader(http.StatusOK)
		loggerInst.Info(r.Context(), "Token verified.")
		w.Write([]byte("Token verified"))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		loggerInst.Error(r.Context(), "Token verification failed.")
		w.Write([]byte("Token verification failed."))
	}

}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	loggerInst := r.Context().Value(logger.LoggerKey).(*logger.LogInstance)
	var updatedUser models.AuthUser
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		loggerInst.Error(r.Context(), "Invalid request payload.", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var checkUser models.AuthUser
	result := database.DB.DB.Where("username = ?", updatedUser.Username).First(&checkUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			loggerInst.Error(r.Context(), "The user does not exist, you cant change password..", result.Error)
			http.Error(w, "The user does not exist, you cant change password..", http.StatusInternalServerError)
			return
		} else {
			// another error is occured
			loggerInst.Error(r.Context(), "Error checking user registration", result.Error)
			http.Error(w, "Error checking user registration", http.StatusInternalServerError)
			return
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		loggerInst.Error(r.Context(), "Error in generating hashed password.")
		return
	}
	checkUser.Password = string(hashedPassword)
	if err := database.DB.DB.Save(&checkUser).Error; err != nil {
		loggerInst.Error(r.Context(), "Error while updating password.", err)
		http.Error(w, "Error while updating password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password updated successsfully."))
	loggerInst.Info(r.Context(), "Password updated successfully.")

}

func ChangeUsername(w http.ResponseWriter, r *http.Request) {
	loggerInst := r.Context().Value(logger.LoggerKey).(*logger.LogInstance)
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		loggerInst.Error(r.Context(), "The username must be entered.")
		http.Error(w, "The username must be entered.", http.StatusBadRequest)
		return
	}
	var newUsernamePayload struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newUsernamePayload); err != nil {
		loggerInst.Error(r.Context(), "Error in decode operation. %v", err)
	}
	newUsername := newUsernamePayload.Username
	//check the user is exist ?
	var existsUser models.AuthUser
	result := database.DB.DB.Where("username = ?", username).First(&existsUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			loggerInst.Error(r.Context(), "The user does not exist, you cant change username..", result.Error)
			http.Error(w, "The user does not exist, you cant change username..", http.StatusNotFound)
			return
		} else {
			// another error is occured
			loggerInst.Error(r.Context(), "Error checking user registration.", result.Error)
			http.Error(w, "Error checking user registration.", http.StatusInternalServerError)
			return
		}
	}
	var duplicateUser models.AuthUser
	result = database.DB.DB.Where("username = ?", newUsername).First(&duplicateUser)
	if result.Error == nil {
		loggerInst.Error(r.Context(), "The new username is already taken.")
		http.Error(w, "The new username is already taken.", http.StatusConflict)
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		loggerInst.Error(r.Context(), "Error checking new username.", result.Error)
		http.Error(w, "Error checking new username.", http.StatusInternalServerError)
		return
	}

	// Update the username
	existsUser.Username = newUsername
	if err := database.DB.DB.Save(&existsUser).Error; err != nil {
		loggerInst.Error(r.Context(), "Error while updating username.", err)
		http.Error(w, "Error while updating username.", http.StatusInternalServerError)
		return
	}

	webUser := models.RequestPayload{
		OldUsername: username,
		NewUsername: existsUser.Username,
	}

	// Sync web db in here
	kafka.ProduceEvent(webUser, "change-username")
	loggerInst.Info(r.Context(), "The message was successfully sent to Kafka.")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The username updated successfully."))
	loggerInst.Info(r.Context(), "Username updated successfully.")
}

func Register(w http.ResponseWriter, r *http.Request) {
	loggerInst := r.Context().Value(logger.LoggerKey).(*logger.LogInstance)
	var authUser models.AuthUser
	if err := json.NewDecoder(r.Body).Decode(&authUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		loggerInst.Error(r.Context(), "Invalid request payload", err)
		return
	}
	//check whether the user is already registered
	var existingUser models.AuthUser
	result := database.DB.DB.Where("username = ?", authUser.Username).First(&existingUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("The user does not exist")
		} else {
			// another error is occured
			http.Error(w, "Error checking user registration", http.StatusInternalServerError)
			loggerInst.Error(r.Context(), "Error checking user registration", result.Error)
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
		loggerInst.Error(r.Context(), "Error in generating hashed password.")
		return
	}

	authUser.Password = string(hashedPassword)
	if err := database.DB.DB.Create(&authUser).Error; err != nil {
		http.Error(w, "Could not create user", http.StatusBadRequest)
		loggerInst.Error(r.Context(), "Could not create user")
		return
	}

	loggerInst.Info(r.Context(), "The user is registered.")

	//sync web db in here
	kafka.ProduceEvent(authUser.Username, "register-user")
	loggerInst.Info(r.Context(), "The msg succ. sent to kafka...")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User is registered."))
}

func Login(w http.ResponseWriter, r *http.Request) {
	loggerInst := r.Context().Value(logger.LoggerKey).(*logger.LogInstance)
	var user models.AuthUser
	var input models.AuthUser
	json.NewDecoder(r.Body).Decode(&input)

	if err := database.DB.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		http.Error(w, "There is no such user.", http.StatusUnauthorized)
		loggerInst.Error(r.Context(), "There is no such user.", err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		loggerInst.Error(r.Context(), "Invalid username or password.", err)

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
		loggerInst.Error(r.Context(), "Internal Error.", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	loggerInst.Info(r.Context(), "Successfully logged in.")

	w.Write([]byte("Successfully logged in."))

}
