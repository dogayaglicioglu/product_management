package main

import (
	handler "auth-service/handlers"
	"auth-service/middleware"
	"auth-service/routes"
	"net/http"

	"auth-service/database"

	"github.com/gorilla/mux"
)

func main() {
	//loggerInstance := logger.Logger()
	dbInst := database.ConnectDb()

	//fLogger := loggerInstance.Logger
	handler.InitDb(dbInst)

	router := mux.NewRouter()
	router.Use(middleware.LoggerMiddleware) //add middleware to the router

	routes.SetUpRoutes(router)
	http.ListenAndServe(":8082", router)
}
