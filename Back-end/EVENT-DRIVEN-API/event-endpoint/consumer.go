package event

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	//queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()

	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()

	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name          string `json:"name"`
	Data          string `json:"data"`
	CorrelationId string `json:"correlationId"`
	ReplyTo       string `json:"replyTo"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	ch.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)

	q, _ := ch.QueueDeclare("", false, true, true, false, nil)

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [exchange: %s, queue: %s]\n", "logs_topic", q.Name)
	<-forever
	return nil
}

type DataParsed struct {
	Id_Group    int
	HttpType    string
	Url         string
	RequestType string
	Request     map[string]interface{}
	Header      map[string]interface{}
	Token       string
}

func handlePayload(payload Payload) {
	fmt.Println(payload)
	log.Print(payload)
	switch payload.Name {
	case "test", "event":
		// log whatever we get
		result, err := testApp(payload)
		if err != nil {
			log.Panic(err)
		}
		log.Println(result)
	}
}

func testApp(dataParseds Payload) (result string, err error) {
	var dataParsed DataParsed
	err = json.Unmarshal([]byte(dataParseds.Data), &dataParsed)
	if err != nil {
		panic(err)
	}
	var response string

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return "", err
	}

	ch, err := conn.Channel()
	if err != nil {
		return "", err
	}

	defer ch.Close()
	defer conn.Close()

	switch dataParsed.HttpType {
	case "POST":
		response, err = PostTest(dataParsed)
		if err != nil {
			panic(err)
		}
		ch.Publish(
			"",                  // default exchange
			dataParseds.ReplyTo, // la cola que mandó el request
			false,
			false,
			amqp.Publishing{
				ContentType:   "application/json",
				Body:          []byte(response),
				CorrelationId: dataParseds.CorrelationId,
			},
		)

	case "GET":
		response, err = GetTest(dataParsed)
		if err != nil {
			panic(err)
		}
		ch.Publish(
			"",                  // default exchange
			dataParseds.ReplyTo, // la cola que mandó el request
			false,
			false,
			amqp.Publishing{
				ContentType:   "application/json",
				Body:          []byte(response),
				CorrelationId: dataParseds.CorrelationId,
			},
		)
	}

	return response, nil
}

func PostTest(dataP DataParsed) (result string, err error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(dataP.Request).
		Post(dataP.Url)

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func GetTest(dataP DataParsed) (result string, err error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Get(dataP.Url)

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}
