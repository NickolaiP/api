package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NickolaiP/api/config"
	"github.com/NickolaiP/api/dbinit"
	"github.com/NickolaiP/api/handlers"
	"github.com/NickolaiP/api/repository"

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

	if err := dbinit.InitDB(db, "./scripts/init.sql"); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	kafkaProducer, err := sarama.NewSyncProducer(conf.KafkaBrokers, nil)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	messageRepo := repository.NewMessageRepository(db)
	messageHandler := handlers.NewMessageHandler(messageRepo, kafkaProducer)

	router := mux.NewRouter()
	router.HandleFunc("/messages", messageHandler.CreateMessage).Methods("POST")
	router.HandleFunc("/stats", messageHandler.GetStats).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	log.Println("Server is running on port 8080")

	<-ctx.Done()

	stop()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
