package main

import (
	"bufio"
	"context"
	"log"
	"net/http"

	"github.com/segmentio/kafka-go"
)

const (
	KAFKA_BOOTSTRAP_SERVER = "192.168.88.129:9092"
	KAFKA_TOPIC            = "wiki"
)

func WriteMessageToKafka(kafkaWriter *kafka.Writer, message []byte) error {
	err := kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   nil,
		Value: message,
	})
	return err
}

func main() {
	// init kafka producer
	p := &kafka.Writer{
		Addr:     kafka.TCP(KAFKA_BOOTSTRAP_SERVER),
		Topic:    KAFKA_TOPIC,
		Balancer: &kafka.LeastBytes{},
	}

	// Stream data
	var stream string = "https://stream.wikimedia.org/v2/stream/eventgate-main.test.event"
	resp, err := http.Get(stream)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	defer resp.Body.Close()

	// Create a scanner to read response
	scanio := bufio.NewScanner(resp.Body)

	// Split response if data contain double newline character
	scanio.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		for i := 0; i < len(data)-2; i++ {
			if data[i] == '\n' && data[i+1] == '\n' {
				return i + 2, data[0:i], nil
			}
		}

		if atEOF {
			return len(data), data, nil
		}
		return
	})
	for {
		if scanio.Scan() {
			res := scanio.Bytes()
			err := WriteMessageToKafka(p, res)
			if err != nil {
				log.Println("Error: ", err)
			} else {
				log.Println("Wtite success")
			}
		}
	}
}
