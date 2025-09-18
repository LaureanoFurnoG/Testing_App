package main

import (
	"eventDriven/event-endpoint"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//start listening
	defer rabbitConn.Close()
	log.Println("CONNECTED TO RABBITMQ")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	topics := []string{"info", "test", "event"}
	err = consumer.Listen(topics)
	if err != nil{
		log.Fatal(err)
	}
	/// watch the queue and consume events
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbitis ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("RabbitMQ nt yet ready")
			counts++
		} else {
			connection = c
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off...")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
