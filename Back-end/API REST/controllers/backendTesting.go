package controllers

import (
	"Go-API-T/middlewere"
	"context"

	//"strings"
	//"time"
	"encoding/json"
	"log"
	"time"

	//"example/Go-API-T/services"
	//"encoding/base64"
	//"encoding/json"
	//"context"
	//"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	//"github.com/golang-jwt/jwt/v5"
	//"golang.org/x/crypto/bcrypt"
	//"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",      // type
		true,         // durable
		false,        // auto deleted ?
		false,        // internal ?
		false,        // no wait ?
		nil,          // arguments ?
	)
}

// function to save all endpoints in a "router".
func TestsRoutes(rg *gin.RouterGroup, handler *HandlerAPI, mw *middlewere.Middleware) {
	user := rg.Group("/tests") //prefix that all routes(endpoints)

	user.POST("/test-event", handler.TestEvent)

	//user.POST("/TwoStep", twoStep)
}

// define a

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func (e *Emitter) Push(event string, severity string, dataJson []byte) error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	payload := Payload{
		Name: event,
		Data: string(dataJson),
	}

	body, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"logs_topic", // exchange
		severity,     // routing key
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: uuid.New().String(), // identificador
			ReplyTo:       "response_queue",
		},
	)
	if err != nil {
		return err
	}

	log.Println("Sent message:", string(body))
	return nil
}
type PayloadsR struct {
    Name          string `json:"name"`
    Data          string `json:"data"`
    CorrelationId string `json:"correlationId"`
    ReplyTo       string `json:"replyTo"`
}
func (h *HandlerAPI) TestEvent(c *gin.Context) {
	// Connect rabbitMQ
	var jsonData struct {
		Id_Group    int
		HttpType    string
		Url         string
		RequestType string
		Request     map[string]interface{}
		Header      map[string]interface{}
		Token       string
	}
	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to RabbitMQ"})
		return
	}
	defer conn.Close()
	//create channel
	ch, err := conn.Channel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open channel"})
		return
	}
	defer ch.Close()
	declareExchange(ch)
	//declare response queue
	msgs, err := ch.QueueDeclare(
		"response_queue", // request
		true,             // autoAck
		false,            // exclusive
		false,            // noLocal
		false,            // noWait
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to consume response"})
		return
	}

	/* create emitter (replace for msgs, consume channel, because is a channel of response and request)
	emitter, err := NewEmitter(conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	*/

	reqJson, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal request"})
		return
	}

	corrID := uuid.New().String()
	// Publish message with replyto and correlationid

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	payload := PayloadsR{
		Name:          "test",
		Data:          string(reqJson), 
		ReplyTo:       msgs.Name,       // cola de respuesta
		CorrelationId: corrID,          // id único
	}
	body, _ := json.Marshal(payload)

	err = ch.PublishWithContext(
		ctx,
		"logs_topic", // exchange
		"test",       // routing key
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish message"})
		return
	}

	//consume the response queue
	res, err := ch.Consume(
		msgs.Name,
		"",
		true,  //autoAck
		false, //exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	responseCH := make(chan string)
	//sincronize channel
	go func() {
		for d := range res {
			if d.CorrelationId == corrID {
				responseCH <- string(d.Body)
				break
			}
		}
	}()
	//wait timeout or response
	var jsonMap map[string]interface{} //for parse to json
	select {
	case res := <-responseCH:
		json.Unmarshal([]byte(res), &jsonMap) //parse to json
		c.JSON(http.StatusOK, gin.H{"result": jsonMap}) //response a json
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "timeout waiting for response"})
	}

}
