# JHA Code Test - Simple Weather Forecast Service
Demo Purpose Only

## Requirements

- `go1.23.2`


## Installation and Usage
Install Go Runtime 1.23.2
https://go.dev/dl/

```bash
$ export OW_API_KEY=<your_api_key_from_openweathermap>
```
Clone this project:
```bash
$ git clone https://github.com/markxfl/jha-codetest.git
```

Run the server `main.go`

```bash
$ go run main.go
Server is running on port 8080...
```

Run curl command on the other terminal
(or past the following url to a web browser if don't have curl command) 

```bash
$ curl "http://localhost:8080/forecast?lat=27.97722457369917&lon=-82.53023750466228"
{"latitude":27.97722457369917,"longitude":-82.53023750466228,"shortForecast":"Partly Cloudy","temperature":55,"temperatureDescription":"moderate"}
```
