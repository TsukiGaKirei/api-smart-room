package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"api-smart-room/database"
	"api-smart-room/routes"
	"api-smart-room/static"
)

// Constants for MQTT
const (
	mqttBrokerHost     = "34.101.160.103"
	mqttBrokerPort     = 1883
	mqttBrokerUsername = "user1"
	mqttBrokerPassword = "qweasd123"
	mqttTopic          = "esp32-topic"
)

func main() {
	// Initialize the database
	database.Init()

	// Set up Echo instance for API
	e := routes.Init()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Custom error message handling
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Initialize the MQTT client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBrokerHost, mqttBrokerPort))
	opts.SetClientID("your-client-id") // Set a unique client ID

	opts.SetUsername(mqttBrokerUsername)
	opts.SetPassword(mqttBrokerPassword)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}

	// Subscribe to the MQTT topic
	if token := client.Subscribe(mqttTopic, 0, handleMQTTMessage); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to MQTT topic: %v", token.Error())
	}

	// Start HTTP server
	port := getPort()
	go func() {
		if err := e.Start(":" + port); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Graceful shutdown on interrupt or terminate signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shut down the HTTP server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// Unsubscribe and disconnect from MQTT broker
	if token := client.Unsubscribe(mqttTopic); token.Wait() && token.Error() != nil {
		log.Fatal("Error unsubscribing from MQTT topic:", token.Error())
	}

	client.Disconnect(0)
}

// Custom error handling for HTTP errors
func customHTTPErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				report.Message = fmt.Sprintf("Mohon isi %s", err.Field())
			case "email":
				report.Message = fmt.Sprintf("%s bukanlah email yang valid", err.Field())
			case "gte":
				report.Message = fmt.Sprintf("%s harus lebih besar dari %s", err.Field(), err.Param())
			case "lte":
				report.Message = fmt.Sprintf("%s harus lebih kurang dari %s", err.Field(), err.Param())
			}

			break
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		report = echo.NewHTTPError(http.StatusNotFound, "Data tidak ditemukan")
	} else if errors.Is(err, &static.AuthError{}) {
		report = echo.NewHTTPError(http.StatusUnauthorized, "User tidak ditemukan")
	} else if errors.Is(err, echo.ErrInternalServerError) {
		report = echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	} else if errors.Is(err, echo.ErrNotFound) {
		report = echo.NewHTTPError(http.StatusNotFound, "Halaman tidak ditemukan")
	}

	c.Logger().Error(report)
	errObj := static.ResponseError{
		Error:   true,
		Message: fmt.Sprintf("%v", report.Message),
	}
	c.JSON(report.Code, errObj)
}

// Validation
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Handle MQTT messages
func handleMQTTMessage(client mqtt.Client, msg mqtt.Message) {
	// Decode the message data
	data := string(msg.Payload())

	// Split the message into parts using a comma as the delimiter
	parts := strings.Split(data, ";")

	// Ensure that there are three parts in the message
	if len(parts) != 2 {
		log.Printf("Invalid message format: %s\n", data)
		return
	}

	// Parse each part into the respective variables
	RID, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Printf("Error parsing RID: %v\n", err)
		return
	}

	CountPerson, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Error parsing CountPerson: %v\n", err)
		return
	}

	// Now you have RID, RoomTemp, and CountPerson as variables
	log.Printf("RID: %d, CountPerson: %d\n", RID, CountPerson)
	db := database.GetDBInstance()
	if err = db.Exec(`update rooms set count_person=? where rid=?`, CountPerson, RID).Error; err != nil {
		log.Printf("Error post to database: %v\n", err)
		return
	}
}

// Get the port from the environment variable or use the default
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	return port
}
