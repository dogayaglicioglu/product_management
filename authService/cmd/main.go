package main

import (
	"auth-service/logger"
	"auth-service/middleware"
	"auth-service/repository"
	"auth-service/routes"
	"net/http"

	"auth-service/database"

	"github.com/gorilla/mux"
)

func main() {
	dbCreated := make(chan bool)
	go func() {
		database.ConnectDb(dbCreated)
	}()

	logger.InitLog()
	<-dbCreated
	router := mux.NewRouter()
	router.Use(middleware.AccessLogger)
	authRepo := repository.NewAuthRepository(database.GetDb())

	routes.SetUpRoutes(router, authRepo)
	http.ListenAndServe(":8082", router)
}
