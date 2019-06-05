package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// NB: this requires setting an OWM_KEY environment variable with your API key (see https://openweathermap.org/api)
const owm = "https://api.openweathermap.org/data/2.5/weather?appid=%s&units=imperial&q=%s"

// Response is of type APIGatewayProxyResponse since we're leveraging the AWS Lambda Proxy Request functionality (default behavior)
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest since we're leveraging the AWS Lambda Proxy Request functionality (default behavior)
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Request events.APIGatewayProxyRequest

// Coord JSON weather response
type Coord struct {
	Longitude float32 `json:"lon"`
	Latitude  float32 `json:"lat"`
}

// Weather JSON weather response
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Main JSON weather response
type Main struct {
	Temp     float32 `json:"temp"`
	Pressure int     `json:"pressure"`
	Humidity int     `json:"humidity"`
	TempMin  float32 `json:"temp_min"`
	TempMax  float32 `json:"temp_max"`
}

// Wind JSON weather response
type Wind struct {
	Speed   float32 `json:"speed"`
	Degrees int     `json:"deg"`
}

// Cloud JSON weather response
type Cloud struct {
	All int `json:"all"`
}

// WeatherResponse JSON
type WeatherResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Visibility  int       `json:"visibility"`
	Coordinates Coord     `json:"coord"`
	Weather     []Weather `json:"weather"`
	Main        Main      `json:"main"`
	Wind        Wind      `json:"wind"`
	Clouds      Cloud     `json:"clouds"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, req Request) (Response, error) {
	fmt.Printf("id: %s, path: %s\n", req.RequestContext.RequestID, req.Path)

	//construct URL to API
	if req.PathParameters == nil || len(req.PathParameters["id"]) == 0 {
		return Response{StatusCode: 400}, errors.New("no location provided")
	}
	api := fmt.Sprintf(owm, os.Getenv("OWM_KEY"), req.PathParameters["id"])
	fmt.Printf("api: %s\n", api)

	//make request and read data
	response, err := http.Get(api)
	if err != nil {
		return failure(err)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return failure(err)
	}

	//useful for debugging
	// fmt.Println(string(data))

	//convert data to our struct
	var wr WeatherResponse
	err = json.Unmarshal(data, &wr)
	if err != nil {
		return failure(err)
	}

	//build our response body
	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{
		"message": fmt.Sprintf("Weather for %s is %.1f° F with %s, %.1f MPH winds and %d%% humidity", wr.Name, wr.Main.Temp, wr.Weather[0].Description, wr.Wind.Speed, wr.Main.Humidity),
	})
	if err != nil {
		return failure(err)
	}
	json.HTMLEscape(&buf, body)

	//return our response
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}

//very basic error handling
func failure(err error) (Response, error) {
	fmt.Printf("ERROR: %v\n", err)
	return Response{StatusCode: 500}, err
}

func main() {
	lambda.Start(Handler)
}
