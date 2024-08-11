package main

import (
	"fmt"
	"net/http"
	"product_management/database"
	"product_management/kafka"
	"product_management/repository"
	"product_management/routes"

	"github.com/gorilla/mux"
)

func main() {
	kafkaCreated := make(chan bool)
	dbCreated := make(chan bool)
	go func() {
		database.ConnectDb(dbCreated)

	}()

	fmt.Print("DB IS CREATED...")
	go func() {
		fmt.Print("INIT CONSUMER")
		kafka.InitConsumer(kafkaCreated)
	}()
	<-dbCreated
	<-kafkaCreated
	webRepo := repository.NewWebRepository(database.GetDb())
	router := mux.NewRouter()
	routes.SetUpRoutes(router, webRepo)

	http.ListenAndServe(":8081", router)
}
