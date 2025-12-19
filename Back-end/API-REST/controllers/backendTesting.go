package controllers

import (
	"Go-API-T/middlewere"
	"Go-API-T/services"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"Go-API-T/initializers"
	"Go-API-T/models"
	"encoding/json"
	"log"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
	test := rg.Group("/tests/:groupId") //prefix that all routes(endpoints)

	test.POST("/test-event", mw.RequireAuth(), mw.BelongsGroup(), handler.TestEvent)
	test.GET("/find-tests", mw.RequireAuth(), mw.BelongsGroup(), handler.FindTest)
	test.GET("/find-testsSave", mw.RequireAuth(), mw.BelongsGroup(), handler.FindTestSave)

	test.POST("/test-all", mw.RequireAuth(), mw.BelongsGroup(), handler.RunAllEnd)
	test.POST("/create-save", mw.RequireAuth(), mw.BelongsGroup(), handler.saveTestEndpointsDoc)
	test.DELETE("/delete-testSave", mw.RequireAuth(), mw.BelongsGroup(), handler.DeleteEndpointSave)
	test.DELETE("/delete-test", mw.RequireAuth(), mw.BelongsGroup(), handler.DeleteEndpoint)
	test.POST("/test-front", mw.RequireAuth(), mw.BelongsGroup(), handler.TestFrontEnd)
	test.POST("/make-pdf", mw.RequireAuth(), mw.BelongsGroup(), handler.MakePDF)

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
type jsonMapResponse struct {
	HttpResponseCode string
	Response         string
}
type jsonData struct {
	HttpType    string
	Urlapi      string
	Name        string
	RequestType string
	Request     map[string]interface{}
	Header      map[string]interface{}
	Token       string

	Response         map[string]interface{}
	ResponseHttpCode int
}

func (h *HandlerAPI) TestEvent(c *gin.Context) {
	// Connect rabbitMQ
	groupId := c.Param("groupId")
	var jsonDataRe struct {
		Id_Group    int
		Name        string
		HttpType    string
		Urlapi      string
		RequestType string
		Request     map[string]interface{}
		Header      map[string]interface{}
		Token       string
	}
	convStringToInt, err := strconv.Atoi(groupId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse groupID",
		})
		return
	}
	jsonDataRe.Id_Group = convStringToInt
	if c.ShouldBindJSON(&jsonDataRe) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	accessHeader := c.GetHeader("Authorization")
	if accessHeader == "" || !strings.HasPrefix(accessHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Access token not found or invalid format"})
		c.Abort()
		return
	}

	accessToken := strings.TrimPrefix(accessHeader, "Bearer ")
	accessToken = strings.TrimSpace(accessToken)
	
	fmt.Println(accessToken)
	testDriven(jsonDataRe, accessToken, c, h)

}

func saveDataTest(Id_Group int, values jsonData, accessToken string, c *gin.Context, h *HandlerAPI) error {
	var group models.Groups
	searchGroup := initializers.DB.Find(&group, "id = ?", Id_Group)

	if searchGroup.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group Missing",
		})
		return searchGroup.Error
	}

	//search user
	var userF models.Users
	userKeycloak, err := h.clientKC.UserInfo(c.Request.Context(), accessToken)
	fmt.Println(accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return err
	}

	userFind := initializers.DB.First(&userF, "keycloak_id = ?", userKeycloak.ID)

	if userFind.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User Missing",
		})
		return err
	}
	//search if the test exist already

	//create instance of test data in local database
	requestJSON, err := json.Marshal(values.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to marshal request data",
		})
		return err
	}

	responseJSON, err := json.Marshal(values.Response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to marshal response data",
		})
		return err
	}

	headerJSON, err := json.Marshal(values.Header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to marshal header data",
		})
		return err
	}
	testS := models.Backendtests{
		Idgroup:          uint(Id_Group),
		Name:             values.Name,
		Httptype:         values.HttpType,
		Urlapi:           values.Urlapi,
		Requesttype:      values.RequestType,
		Request:          requestJSON,
		Response:         responseJSON,
		ResponseHttpCode: values.ResponseHttpCode,
		Header:           headerJSON,
		Token:            values.Token,
	}

	findTest := initializers.DB.Where("Urlapi = ? AND response_http_code = ?", testS.Urlapi, testS.ResponseHttpCode).Find(&testS)
	if findTest.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search data",
		})
		return err
	}

	if findTest.RowsAffected == 0 {
		createTest := initializers.DB.Create(&testS)
		if createTest.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Test not created in the database",
			})
			return err
		}
	} else {
		updateTest := initializers.DB.Model(&testS).Updates(models.Backendtests{
			Httptype:         values.HttpType,
			Name:             values.Name,
			Requesttype:      values.RequestType,
			Request:          requestJSON,
			Response:         responseJSON,
			Header:           headerJSON,
			Token:            values.Token,
			Urlapi:           values.Urlapi,
			ResponseHttpCode: values.ResponseHttpCode,
		})
		if updateTest.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test"})
			return updateTest.Error
		}
	}

	return nil
}

