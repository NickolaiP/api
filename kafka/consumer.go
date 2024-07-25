package kafka

import (
	"database/sql"
	"log"

	"github.com/IBM/sarama"
)

func StartConsumer(db *sql.DB, brokers []string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("messages", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start consumer partition: %v", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Received message: %s", string(msg.Value))
			query := "UPDATE messages SET processed = $1 WHERE content = $2"
			_, err := db.Exec(query, true, string(msg.Value))
			if err != nil {
				log.Printf("Failed to update message status in database: %v", err)
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Error in consumer: %v", err)
		}
	}
}
