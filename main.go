package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	OpenWeatherApiConfigKey string `json:"OpenWeatherApiConfigKey"`
}

type weatherData struct {
	Name        string `json:"name"`
	Coordinates struct {
		Longitude float64 `json:"lon"`
		Latitude  float64 `json:"lat"`
	} `json:"coord"`

	Main struct {
		Temperature_K float64 `json:"temp"`
		Pressure_hPa  float64 `json:"pressure"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println("error while reading filename...", err)
		return apiConfigData{}, err
	}

	var c = apiConfigData{}

	err1 := json.Unmarshal(bytes, &c)

	if err1 != nil {
		fmt.Println("Error during json unmarshaling: ", err)
		return apiConfigData{}, err
	}

	return c, nil
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Need Weather Details of any city? You're Covered!!!")
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig("apiConfig")

	if err != nil {
		fmt.Println("Error: could not load the api key from .apiConfig file...", err)
		return weatherData{}, err
	}

	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherApiConfigKey + "&q=" + city)

	if err != nil {
		fmt.Println("Error in fetching data through the url: ", err)
		return weatherData{}, err
	}

	defer res.Body.Close()

	var weather weatherData

	err = json.NewDecoder(res.Body).Decode(&weather)

	if err != nil {
		fmt.Println("Error while decoding response body", err)
		return weatherData{}, err
	}

	return weather, nil

}
func main() {

	http.HandleFunc("/status", status)

	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			fmt.Println("Error in executing query: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8001", nil)
}
