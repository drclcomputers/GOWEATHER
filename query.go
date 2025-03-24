package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Weather struct {
	Name string `json:"name"`
	Sys  struct {
		Country string  `json:"country"`
		Sunrise float64 `json:"sunrise"`
		Sunset  float64 `json:"sunset"`
	} `json:"sys"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		MaxTemp   float64 `json:"temp_max"`
		MinTemp   float64 `json:"temp_min"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weatherdata []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Cod        int `json:"cod"`
	Visibility int `json:"visibility"`
}

type WeatherForecast struct {
	Cod  string `json:"cod"`
	List []struct {
		Dt   int64 `json:"dt"`
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		} `json:"weather"`
	} `json:"list"`
	City struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	}
}

type Location struct {
	City string `json:"city"`
}

const (
	CelsiusToKelvin = 273.15
	RealtimeURL     = "http://api.openweathermap.org/data/2.5/weather?&appid=%s&q=%s"
	ForecastURL     = "https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s"
	idapiURL        = "http://ip-api.com/json/"
)

var use_metric = true
var ApiKey = "72a8cad5f37d68dbf24ad918aca7ef41"

func getIpLocation() string {
	resp, err := http.Get(idapiURL)
	if err != nil {
		return fmt.Sprintf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("error reading response body: %v", err)
	}

	var locData Location
	if err := json.Unmarshal(body, &locData); err != nil {
		return fmt.Sprintf("error unmarshalling response: %v", err)
	}
	return locData.City
}

func convertKelvinToTemperature(kelvin float64) string {
	if use_metric {
		return fmt.Sprintf("%.0f°C", kelvin-CelsiusToKelvin)
	}
	return fmt.Sprintf("%.0f°F", (kelvin-CelsiusToKelvin)*1.8+32)
}

func getWindDirection(degree float64) string {
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((degree+22.5)/45.0) % 8
	return directions[index]
}

func getCurrentWeather(city string) (string, error) {
	url := fmt.Sprintf(RealtimeURL, ApiKey, city)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var weatherData Weather
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	if weatherData.Cod != 200 {
		return "", fmt.Errorf("error fetching weather data, code: %d", weatherData.Cod)
	}

	message := fmt.Sprintf("Location: %s, %s\n", weatherData.Name, weatherData.Sys.Country)
	message += fmt.Sprintf("Weather: %s\n", weatherData.Weatherdata[0].Main)
	message += fmt.Sprintf("Temperature: %s\n", convertKelvinToTemperature(weatherData.Main.Temp))
	message += fmt.Sprintf("Feels like: %s\n", convertKelvinToTemperature(weatherData.Main.FeelsLike))
	message += fmt.Sprintf("Max temperature: %s\n", convertKelvinToTemperature(weatherData.Main.MaxTemp))
	message += fmt.Sprintf("Min temperature: %s\n", convertKelvinToTemperature(weatherData.Main.MinTemp))
	message += fmt.Sprintf("Pressure: %d hPa -> %s\n", weatherData.Main.Pressure, getPressureDescription(weatherData.Main.Pressure))
	message += fmt.Sprintf("Humidity: %d%%\n", weatherData.Main.Humidity)
	message += fmt.Sprintf("Wind: %.1f %s -> From: %s\n", convertWindSpeed(weatherData.Wind.Speed), getWindSpeedUnit(), getWindDirection(weatherData.Wind.Deg))
	message += fmt.Sprintf("Visibility: %.1f %s\n", convertVisibility(weatherData.Visibility), getVisibilityUnit())
	message += fmt.Sprintf("Sunrise: %s GMT+2\n", time.Unix(int64(weatherData.Sys.Sunrise), 0).Format("15:04"))
	message += fmt.Sprintf("Sunset: %s GMT+2\n", time.Unix(int64(weatherData.Sys.Sunset), 0).Format("15:04"))

	return message, nil
}

func getPressureDescription(pressure int) string {
	switch {
	case pressure < 980:
		return "Low pressure -> Cloudy/Raining!"
	case pressure > 1000:
		return "High pressure -> Sunny!"
	default:
		return "Normal pressure"
	}
}

func convertWindSpeed(speed float64) float64 {
	if use_metric {
		return speed * 3.6
	}
	return speed * 2.23694
}

func getWindSpeedUnit() string {
	if use_metric {
		return "km/h"
	}
	return "mph"
}

func convertVisibility(visibility int) float64 {
	if use_metric {
		return float64(visibility) / 1000
	}
	return float64(visibility) * 0.0006213712
}

func getVisibilityUnit() string {
	if use_metric {
		return "km"
	}
	return "miles"
}

func getWeatherForecast(city string) (string, error) {
	url := fmt.Sprintf(ForecastURL, city, ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching forecast data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var forecastData WeatherForecast
	if err := json.Unmarshal(body, &forecastData); err != nil {
		return "", fmt.Errorf("error unmarshalling forecast response: %v", err)
	}

	if forecastData.Cod != "200" {
		return "", fmt.Errorf("error fetching forecast data, code: %s", forecastData.Cod)
	}

	message := "5 day forecast:\n"

	for _, forecast := range forecastData.List {
		t := time.Unix(forecast.Dt, 0).UTC()
		if t.Hour() == 12 { // 12:00 PM forecast
			message += fmt.Sprintf("Date: %s, Temp: %s, Weather: %s\n", t.Format("2006-01-02"), convertKelvinToTemperature(forecast.Main.Temp), forecast.Weather[0].Main)
		}
	}
	return message, nil
}

func querycity(cityInput string) string {
	escapedCity := url.QueryEscape(cityInput)
	currentWeather, err := getCurrentWeather(escapedCity)
	if err != nil {
		return ""
	}

	forecast, err := getWeatherForecast(escapedCity)
	if err != nil {
		return ""
	}

	return currentWeather + forecast
}
