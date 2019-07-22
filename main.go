package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hongthang152/ShareItBackend/driver"
	"github.com/hongthang152/ShareItBackend/handler"
	"github.com/hongthang152/ShareItBackend/handler/middlewares"
	"github.com/hongthang152/ShareItBackend/handler/schedulers"
	"github.com/hongthang152/ShareItBackend/utils"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := driver.ConnectMongo(context.Background(), os.Getenv("MONGO_URL"))
	if err != nil {
		panic(err)
	}
	utils.SetUniqueIndex(client)

	FileHandler := handler.NewFileHandler(client)

	router := mux.NewRouter()
	router.Use(middlewares.SetCORSPolicy)
	router.Use(middlewares.Recovery)

	router.HandleFunc("/files", FileHandler.GetAll).Methods("GET", "OPTIONS")
	router.HandleFunc("/files/get", FileHandler.GetFile).Methods("GET", "OPTIONS")
	router.HandleFunc("/files/create", FileHandler.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/", FileHandler.CheckHealth).Methods("GET", "OPTIONS")

	fileScheduler := schedulers.NewFileScheduler(client)
	fileScheduler.DeleteExpiredFiles("*/30 * * * *")

	log.Println("Listening on port :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
