package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	KAFKA_BOOTSTRAP_SERVER = "192.168.88.129:9092"
	KAFKA_TOPIC            = "wiki"
)

func WriteMessagesToKafka(kafkaWriter *kafka.Writer, batchmessage [][]byte) error {
	err := kafkaWriter.WriteMessages(context.Background(), func(bmess [][]byte) []kafka.Message {
		result := []kafka.Message{}
		for _, v := range bmess {
			result = append(result, kafka.Message{Value: v})
		}
		return result
	}(batchmessage)...)
	return err
}

func main() {
	// init kafka producer
	p := &kafka.Writer{
		Addr:                   kafka.TCP(KAFKA_BOOTSTRAP_SERVER),
		Topic:                  KAFKA_TOPIC,
		RequiredAcks:           kafka.RequireAll,
		WriteTimeout:           time.Second * 1,
		Compression:            kafka.Gzip,
		AllowAutoTopicCreation: true,
	}

	// Stream data
	var streamURL string = "https://stream.wikimedia.org/v2/stream/recentchange"
	resp, err := http.Get(streamURL)
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

	var batch = [][]byte{}
	const SIZE_BATCH int = 100

	for {
		if scanio.Scan() {
			res := scanio.Bytes()
			batch = append(batch, res)
			if len(batch) < SIZE_BATCH {
				continue
			}
			err := WriteMessagesToKafka(p, batch)
			if err != nil {
				log.Println("Error: ", err)
			} else {
				log.Println("Wtite success")
			}
			batch = nil
		}
	}
}