func (h *HandlerAPI) FindTest(c *gin.Context) {
	var jsonDataRe struct {
		Name string
	}

	if c.ShouldBindJSON(&jsonDataRe) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	testS := []models.Backendtests{}
	find := initializers.DB.Where("name LIKE ? AND Idgroup = ?", "%"+jsonDataRe.Name+"%", c.Param("groupId")).Find(&testS)

	if find.Error != nil {
		c.JSON(404, gin.H{
			"error": "Search endpoint error ",
		})
		return
	}
	c.JSON(200, gin.H{
		"Group": testS,
	})
}

type jsonDataRe struct {
	Id_Group    int
	Name        string
	HttpType    string
	Urlapi      string
	RequestType string
	Request     map[string]interface{}
	Header      map[string]interface{}
	Token       string
}

// body
type TestsRequest struct {
	Tests []jsonDataRe `json:"tests"`
}

func (h *HandlerAPI) RunAllEnd(c *gin.Context) {
	accessToken := c.Request.Header.Get("Access-Token")

	var req TestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalid"})
		return
	}

	for _, value := range req.Tests {
		b, _ := json.Marshal(value)

		var re jsonDataRe

		if err := json.Unmarshal(b, &re); err != nil {
			c.JSON(156, gin.H{"error": err.Error()})
			continue
		}

		testDriven(re, accessToken, c, h)
	}
}

func testDriven(jsonDataRe jsonDataRe, accessToken string, c *gin.Context, h *HandlerAPI) {
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
	reqJson, err := json.Marshal(jsonDataRe)
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
		ReplyTo:       msgs.Name, // cola de respuesta
		CorrelationId: corrID,    // id único
	}
	body, _ := json.Marshal(payload)

	err = ch.PublishWithContext(
		ctx,
		"logs_topic", // exchange
		"test",       // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
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
	var jsonMap jsonMapResponse //for parse to json
	var jsonResult map[string]interface{}
	select {
	case res := <-responseCH:
		json.Unmarshal([]byte(res), &jsonMap)
		json.Unmarshal([]byte(jsonMap.Response), &jsonResult) //parse to json
		// Parse HTTP response code
		httpCode, err := strconv.Atoi(jsonMap.HttpResponseCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Your endpoint, maybe not response with a status code",
			})
			return
		}
		saveData := jsonData{
			HttpType:    jsonDataRe.HttpType,
			Name:        jsonDataRe.Name,
			Urlapi:      jsonDataRe.Urlapi,
			RequestType: jsonDataRe.RequestType,
			Request:     jsonDataRe.Request,
			Header:      jsonDataRe.Header,
			Token:       jsonDataRe.Token,

			Response:         jsonResult,
			ResponseHttpCode: httpCode,
		}
		saveD := saveDataTest(jsonDataRe.Id_Group, saveData, accessToken, c, h)
		if saveD != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error in saved the tested",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": map[string]interface{}{"HTTP_Code": jsonMap.HttpResponseCode, "Response": jsonResult}}) //response a json
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "timeout waiting for response"})
	}
}

