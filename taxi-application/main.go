package main

import "github.com/segmentio/kafka-go"

type Service interface {
	Serve() error
}

func CreateProducer() *kafka.Writer {
	return nil
}

func main() {

}
