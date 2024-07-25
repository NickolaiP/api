package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/IBM/sarama"
)

type Message struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Processed bool   `json:"processed"`
}

func CreateMessageHandler(db *sql.DB, kafkaProducer sarama.SyncProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Устанавливаем processed в true перед записью в базу данных
		msg.Processed = true

		// Insert message into the database
		query := "INSERT INTO messages (content, processed) VALUES ($1, $2) RETURNING id"
		err := db.QueryRow(query, msg.Content, msg.Processed).Scan(&msg.ID)
		if err != nil {
			log.Printf("Failed to insert message into database: %v", err)
			http.Error(w, "Failed to save message", http.StatusInternalServerError)
			return
		}

		// Send message to Kafka
		kafkaMsg := &sarama.ProducerMessage{
			Topic: "messages",
			Value: sarama.StringEncoder(msg.Content),
		}
		_, _, err = kafkaProducer.SendMessage(kafkaMsg)
		if err != nil {
			log.Printf("Failed to send message to Kafka: %v", err)
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)
	}
}

func GetStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var count int
		query := `SELECT COUNT(*) FROM messages WHERE processed = TRUE`
		err := db.QueryRow(query).Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"processed_messages": count})
	}
}
