package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type weather struct {
	Name  string `json:"name"`
	Coord struct {
		LONG float64 `json:"lon"`
		LAT  float64 `json:"lat"`
	}
	Sys struct {
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
}

type weatherforecast struct {
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

const C0 = 273.15 //0°C = 273K

var metric bool = true
var showforecast bool = false

func conv_K(degrees float64) string { //convert K to C or to F
	if metric {
		return fmt.Sprintf("%.0f°C", degrees-C0)
	}
	return fmt.Sprintf("%.0f°F", (degrees-C0)*1.8+32)
}

func directwind(deg float64) string { //Get wind direction
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((deg+22.5)/45.0) % 8
	return directions[index]
}

func main() {

	fmt.Println("GOWEATHER - ver 0.2.1")
	apikey, err := os.ReadFile("config.txt")
	if err != nil {
		fmt.Println("Can't read API Key! 'config.txt' does not exist!")
		return
	}
	apikeystr := strings.TrimSpace(string(apikey))

	var cityparts []string
	for _, arg := range os.Args[1:] {
		if arg == "-i" {
			metric = false
		} else if arg == "-f" {
			showforecast = true
		} else {
			cityparts = append(cityparts, arg)
		}
	}

	var city string
	if len(cityparts) == 0 {
		fmt.Print("Input city to view forecast: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		city = scanner.Text()
		for city == "" {
			fmt.Print("Input city to view forecast: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			city = scanner.Text()
		}
	} else {
		city = strings.Join(cityparts, " ")
	}

	urlcity := url.QueryEscape(city)

	if strings.Contains(apikeystr, "INSERT") || apikeystr == "" {
		fmt.Println("Go to https://home.openweathermap.org/, create an account and get an API key. Afterwards, paste it into 'config.txt'.")
		return
	}

	var message string

	if !showforecast {
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

		var data weather

		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("Error parsing JSON response.")
			return
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
		message += "Sunrise: " + time.Unix(int64(data.Sys.Sunrise), 0).Format("15:04") + " GMT+2\n"
		message += "Sunset: " + time.Unix(int64(data.Sys.Sunset), 0).Format("15:04") + " GMT+2\n"
	} else {
		fullurl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s", urlcity, apikeystr)
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

		var data weatherforecast

		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("Error parsing JSON response.")
			//return
		}

		message += "Location: " + data.City.Name + ", " + data.City.Country + "\n"
		for _, forecast := range data.List {
			t := time.Unix(forecast.Dt, 0).UTC()
			if t.Hour() == 12 { // Check if the time is 12:00 PM
				message += fmt.Sprintf("Date: %s, Temp: %s, Weather: %s\n", t.Format("2006-01-02"), conv_K(forecast.Main.Temp), forecast.Weather[0].Description)
			}
		}

	}

	fmt.Println(message)
	os.WriteFile("cache.txt", []byte(message), 0644)

	fmt.Println("Press any key to continue...")
	fmt.Scanln()

}
