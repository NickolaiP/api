package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/NickolaiP/api/repository"
)

type MessageHandler struct {
	repo     repository.MessageRepository
	producer sarama.SyncProducer
}

func NewMessageHandler(repo repository.MessageRepository, producer sarama.SyncProducer) *MessageHandler {
	return &MessageHandler{
		repo:     repo,
		producer: producer,
	}
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var msg repository.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg.Processed = true

	if err := h.repo.Save(context.Background(), msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *MessageHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	count, err := h.repo.GetProcessedCount(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int{"processed_messages": count}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
