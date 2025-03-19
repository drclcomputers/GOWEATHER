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

const C0 = 273.15 //0°C = 273K

func conv_K_C(degrees float64) float64 { //convert K to C
	return degrees - C0
}

func directwind(deg float64) string { //Get wind direction
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((deg+22.5)/45.0) % 8
	return directions[index]
}

func main() {
	//Welcome page and city input **********************************************************************
	fmt.Println("GOWEATHER - ver 0.1.2")
	apikey := "INSERT WEATHER API KEY HERE"
	var city string
	if len(os.Args) == 1 {
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
	} else if len(os.Args) > 1 {
		city = os.Args[1]
	}

	urlcity := url.QueryEscape(city)
	//**************************************************************************************************

	//Fetching data from API ***************************************************************************
	if apikey == "INSERT YOUR API KEY HERE" {
		fmt.Println("Go to https://home.openweathermap.org/, create an account and get an API key. Afterwards, insert it here, into the apikey variable.")
		return
	}

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
	//**************************************************************************************************

	//Parse JSON to string ******************************
	var data weather

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error parsing JSON response.")
		return
	}
	//****************************************************

	var message string
	message += "Location: " + data.Name + ", " + data.Sys.Country + "\n"
	if len(data.Weatherdata) > 0 {
		message += "Weather: " + data.Weatherdata[0].Description + "\n"
	} else {
		message += "Weather description unavailable!\n"
	}
	message += "Temperature: " + fmt.Sprintf("%.0f", conv_K_C(data.Main.Temp)) + "°C\n"
	message += "Feels like: " + fmt.Sprintf("%.0f", conv_K_C(data.Main.FeelsLike)) + "°C\n"
	message += "Max temperature: " + fmt.Sprintf("%.0f", conv_K_C(data.Main.MaxTemp)) + "°C\n"
	message += "Min temperature: " + fmt.Sprintf("%.0f", conv_K_C(data.Main.MinTemp)) + "°C\n"
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
	message += "Wind: " + fmt.Sprintf("%.1f", data.Wind.Speed*3.6) + "km/h -> From: "
	message += directwind(data.Wind.Deg) + "\n"
	message += "Sunrise: " + time.Unix(int64(data.Sys.Sunrise), 0).Format("15:04") + " GMT+2\n"
	message += "Sunset: " + time.Unix(int64(data.Sys.Sunset), 0).Format("15:04") + " GMT+2\n"

	//Print Data *************************************************************************
	fmt.Println(message)

	fmt.Println("Press any key to continue...")
	fmt.Scanln()

}
