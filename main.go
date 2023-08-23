package main

import (
	"api-smart-room/api/auth"
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"

	"api-smart-room/database"
	"api-smart-room/routes"
	"api-smart-room/static"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

//keep running to receive message from microcontroller
const (
	projectID      = "delta-coil-393803"
	topicID        = "api-subscribe"
	subscriptionID = "api-subscribe-sub"
)

func main() {
	database.Init()
	ctx := context.Background()

	// Initialize a Google Cloud Pub/Sub client with credentials
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("./delta-coil-393803-555a240284e0.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	auth.ReceiveMessage(ctx, client, subscriptionID)
	//
	e := routes.Init()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Custom error message
	e.HTTPErrorHandler = func(err error, c echo.Context) {
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
			report = echo.NewHTTPError(http.StatusUnauthorized, "Internal Server Error")
		} else if errors.Is(err, echo.ErrNotFound) {
			report = echo.NewHTTPError(http.StatusUnauthorized, "Halaman tidak ditemukan")
		}

		c.Logger().Error(report)
		errObj := static.ResponseError{
			Error:   true,
			Message: fmt.Sprintf("%v", report.Message),
		}
		c.JSON(report.Code, errObj)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

// Validation
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
