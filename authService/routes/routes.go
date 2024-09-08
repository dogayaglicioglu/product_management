package routes

import (
	handler "auth-service/handlers"
	"auth-service/middleware"
	"auth-service/repository"
	"net/http"

	"github.com/gorilla/mux"
)

func SetUpRoutes(r *mux.Router, authRepo repository.AuthRepository) {
	handler.SetAuthRepo(authRepo)
	r.HandleFunc("/auth/register", handler.Register).Methods("POST") //sync. web db
	r.HandleFunc("/auth/login", handler.Login).Methods("POST")
	r.HandleFunc("/auth/verify", handler.Verify).Methods("GET")
	r.Handle("/auth/changepassword", middleware.Authorization(http.HandlerFunc(handler.ChangePassword))).Methods("POST")            //no sync. web db
	r.Handle("/auth/changeusername/{username}", middleware.Authorization(http.HandlerFunc(handler.ChangeUsername))).Methods("POST") //sync. web db
	r.Handle("/auth/update/{username}", middleware.Authorization(http.HandlerFunc(handler.UpdateUsernameAndPasswd))).Methods("PUT") //sync. web db
	r.Handle("/auth/delete/{username}", middleware.Authorization(http.HandlerFunc(handler.Delete))).Methods("DELETE")               //sync. web db
}
