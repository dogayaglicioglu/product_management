package main

import (
	"fmt"
	"net/http"
	"product_management/database"
	"product_management/kafka"
	"product_management/routes"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDb()
	fmt.Print("DB IS CREATED...")
	go func() {
		fmt.Print("INIT CONSUMER")
		kafka.InitConsumer()
	}()
	router := mux.NewRouter()
	routes.SetUpRoutes(router)

	http.ListenAndServe(":8081", router)
}
