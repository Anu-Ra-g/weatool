package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	q := ""
	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	if q == ""{
		log.Fatal("Add your city name")
	}
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=c2dfb64fba804ff994f193841230708&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		log.Fatal("Check your connection")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("Cannot call the api")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Cannot Read the response Body")
	}

	var weather Weather
	if err := json.Unmarshal(body, &weather); err != nil {
		log.Fatal(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf("%s, %s: %.0fC, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s - %.0fC, %.0f%%, %s", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)

		if hour.ChanceOfRain < 40 {
			color.Green(message)
		} else if hour.ChanceOfRain > 40 && hour.ChanceOfRain < 80{
			color.Yellow(message)
		}else {
			color.Red(message)
		}
	}

}
