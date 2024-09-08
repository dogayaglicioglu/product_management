package main

import (
	"fmt"
	"net/http"
	"product_management/database"
	"product_management/kafka"
	"product_management/redis"
	"product_management/repository"
	"product_management/routes"

	"github.com/gorilla/mux"
)

func main() {
	kafkaCreated := make(chan bool)
	dbCreated := make(chan bool)
	redisCreated := make(chan bool)
	go func() {
		database.ConnectDb(dbCreated)
	}()
	go func() {
		fmt.Print("INIT CONSUMER")
		kafka.InitConsumer(kafkaCreated)
	}()

	go func() {
		redis.ConnectRedis(redisCreated)
	}()
	<-dbCreated
	fmt.Print("DB IS CREATED...")
	<-kafkaCreated
	<-redisCreated
	db := database.GetDb()
	if db == nil {
		fmt.Println("db empty")
	} else {
		fmt.Println("db is initialized")
	}
	webRepo := repository.NewWebRepository(db)
	router := mux.NewRouter()
	routes.SetUpRoutes(router, webRepo)

	http.ListenAndServe(":8081", router)
}
