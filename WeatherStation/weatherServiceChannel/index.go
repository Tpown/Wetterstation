package weatherServiceChannel

import (
	"context"
	"fmt"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"

	ws "../gen-go/weatherService"
)

type WeatherServiceChannel struct {
	weatherClient *ws.WeatherClient
	socket        *thrift.TSocket
}

func (this *WeatherServiceChannel) Connect(hostPort string) error {
	socket, err := thrift.NewTSocket(hostPort)
	if nil != err {
		fmt.Println("Unable to create socket: " + err.Error())
		return err
	}
	err = socket.Open()
	if nil != err {
		fmt.Println("Unable to open socket: " + err.Error())
		return err
	}

	this.socket = socket

	// we are using the JSON Protocol
	protocolFactory := thrift.NewTJSONProtocolFactory()

	// set up the client
	this.weatherClient = ws.NewWeatherClientFactory(socket, protocolFactory)
	return nil
}

func (this *WeatherServiceChannel) Close() error {
	this.weatherClient = nil
	return this.socket.Close()
}

// TODO: We are passing the location to this function in future
func (this *WeatherServiceChannel) Login() {
	location := ws.NewLocation()
	location.LocationID = 1
	location.Name = "Darmstadt"
	location.Latitude = 34.5
	location.Longitude = 23.4
	var locationDescription string
	locationDescription = "Nice"
	location.Description = &locationDescription

	r, err := this.weatherClient.Login(context.Background(), location)

	if nil != err {
		fmt.Println("Unable to create task: " + err.Error())
		return
	}
	fmt.Print("Done with: ")
	fmt.Println(r)

	fmt.Println("Done")
}

// TODO: Add Session Token from Login
// @return r - if the SendReport was successful
// @return err - corresponding error on fail
func (this *WeatherServiceChannel) SendReport(temp float64, humidity int16, windStrength int16, rainfall float64) (r bool, err error) {
	// initiliaze a report
	report := ws.NewWeatherReport()
	// set now date
	report.DateTime = time.Now().String()
	// skipping Location
	report.Location = nil

	report.Report = 1

	// filling report data
	report.Temperature = temp
	report.Humidity = humidity
	report.WindStrength = windStrength
	report.Rainfall = rainfall

	// TODO: replace -1 with the right session token
	return this.weatherClient.SendWeatherReport(context.Background(), report, -1)
}

func (this *WeatherServiceChannel) GetLastReport() (r *ws.WeatherReport, err error) {
	return this.weatherClient.ReceiveForecastFor(context.Background(), -1, "")
}
