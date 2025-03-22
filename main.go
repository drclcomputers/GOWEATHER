package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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
			Description string `json:"description"`
		} `json:"weather"`
	} `json:"list"`
	City struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	}
}

const (
	C0                   = 273.15 //0°C = 273K
	apifile       string = "config.txt"
	cache_file    string = "cache.txt"
	realtime_link string = "http://api.openweathermap.org/data/2.5/weather?&appid=%s&q=%s"
	forecast_link string = "https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s"
)

var (
	metric       bool
	showforecast bool
	cityinp      string
	api          string
)

// convert K to C or to F
func conv_K(degrees float64) string {
	if metric {
		return fmt.Sprintf("%.1f°C", degrees-C0)
	}
	return fmt.Sprintf("%.1f°F", (degrees-C0)*1.8+32)
}

// Get wind direction
func directwind(deg float64) string {
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((deg+22.5)/45.0) % 8
	return directions[index]
}

func pause() {
	fmt.Println("Press any key to continue...")
	fmt.Scanln()
}

// get api
func api_key_check() string {
	apikey, err := os.ReadFile(apifile)
	apikeystr := strings.TrimSpace(string(apikey))
	if err != nil {
		log.Fatalf("Error reading API key: %v", err)
		os.Exit(2)
	}
	if strings.Contains(apikeystr, "INSERT") || apikeystr == "" {
		fmt.Println("Unable to find API Key. Go to https://home.openweathermap.org/, create an account and get an API key. Afterwards, paste it into 'config.txt'.")
		os.Exit(2)
	}
	return apikeystr
}

// input city
func read_city() string {
	fmt.Print("Input city to view forecast: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// parameters
func parameters() {
	flag.BoolVar(&metric, "i", false, "Use imperial units")
	flag.BoolVar(&showforecast, "f", false, "Show forecast instead of current weather")
	flag.StringVar(&cityinp, "c", "", "The name of the city (exp: -c \"los angeles\")")
	flag.Parse()
	metric = !metric
}

// current weather
func current(city string) string {
	var message string

	fullurl := fmt.Sprintf(realtime_link, api, city)

	resp, err := http.Get(fullurl)
	if err != nil {
		fmt.Println("Error fetching data! Can't connect to network or the city doesn't exist!")
		os.Exit(2)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error fetching data!")
		os.Exit(2)
	}

	var data Weather

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error parsing JSON response.")
		fmt.Println("Response body:", string(body))
		os.Exit(2)
	}

	if data.Cod == 404 {
		fmt.Println("City not found!")
		os.Exit(2)
	} else if data.Cod != 200 {
		fmt.Printf("Error code: %d", data.Cod)
		os.Exit(2)
	}

	message += "Location: " + data.Name + ", " + data.Sys.Country + "\n"
	if len(data.Weatherdata) > 0 {
		message += "Weather: " + data.Weatherdata[0].Description + "\n"
	} else {
		message += "Weather description unavailable!\n"
	}

	message += "Temperature: " + conv_K(data.Main.Temp) + "\n"
	message += "Feels like: " + conv_K(data.Main.FeelsLike) + "\n"

	message += "Max temperature: " + conv_K(data.Main.MaxTemp) + "\n"
	message += "Min temperature: " + conv_K(data.Main.MinTemp) + "\n"

	message += "Pressure: " + fmt.Sprintf("%d", data.Main.Pressure) + "hPa -> "
	switch {
	case data.Main.Pressure < 980:
		message += "Low pressure -> Cloudy/Raining!"
	case data.Main.Pressure > 1000:
		message += "High pressure -> Sunny!"
	default:
		message += "Normal pressure"
	}
	message += "\n"
	message += "Humidity: " + fmt.Sprintf("%d", data.Main.Humidity) + "%\n"
	if metric {
		message += "Wind: " + fmt.Sprintf("%.1f", data.Wind.Speed*3.6) + "km/h -> From: "
	} else {
		message += "Wind: " + fmt.Sprintf("%.1f", data.Wind.Speed*2.23694) + "mph -> From: "
	}
	message += directwind(data.Wind.Deg) + "\n"
	if metric {
		message += fmt.Sprintf("Visibility: %d km\n", data.Visibility/1000)
	} else {
		message += fmt.Sprintf("Visibility: %.1f miles\n", float64(data.Visibility)*0.0006213712)
	}
	message += "Sunrise: " + time.Unix(int64(data.Sys.Sunrise), 0).Format("15:04") + " GMT+2\n"
	message += "Sunset: " + time.Unix(int64(data.Sys.Sunset), 0).Format("15:04") + " GMT+2\n"
	return message
}

// forecast
func show_forecast(city string) string {
	fullurl := fmt.Sprintf(forecast_link, city, api)
	resp, err := http.Get(fullurl)
	if err != nil {
		fmt.Println("Error fetching data! Can't connect to network or the city doesn't exist!")
		os.Exit(2)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error fetching data!")
		os.Exit(2)
	}

	var message string
	var data WeatherForecast

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error parsing JSON response.")
		fmt.Println("Response body:", string(body))
		os.Exit(2)
	}

	if data.Cod == "404" {
		fmt.Println("City not found!")
		os.Exit(2)
	} else if data.Cod != "200" {
		fmt.Printf("Error code: %s", data.Cod)
		os.Exit(2)
	}

	message += "Location: " + data.City.Name + ", " + data.City.Country + "\n"
	for _, forecast := range data.List {
		t := time.Unix(forecast.Dt, 0).UTC()
		if t.Hour() == 12 { // Check if the time is 12:00 PM
			message += fmt.Sprintf("Date: %s, Temp: %s, Weather: %s\n", t.Format("2006-01-02"), conv_K(forecast.Main.Temp), forecast.Weather[0].Description)
		}
	}
	return message
}

// cache file
func cache(message string) {
	currentTime := time.Now().String()
	if err := os.WriteFile(cache_file, []byte(currentTime+"\n\n"+message), 0644); err != nil {
		log.Printf("Failed to write cache file: %v", err)
	}
}

func main() {
	fmt.Println("GOWEATHER - ver 0.3.0")
	api = api_key_check()
	parameters()

	for cityinp == "" {
		cityinp = read_city()
	}

	urlcity := url.QueryEscape(cityinp)

	var message string

	if !showforecast {
		message = current(urlcity)
	} else {
		message = show_forecast(urlcity)
	}

	fmt.Println(message)
	cache(message)

	pause()

}
