package routes

import (
	handler "auth-service/handlers"
	"auth-service/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetUpRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", handler.Register).Methods("POST")
	r.HandleFunc("/auth/login", handler.Login).Methods("POST")
	r.HandleFunc("/auth/verify", handler.Verify).Methods("GET")
	r.Handle("/auth/changepassword", middleware.Authorization(http.HandlerFunc(handler.ChangePassword))).Methods("POST") //no sync. web db
	r.Handle("/auth/changeusername{username}", middleware.Authorization(http.HandlerFunc(handler.ChangeUsername))).Methods("POST")
	r.HandleFunc("/auth/update/{username}", handler.Update).Methods("PUT")
	r.HandleFunc("/auth/delete/{username}", handler.Delete).Methods("DELETE")
	//TO DO:
	//when changeUsername update and delete user web db must be sync. with kafka
}
