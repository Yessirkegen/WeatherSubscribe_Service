package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WeatherData struct {
	Temperature float64
	Summary     string
}

type WeatherAPIClient interface {
	GetWeather(city string) (WeatherData, error)
}

type WeatherServiceClient struct {
	baseURL string
}

func NewWeatherServiceClient(baseURL string) *WeatherServiceClient {
	return &WeatherServiceClient{baseURL: baseURL}
}

func (c *WeatherServiceClient) GetWeather(city string) (WeatherData, error) {
	url := fmt.Sprintf("%s/weather?city=%s", c.baseURL, city)
	resp, err := http.Get(url)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherData{}, fmt.Errorf("failed to get weather data: %s", resp.Status)
	}

	var data WeatherData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return WeatherData{}, err
	}

	return data, nil
}
