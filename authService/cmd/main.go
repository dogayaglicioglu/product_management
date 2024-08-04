package main

import (
	"auth-service/logger"
	"auth-service/middleware"
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

	routes.SetUpRoutes(router)
	http.ListenAndServe(":8082", router)
}
