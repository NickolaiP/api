package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/NickolaiP/api/config"
	"github.com/NickolaiP/api/dbinit"
	"github.com/NickolaiP/api/handlers"
	"github.com/NickolaiP/api/kafka"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", conf.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация базы данных
	if err := dbinit.InitDB(db, "./scripts/init.sql"); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	kafkaProducer, err := sarama.NewSyncProducer(conf.KafkaBrokers, nil)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	go kafka.StartConsumer(db, conf.KafkaBrokers)

	router := mux.NewRouter()
	router.HandleFunc("/messages", handlers.CreateMessageHandler(db, kafkaProducer)).Methods("POST")
	router.HandleFunc("/stats", handlers.GetStatsHandler(db)).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
