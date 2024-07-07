package main

import (
	handler "auth-service/handlers"
	"auth-service/routes"
	"net/http"

	"auth-service/database"

	"github.com/gorilla/mux"
)

func main() {
	dbInst := database.ConnectDb()
	handler.Init(dbInst)

	router := mux.NewRouter()
	routes.SetUpRoutes(router)
	http.ListenAndServe(":8082", router)
}
