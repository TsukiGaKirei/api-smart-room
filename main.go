package main

import (
	"api-smart-room/static"
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

	"google.golang.org/api/option"
	"gorm.io/gorm"

	"api-smart-room/database"
	"api-smart-room/routes"

	"cloud.google.com/go/pubsub"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Constants for Pub/Sub
const (
	projectID      = "smart-room-final-project"
	topicID        = "microcontroller-topic"
	subscriptionID = "api-sub"
)

func main() {
	// Initialize the database
	database.Init()

	// Set up Echo instance for API
	e := routes.Init()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Custom error message handling
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Set up Pub/Sub subscriber
	ctx := context.Background()
	client, subscription, err := initSubscriber(ctx, projectID, subscriptionID, "./creds.json", handlePubSubMessage, e)
	if err != nil {
		log.Fatalf("Error setting up Pub/Sub subscriber: %v", err)
	}

	// Start Pub/Sub subscriber in a separate goroutine
	go func() {
		err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			handlePubSubMessage(ctx, msg)
			msg.Ack()
		})
		if err != nil {
			log.Fatalf("Error receiving Pub/Sub messages: %v", err)
		}
	}()

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

	// Close Pub/Sub client when shutting down
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

// Custom error handling for HTTP errors
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

// Initialize Pub/Sub subscriber
func initSubscriber(ctx context.Context, projectID, subscriptionID, credsPath string, handlerFunc func(context.Context, *pubsub.Message), e *echo.Echo) (*pubsub.Client, *pubsub.Subscription, error) {
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("./creds.json"))
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating Pub/Sub client: %v", err)
	}

	subscription := client.Subscription(subscriptionID)
	subscription.ReceiveSettings.MaxOutstandingMessages = 10
	subscription.ReceiveSettings.NumGoroutines = 10

	return client, subscription, nil
}

// Handle Pub/Sub messages// Handle Pub/Sub messages
func handlePubSubMessage(ctx context.Context, msg *pubsub.Message) {
	// Decode the message data
	data := string(msg.Data)

	// Split the message into parts using a comma as the delimiter
	parts := strings.Split(data, ",")

	// Ensure that there are three parts in the message
	if len(parts) != 3 {
		log.Printf("Invalid message format: %s\n", data)
		return
	}

	// Parse each part into the respective variables
	RID, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Printf("Error parsing RID: %v\n", err)
		return
	}

	RoomTemp, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		log.Printf("Error parsing RoomTemp: %v\n", err)
		return
	}

	CountPerson, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("Error parsing CountPerson: %v\n", err)
		return
	}

	// Now you have RID, RoomTemp, and CountPerson as variables
	log.Printf("RID: %d, RoomTemp: %f, CountPerson: %d\n", RID, RoomTemp, CountPerson)
	db := database.GetDBInstance()
	if err = db.Raw(`update rooms set `, RID, RoomTemp, CountPerson).Error; err != nil {
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
