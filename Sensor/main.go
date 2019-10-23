package main

import (
	"fmt"
	"os"
	"time"

	"./sensor"
)

func main() {
	arg := os.Args[1]

	var s sensor.Sensor
	switch arg {
	case "-t":
		s = sensor.New(sensor.TEMPERATURE, "", "localhost", 6001)
		s.Connect("client-t")
		fmt.Println("Starting Sensor:" + sensor.TEMPERATURE)
	case "-h":
		s = sensor.New(sensor.HUMIDITY, "", "localhost", 6002)
		s.Connect("client-h")
		fmt.Println("Starting Sensor:" + sensor.HUMIDITY)
	case "-w":
		s = sensor.New(sensor.WINDSPEED, "", "localhost", 6003)
		s.Connect("client-w")
		fmt.Println("Starting Sensor:" + sensor.WINDSPEED)
	case "-r":
		s = sensor.New(sensor.RAINFALL, "", "localhost", 6004)
		s.Connect("client-r")
		fmt.Println("Starting Sensor:" + sensor.RAINFALL)
	default:
		os.Exit(1)
	}

	// start running the goroutine
	for {
		s.CreateRandomValues()
		// comment the fllowing to use UDP Channel
		//s.ClientSensorUPDConnection()
		s.MqttPublish()
		time.Sleep(time.Second)
	}
}
