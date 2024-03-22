package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

type OrderService struct {
	Config ConfigOrderService
}

type ConfigOrderService struct {
	ServerUrl  string
	ServerPort string
}

type Location struct {
	Long float64
	Lat  float64
}

type OrderRequest struct {
	UserId             string
	Position           Location
	Destination        Location
	AddressPosition    string
	AddressDestination string
}

func NewOrderService(config ConfigOrderService) *OrderService {
	return &OrderService{Config: config}
}

func orderRequestHandle(r *gin.Context) {
	var request OrderRequest

	err := r.ShouldBindJSON(&request)
	if err != nil {
		log.Println("Order Service Handle Request: ", err)
	}

	// Create produceer
	kafkaUrl := os.Getenv("kafka.bootstrapserver")
	kafkaTopic := os.Getenv("kafka.topic.order")
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{Brokers: []string{kafkaUrl}, Topic: kafkaTopic})

	// write message to kafka
	message, err := json.Marshal(request)
	if err != nil {
		log.Println("Order Service Handle Request invalid value: ", err)
	}
	kafkaWriter.WriteMessages(context.Background(), kafka.Message{Key: []byte(request.UserId), Value: message})
}

func (s *OrderService) Serve() error {
	// Run a gin rest api server
	r := gin.Default()

	// handler request
	r.POST("/order", orderRequestHandle)

	// listen request order from user
	err := r.Run(fmt.Sprintf("%v:%v", s.Config.ServerUrl, s.Config.ServerPort))
	if err != nil {
		log.Println("Order Service: ", err)
		return err
	}
	return nil
}
