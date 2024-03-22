package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type EstimateCostStreamService struct{}

type EstimateCost struct {
	UserId string
	Cost   string
}

func CaculateCost(order OrderRequest) EstimateCost {
	return EstimateCost{UserId: order.UserId, Cost: "10000"}
}

func (s *EstimateCostStreamService) Serve() error {
	// Create a consummer to reciv order
	kafkaUrl := os.Getenv("kafka.bootstrapserver")
	kafkaTopicOrder := os.Getenv("kafka.topic.order")
	kafkaConsumer := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{kafkaUrl}, Topic: kafkaTopicOrder})

	// Create a producer to post cost of delivery
	kafkaTopicCost := os.Getenv("kafka.topic.estimatevalue")
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{Brokers: []string{kafkaUrl}, Topic: kafkaTopicCost})

	for {
		mess, err := kafkaConsumer.ReadMessage(context.Background())
		if err != nil {
			log.Println("Estimate Cost:", err)
			return err
		}
		var req OrderRequest
		json.Unmarshal(mess.Value, &req)
		cost, _ := json.Marshal(CaculateCost(req))
		kafkaWriter.WriteMessages(context.Background(), kafka.Message{Key: []byte(req.UserId), Value: cost})
	}
}
