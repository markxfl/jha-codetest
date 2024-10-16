package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// MetadataResponse holds the structure of the initial API response
type MetadataResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

// ForecastPeriod represents a single period in the weather forecast
type ForecastPeriod struct {
	Name          string `json:"name"`
	Temperature   int    `json:"temperature"`
	ShortForecast string `json:"shortForecast"`
}

// ForecastResponse holds the structure of the forecast API response
type ForecastResponse struct {
	Properties struct {
		Periods []ForecastPeriod `json:"periods"`
	} `json:"properties"`
}

// TemperatureRange defines a structure for mapping temperature ranges
type TemperatureRange struct {
	Min int
	Max int
}

// Define the temperature categories based on ranges
var temperatureDescriptions = map[TemperatureRange]string{
	{Min: 12, Max: 32}:  "very cold",
	{Min: 32, Max: 50}:  "cold",
	{Min: 40, Max: 60}:  "moderate",
	{Min: 60, Max: 80}:  "warm",
	{Min: 80, Max: 95}:  "hot",
	{Min: 95, Max: 120}: "very hot",
}

// GetTemperatureDescription returns the temperature description based on input temperature
func GetTemperatureDescription(temp int) string {
	for rangeKey, description := range temperatureDescriptions {
		if temp >= rangeKey.Min && temp <= rangeKey.Max {
			return description
		}
	}
	return "unknown temperature range"
}

// FetchURL makes a GET request to the given URL and returns the response body as bytes
func FetchURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// GetForecastURL fetches the forecast endpoint from the metadata API
func GetForecastURL(metadataURL string) (string, error) {
	body, err := FetchURL(metadataURL)
	if err != nil {
		return "", err
	}

	var metadata MetadataResponse
	if err := json.Unmarshal(body, &metadata); err != nil {
		return "", err
	}

	return metadata.Properties.Forecast, nil
}

// GetTodaysWeatherForecast fetches the forecast data and extracts today's weather forecast
func GetTodaysWeatherForcast(forecastURL string) (*ForecastPeriod, error) {
	body, err := FetchURL(forecastURL)
	if err != nil {
		return nil, err
	}

	var forecast ForecastResponse
	if err := json.Unmarshal(body, &forecast); err != nil {
		return nil, err
	}

	if len(forecast.Properties.Periods) > 0 {
		// return first period in the array (should be today's weather forecast)
		return &forecast.Properties.Periods[0], nil
	}
	return nil, fmt.Errorf("no forecast periods found")
}

// API handler for weather forecast
func forecastHandler(w http.ResponseWriter, r *http.Request) {
	// Get latitude and longitude from query parameters
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	// Validate parameters
	if latStr == "" || lonStr == "" {
		http.Error(w, "Missing latitude or longitude", http.StatusBadRequest)
		return
	}

	// Convert latitude and longitude to float64
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	// Construct the metadata URL
	metadataURL := fmt.Sprintf("https://api.weather.gov/points/%f,%f", lat, lon)

	// Get the forecast URL
	forecastURL, err := GetForecastURL(metadataURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching forecast URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Get today's Weather Forecast
	var period *ForecastPeriod
	period, err = GetTodaysWeatherForcast(forecastURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching today's temperature: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the temperature in JSON format
	response := map[string]interface{}{
		"latitude":               lat,
		"longitude":              lon,
		"shortForecast":          period.ShortForecast,
		"temperature":            period.Temperature,
		"temperatureDescription": GetTemperatureDescription(period.Temperature),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Register the forecast handler
	http.HandleFunc("/forecast", forecastHandler)

	// Start the HTTP server
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
