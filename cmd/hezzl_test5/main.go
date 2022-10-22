package main

import (
	"context"
	"github.com/GeorgeShibanin/hezzl_test5/internal/config"
	"github.com/GeorgeShibanin/hezzl_test5/internal/handlers"
	"github.com/GeorgeShibanin/hezzl_test5/internal/storage"
	"github.com/GeorgeShibanin/hezzl_test5/internal/storage/postgres"
	"github.com/GeorgeShibanin/hezzl_test5/internal/storage/rediscachedstorage"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	srv := NewServer()
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func NewServer() *http.Server {
	r := mux.NewRouter()

	var store storage.Storage
	var err error

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Redis_URL,
	})

	store = initPostgres()
	store, err = rediscachedstorage.Init(redisClient, store)
	if err != nil {
		log.Fatalf("can't init postgres connection: %s", err.Error())
	}

	handler := handlers.NewHTTPHandler(store)
	r.HandleFunc("/item/create", handler.HandlePostItem).Methods(http.MethodPost)
	r.HandleFunc("/item/update", handler.HandlePatchItem).Methods(http.MethodPatch)
	r.HandleFunc("/item/remove", handler.HandleDeleteItem).Methods(http.MethodDelete)
	r.HandleFunc("/item/list", handler.HandleGetItem).Methods(http.MethodGet)

	return &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func initPostgres() *postgres.StoragePostgres {
	store, err := postgres.Init(
		context.Background(),
		config.PostgresHost,
		config.PostgresUser,
		config.PostgresDB,
		config.PostgresPassword,
		config.PostgresPort,
	)
	if err != nil {
		log.Fatalf("can't init postgres connection: %s", err.Error())
	}
	return store
}
