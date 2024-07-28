package main

import (
	handler "auth-service/handlers"
	"auth-service/logger"
	"auth-service/middleware"
	"auth-service/routes"
	"net/http"

	"auth-service/database"

	"github.com/gorilla/mux"
)

func main() {
	dbInst := database.ConnectDb()
	logger.InitLog()
	handler.InitDb(dbInst)
	router := mux.NewRouter()
	router.Use(middleware.AccessLogger)

	routes.SetUpRoutes(router)
	http.ListenAndServe(":8082", router)
}
