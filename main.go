package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const K0 = 273.15

func conv_K_C(degrees float64) float64 { //convert K to C
	d := float64(degrees)
	return d - K0
}

func directwind(deg float64) string {
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((deg+22.5)/45.0) % 8
	return directions[index]
}

func main() {
	fmt.Println("GOWEATHER - ver 0.9.4")
	apikey := "INSERT YOUR OPENWEATHERMAP API KEY"
	fmt.Print("Input city to view forecast: ")

	var city string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	city = scanner.Text()
	urlcity := url.QueryEscape(city)

	fullurl := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?&appid=%s&q=%s", apikey, urlcity)

	resp, err := http.Get(fullurl)
	if err != nil {
		fmt.Println("Error fetching data! Can't connect to network or the city doesn't exist!")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error fetching data!")
		return
	}

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	if sys, exist := data["sys"].(map[string]interface{}); exist {
		if country, ok := sys["country"].(string); ok {
			fmt.Printf("Location: %s, %s\n", city, country)
		}
	}

	if weather, exist := data["weather"].([]interface{}); exist && len(weather) > 0 {
		if description, ok := weather[0].(map[string]interface{})["description"].(string); ok {
			fmt.Printf("Weather: %s\n", description)
		}
	}

	if temp, exist := data["main"].(map[string]interface{})["temp"].(float64); exist {
		fmt.Printf("Temperature: %.2f°C\n", conv_K_C(temp))
	} else {
		fmt.Println("City doesn't exist!")
		return
	}

	if feel, exist := data["main"].(map[string]interface{})["feels_like"].(float64); exist {
		fmt.Printf("Feels like: %.2f°C\n", conv_K_C(feel))
	}

	if pres, exist := data["main"].(map[string]interface{})["pressure"].(int); exist {
		fmt.Printf("Pressure: %dhPa -> ", pres)
		if pres < 980 {
			fmt.Printf("Low pressure! Cloudy/Raining!")
		} else if pres > 1000 {
			fmt.Printf("High pressure! Sunny!")
		} else {
			fmt.Printf("Normal pressure!")
		}
		fmt.Printf("\n")
	}

	if humid, exist := data["main"].(map[string]interface{})["humidity"].(float64); exist {
		fmt.Printf("Humidity: %.0f%%\n", humid)
	}

	if wind, exist := data["wind"].(map[string]interface{})["speed"].(float64); exist {
		fmt.Printf("Wind: %.2fkm/h -> From: ", wind*3.6)
	}

	if wind, exist := data["wind"].(map[string]interface{}); exist {
		if windDeg, ok := wind["deg"].(float64); ok {
			fmt.Printf("%s\n", directwind(windDeg))
		}
	}

	if sys, exist := data["sys"].(map[string]interface{}); exist {
		if sunrise, ok := sys["sunrise"].(float64); ok {
			fmt.Printf("Sunrise: %s\n", time.Unix(int64(sunrise), 0).Format("15:04"))
		}
		if sunset, ok := sys["sunset"].(float64); ok {
			fmt.Printf("Sunset: %s\n", time.Unix(int64(sunset), 0).Format("15:04"))
		}
	}

}
