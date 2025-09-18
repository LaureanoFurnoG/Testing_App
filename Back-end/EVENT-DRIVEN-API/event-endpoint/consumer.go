package event

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	amqp "github.com/rabbitmq/amqp091-go"
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
	Urlapi         string
	RequestType string
	Request     map[string]interface{}
	Header      map[string]interface{}
	Token       string
}
type ResponseData struct {
	Response         string
	HttpResponseCode string
}

func handlePayload(payload Payload) {
	//fmt.Println(payload)
	//log.Print(payload)
	switch payload.Name {
	case "test", "event":
		// log whatever we get
		_, err := testApp(payload)
		if err != nil {
			log.Panic(err)
		}
		//log.Println(result)
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
	case "PUT":
		response, err = PutTest(dataParsed)
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
	case "PATCH":
		response, err = PatchTest(dataParsed)
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
	case "DELETE":
		
		response, err = DeleteTest(dataParsed)
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

func PutTest(dataP DataParsed) (string, error) {
	client := resty.New()

	headers, err := HeadersSet(dataP)
	if err != nil {
		return "", err
	}

	//log.Println(headers)
	resp, err := client.R().
		SetBody(dataP.Request).
		SetHeaders(headers).
		Put(dataP.Urlapi)
	if err != nil {
		return "", err
	}
	responseJSON, _ := JsonResponse(resp)
	return responseJSON, nil
}

func PatchTest(dataP DataParsed) (string, error) {
	client := resty.New()

	headers, err := HeadersSet(dataP)
	if err != nil {
		return "", err
	}

	//log.Println(headers)
	resp, err := client.R().
		SetBody(dataP.Request).
		SetHeaders(headers).
		Patch(dataP.Urlapi)
	if err != nil {
		return "", err
	}
	responseJSON, _ := JsonResponse(resp)
	return responseJSON, nil
}

func PostTest(dataP DataParsed) (result string, err error) {
	client := resty.New()
	if len(dataP.Header) == 0 {
		resp, err := client.R().
			SetBody(dataP.Request).
			Post(dataP.Urlapi)

		if err != nil {
			return "", err
		}
		responseJSON, _ := JsonResponse(resp)
		return responseJSON, nil

	} else {
		headers, err := HeadersSet(dataP)
		if err != nil {
			return "", err
		}
		//log.Println(headers)
		resp, err := client.R().
			SetBody(dataP.Request).
			SetHeaders(headers).
			Post(dataP.Urlapi)
		if err != nil {
			return "", err
		}
		responseJSON, _ := JsonResponse(resp)
		return responseJSON, nil

	}

}

func DeleteTest(dataP DataParsed) (result string, err error) {
	client := resty.New()
	headers, err := HeadersSet(dataP)
	if err != nil {
		return "", err
	}
	//log.Println(headers)
	resp, err := client.R().
		SetBody(dataP.Request).
		SetHeaders(headers).
		Delete(dataP.Urlapi)
	if err != nil {
		return "", err
	}
	responseJSON, _ := JsonResponse(resp)
	return responseJSON, nil
}

func GetTest(dataP DataParsed) (result string, err error) {
	client := resty.New()

	params := make(map[string]string)
	for k, v := range dataP.Request {
		// if the value is string
		if str, ok := v.(string); ok {
			params[k] = str
		} else {
			// convert to string if the value is an interface or another type
			params[k] = fmt.Sprintf("%v", v)
		}
	}

	headers, err := HeadersSet(dataP)
	if err != nil {
		return "", err
	}

	resp, err := client.R().
		SetQueryParams(params).
		SetHeaders(headers).
		Get(dataP.Urlapi)

	if err != nil {
		return "", err
	}

	responseJSON, _ := JsonResponse(resp)
	return responseJSON, nil
}

func HeadersSet(dataP DataParsed) (result map[string]string, err error) {
	headers := make(map[string]string)
	for k, v := range dataP.Header {
		// if the value is string
		if str, ok := v.(string); ok {
			if k == "Access-Token" {
				headers["Authorization"] = "Bearer " + str
			} else {
				headers[k] = str
			}
		} else {
			// convert to string if the value is an interface or another type
			headers[k] = fmt.Sprintf("%v", v)
		}
	}
	return headers, nil
}

func JsonResponse(resp *resty.Response) (result string, err error) {
	ResponseDatas := ResponseData{
		Response:         resp.String(),
		HttpResponseCode: strconv.Itoa(resp.StatusCode()),
	}

	responseJSON, err := json.Marshal(ResponseDatas)
	if err != nil {
		return "", err
	}

	return string(responseJSON), nil
}
