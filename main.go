package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func extractValue(text, pattern string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return "N/A"
}

func forecastExtract(message string, i int) (string, string, string) {
	pattern := `Date:\s(\d{4}-\d{2}-\d{2}),\sTemp:\s([\d.]+)°C,\sWeather:\s([a-zA-Z]+)`

	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(message, -1)

	if i < 0 || i >= len(matches) {
		return "N/A", "N/A", "N/A"
	}

	return matches[i][1], matches[i][2], matches[i][3]
}

func emojiCond(cond string) string {
	switch {
	case cond == "Clear":
		cond = "⛅"
	case cond == "Clouds":
		cond = "☁️"
	case cond == "Rain":
		cond = "🌧️"
	case cond == "Drizzle":
		cond = "☔"
	case cond == "Thunderstorm":
		cond = "⛈️"
	case cond == "Snow":
		cond = "🌧️"
	case cond == "Mist" || cond == "Fog" || cond == "Smoke" || cond == "Haze":
		cond = "🌫️"
	case cond == "Tornado":
		cond = "🌪️"
	case cond == "Squall":
		cond = "💨"
	case cond == "Ash":
		cond = "🌋"
	}
	return cond
}

func dayofweek(data string) string {
	parsdat1, err := time.Parse("2006-01-02", data)
	parsdat1 = parsdat1.AddDate(0, 0, -1)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "N/A"
	}

	return parsdat1.Weekday().String()
}

func createhtml(message string) string {
	loc := extractValue(message, `Location: (.+)`)
	cond := extractValue(message, `Weather: (.+)`)
	temp := extractValue(message, `Temperature: ([\d.]+)°C`)
	feel := extractValue(message, `Feels like: ([\d.]+)°C`)
	maxt := extractValue(message, `Max temperature: ([\d.]+)°C`)
	mint := extractValue(message, `Min temperature: ([\d.]+)°C`)
	pres := extractValue(message, `Pressure: ([\d.]+) hPa`)
	hum := extractValue(message, `Humidity: ([\d.]+)%`)
	wind := extractValue(message, `Wind: ([\d.]+) km/h`)
	direc := extractValue(message, `From: (\w+)`)
	vis := extractValue(message, `Visibility: ([\d.]+) km`)

	dat1, temp1, cond1 := forecastExtract(message, 1)
	dat2, temp2, cond2 := forecastExtract(message, 2)
	dat3, temp3, cond3 := forecastExtract(message, 3)
	dat4, temp4, cond4 := forecastExtract(message, 4)

	cond = emojiCond(cond)
	cond1 = emojiCond(cond1)
	cond2 = emojiCond(cond2)
	cond3 = emojiCond(cond3)
	cond4 = emojiCond(cond4)

	basehtml, err := os.ReadFile("index.html")
	if err != nil {
		return "Error 404!"
	}

	dat1 = "Tomorrow"
	dat2 = dayofweek(dat2)
	dat3 = dayofweek(dat3)
	dat4 = dayofweek(dat4)

	return fmt.Sprintf(string(basehtml),
		loc, cond, temp, feel, mint, maxt, pres, hum, wind, direc, vis,
		cond1, temp1, dat1,
		cond2, temp2, dat2,
		cond3, temp3, dat3,
		cond4, temp4, dat4)
}

func queryCity(w http.ResponseWriter, r *http.Request) {
	querytext := r.URL.Query().Get("city")
	if querytext == "" {
		fmt.Fprintf(w, "%s", createhtml(querycity(getIpLocation())))
		return
	}
	fmt.Fprintf(w, "%s", createhtml(querycity(strings.TrimRight(querytext, " .,/\\()*&^%$#@!?<>"))))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	basehtml, err := os.ReadFile("index.html")
	if err != nil {
		fmt.Fprintf(w, "%s", "Error 404!")
	}
	fmt.Fprintf(w, "%s", basehtml)
}

func main() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/query", queryCity)

	http.ListenAndServe(":8090", nil)
}
