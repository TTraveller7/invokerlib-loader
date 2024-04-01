package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	address := "kafka:9092"
	topic := "orderline_source"
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	file, err := os.OpenFile("order-line.csv", os.O_RDONLY, 0644)
	if err != nil {
		err = fmt.Errorf("open file failed: err=%v", err)
		log.Printf("%v", err)
		return
	}

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{address}, producerConfig)
	if err != nil {
		err = fmt.Errorf("create producer failed: %v", err)
		log.Printf("%v", err)
		return
	}
	defer producer.Close()

	count := 0
	startTime := time.Now()
	batchSize := 10
	s := bufio.NewScanner(file)
	messages := make([]*sarama.ProducerMessage, 0, batchSize)
	for s.Scan() {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(s.Bytes()),
		}
		messages = append(messages, msg)
		if len(messages) == batchSize {
			if err := producer.SendMessages(messages); err != nil {
				log.Printf("send messages failed: %v", err)
				return
			}
			messages = make([]*sarama.ProducerMessage, 0, batchSize)
			count += batchSize
			if count%10000 == 0 {
				log.Printf("%v messages submitted. Elapsed time: %v", count, time.Since(startTime))
			}
		}
	}
	if len(messages) > 0 {
		if err := producer.SendMessages(messages); err != nil {
			log.Printf("send messages failed: %v", err)
			return
		}
	}
	count += len(messages)
	log.Printf("produce finished: final count=%v", count)
	time.Sleep(1000000)
}
