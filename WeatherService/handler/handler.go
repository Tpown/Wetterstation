package handler

import (
  "fmt"
  "errors"
  "context"
  ws "../gen-go/weatherService"
  "../db"
)

type WeatherServiceHandler struct{
 Persistence *db.Persistence
}

// For now we are ignoring the sessionToken -> login related
func (wsh *WeatherServiceHandler) SendWeatherReport(ctx context.Context, report *ws.WeatherReport, sessionToken int64) (r bool, err error) {
  fmt.Println(report.DateTime)
  fmt.Println(ws.ReportFromString(report.Report.String()))
  fmt.Printf("Temperature: %f\n", report.Temperature)
  fmt.Printf("Humidity: %d\n", report.Humidity)
  fmt.Printf("WindStrength: %d\n", report.WindStrength)
  fmt.Printf("Rainfall: %f\n", report.Rainfall)
  fmt.Printf("Atmosphericpressure: %d\n", report.Atmosphericpressure)
  fmt.Printf("WindDirection: %d\n", report.WindDirection)

  // Persit to DB
  r = wsh.Persistence.SaveReport(report.Temperature, report.Humidity, report.WindStrength, report.Rainfall, report.DateTime)
  if r == false {
    return false, errors.New("Unable to persist the report")
  }

  return r, nil
}

func (wsh *WeatherServiceHandler) ReceiveForecastFor(ctx context.Context, userId int64, time ws.DateTime) (r *ws.WeatherReport, err error) {
  report := wsh.Persistence.GetLastReport()

  // initiliaze a report
  weatherReport := ws.NewWeatherReport()
  // skipping Location
  weatherReport.Location = nil
  weatherReport.Report = 1
  weatherReport.DateTime = report.DateTime
  weatherReport.Temperature = report.Temperature
  weatherReport.Humidity = report.Humidity
  weatherReport.WindStrength = report.WindStrength
  weatherReport.Rainfall = report.Rainfall

  return weatherReport, nil
}

func (wsh *WeatherServiceHandler) CheckWeatherWarnings(ctx context.Context, userId int64) (r ws.WeatherWarning, err error) {
  return -1, nil
}

func (wsh *WeatherServiceHandler) SendWarning(ctx context.Context, systemWarning ws.SystemWarning, userId int64) (r bool, err error) {
  return false, nil
}


// TODO: Still need to implement the following methods!

// @return r - The session token
func (wsh *WeatherServiceHandler) Login(ctx context.Context, location *ws.Location) (r int64, err error) {
  fmt.Println("TODO: Implement DB and connect")
  return -1, nil
}

func (wsh *WeatherServiceHandler) Logout(ctx context.Context, sessionToken int64) (r bool, err error) {
  fmt.Println("TODO: Implement DB and connect")
  return false, nil
}
