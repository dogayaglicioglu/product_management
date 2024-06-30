package main

import (
	"net/http"
	"product_management/database"
	"product_management/routes"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDb()
	router := mux.NewRouter()
	routes.SetUpRoutes(router)

	http.ListenAndServe(":8081", router)
}
