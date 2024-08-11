package handler

import (
	"auth-service/kafka"
	"auth-service/logger"
	"auth-service/models"
	"auth-service/repository"
	"auth-service/verify"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var authRepository repository.AuthRepository

func SetAuthRepo(authRepo repository.AuthRepository) {
	authRepository = authRepo
}

func Delete(w http.ResponseWriter, r *http.Request) {
	loggerInst := r.Context().Value(logger.LoggerKey).(*logger.LogInstance)
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		loggerInst.Error(r.Context(), "The username must be entered.")
		http.Error(w, "The username must be entered.", http.StatusBadRequest)
		return
	}

	err, errMsg := authRepository.DeleteUser(username, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		loggerInst.Error(r.Context(), errMsg, err)
		return
	}

	kafka.ProduceEvent(username, "delete-user") //only send the string value
	loggerInst.Info(r.Context(), "The user succ. deleted.")
	msg := "The user " + username + " deleted successfully."
	w.Write([]byte(msg))
	w.WriteHeader(http.StatusOK)

}

func UpdateUsernameAndPasswd(w http.ResponseWriter, r *http.Request) {
	loggerInst := r.Context().Value(logger.LoggerKey).(*logger.LogInstance)
	vars := mux.Vars(r)
	username := vars["username"] //exists username

	if username == "" {
		loggerInst.Error(r.Context(), "The username must be entered.")
		http.Error(w, "The username must be entered.", http.StatusBadRequest)
		return
	}

	var updatedUser models.AuthUser //bodyde gönderilen yeni username ve password bilgisi
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		loggerInst.Error(r.Context(), "The user couldn't found.", err)
		http.Error(w, "The user couldn't found.", http.StatusInternalServerError)
		return
	}

	if updatedUser.Username == "" && updatedUser.Password == "" {
		loggerInst.Error(r.Context(), "The new username and password must be entered.")
		http.Error(w, "The new username and password must be entered.", http.StatusBadRequest)
		return
	}

	foundedUser, errMsg, err := authRepository.UpdateUser(updatedUser, username, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		loggerInst.Error(r.Context(), errMsg, err)
		return

	}
	webUser := models.RequestPayload{
		NewUsername: foundedUser.Username,
		OldUsername: username,
	}
	//sync. web db
	kafka.ProduceEvent(webUser, "update-user") //send the structure
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

	errMsg, err := authRepository.ChangePassword(updatedUser, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		loggerInst.Error(r.Context(), errMsg, err)
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
	err, errMsg := authRepository.ChangeUsername(username, newUsername, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		loggerInst.Error(r.Context(), errMsg, err)
		return
	}
	webUser := models.RequestPayload{
		OldUsername: username,
		NewUsername: newUsername,
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
	err, errMsg := authRepository.Register(authUser, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		loggerInst.Error(r.Context(), errMsg, err)
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

	var input models.AuthUser
	json.NewDecoder(r.Body).Decode(&input)
	tokenStr, errMsg, expirationTime, err := authRepository.Login(input, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusExpectationFailed)
		loggerInst.Error(r.Context(), errMsg, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expirationTime,
	})

	loggerInst.Info(r.Context(), "Successfully logged in.")
	w.Write([]byte("Successfully logged in."))

}