func (h *HandlerAPI) saveTestEndpointsDoc(c *gin.Context) {
	var request struct {
		Idgroup             int
		Idtest              int
		Testcasedescription string
		Testedinfrontend    bool
		Evidencefrontend    string
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalid"})
		return
	}

	//group
	saveendpointresult := models.Saveendpointresult{
		Idgroup:             uint(request.Idgroup),
		Idtest:              uint(request.Idtest),
		Testcasedescription: request.Testcasedescription,
		Testedinfrontend:    request.Testedinfrontend,
		Evidencefrontend:    request.Evidencefrontend,
	}
	findTest := initializers.DB.Where("Idtest = ?", request.Idtest).Find(&saveendpointresult)
	if findTest.RowsAffected == 0 {
		createSave := initializers.DB.Create(&saveendpointresult)
		if createSave.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The test couldn't create in the database"})
			return
		}
	} else {
		updateT := initializers.DB.Model(&saveendpointresult).Updates(&saveendpointresult)
		if updateT.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The test couldn't create in the database"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"Msg": "The test saved successfull"})
}

func (h *HandlerAPI) DeleteEndpointSave(c *gin.Context) {
	var requestJSON struct {
		Idtest int
	}
	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalid"})
		return
	}
	saveendpointresult := models.Saveendpointresult{}

	deleteTest := initializers.DB.Where("Idtest = ?", requestJSON.Idtest).Delete(&saveendpointresult)
	if deleteTest.Error != nil {
		c.JSON(404, gin.H{"error": "Test not found"})
		return
	}
	c.JSON(204, gin.H{"Msg": "The test deleted successfull"})

}

// this is for
func (h *HandlerAPI) DeleteEndpoint(c *gin.Context) {
	var requestJSON struct {
		Id uint
	}
	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalid"})
		return
	}
	Backendtests := models.Backendtests{}
	saveendpointresult := models.Saveendpointresult{}

	deleteTest := initializers.DB.Where("Id = ?", requestJSON.Id).Find(&Backendtests)
	initializers.DB.Where("Idtest = ?", &Backendtests.ID).Delete(&saveendpointresult)
	deleteTest.Delete(&Backendtests)
	if deleteTest.Error != nil {
		c.JSON(404, gin.H{"error": "Test not found"})
		return
	}

	c.JSON(204, gin.H{"Msg": "The test deleted successfull"})
}

type LighthouseResult struct {
	LighthouseResult struct {
		Categories struct {
			Performance struct {
				Score float64 `json:"score"`
			} `json:"performance"`
			Accessibility struct {
				Score float64 `json:"score"`
			} `json:"accessibility"`
			Seo struct {
				Score float64 `json:"score"`
			} `json:"seo"`
			BestPractices struct {
				Score float64 `json:"score"`
			} `json:"best-practices"`
		} `json:"categories"`
		Audits map[string]struct {
			ID          string      `json:"id"`
			Title       string      `json:"title"`
			Description string      `json:"description"`
			Score       interface{} `json:"score"`
			Details     interface{} `json:"details"`
		} `json:"audits"`
	} `json:"lighthouseResult"`
}

func (h *HandlerAPI) TestFrontEnd(c *gin.Context) {
	var requestJSON struct {
		Url      string
		Strategy string
	}

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalid"})
		return
	}

	apiURL := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=%s&key=%s&strategy=%s&category=performance&category=accessibility&category=seo&category=best-practices", requestJSON.Url, os.Getenv("API_KEY"), requestJSON.Strategy)

	resp, err := http.Get(apiURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result LighthouseResult
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"result": &result})
}

type Section struct {
	Title string
	Page  int
	Level int
}

func (h *HandlerAPI) FindTestSave(c *gin.Context) {
	var jsonDataRe struct {
		Id int
	}
	if c.ShouldBindJSON(&jsonDataRe) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	testSave := []models.Saveendpointresult{}
	find := initializers.DB.Preload("Backendtests").Preload("Group").Find(&testSave)

	if find.Error != nil {
		c.JSON(404, gin.H{
			"error": "Search endpoint saved error ",
		})
		return
	}

	var Result struct {
		Response []models.Saveendpointresult
	}

	Result.Response = testSave

	c.JSON(200, gin.H{"Result": Result})

}

func (h *HandlerAPI) MakePDF(c *gin.Context) {
	var jsonData services.Request

	if c.ShouldBindJSON(&jsonData) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	testsSaved := []models.Saveendpointresult{}
	all_tests_saved := initializers.DB.Preload("Backendtests").Preload("Group").Find(&testsSaved)
	if all_tests_saved.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to return the endpoints save",
		})
		return
	}
	services.CreatePDF(jsonData, testsSaved, c)

	c.JSON(http.StatusOK, gin.H{"msg": "PDF generado"})
}
