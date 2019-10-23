package sensor

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const serverIPaddress string = "weather-station"
const serverPort int = 6000

// Sensor represents a request to run a command
type Sensor struct {
	sType     sensorType
	value     string
	ipaddress string
	port      int
	client    mqtt.Client
}

//Create a new type SensorType
type sensorType string

//Implement some method on it
//which returns the not exported type
func (st sensorType) isSensorType() sensorType {
	return st
}

//SensorType interface which is exported
type SensorType interface {
	isSensorType() sensorType
}

//Enum for Seonsorttypes: TEMPERATURE, HUMIDITY, WINDSPEED, RAINFALL
const (
	TEMPERATURE sensorType = "Temperatur"
	HUMIDITY    sensorType = "Luftfeuchtigkeit"
	WINDSPEED   sensorType = "Windgeschw"
	RAINFALL    sensorType = "Niederschlag"
)

//New Sensor Constructor
func New(sensortype sensorType, value string, ipaddress string, port int) Sensor {
	sensor := Sensor{sensortype, value, ipaddress, port, nil}
	return sensor
}

// Print is only for debugging purposes
func (s Sensor) Print() {
	fmt.Printf("%s test \n", s.sType)
}

func (s *Sensor) CreateRandomValues() {
	rand.Seed(time.Now().UnixNano())
	switch s.sType {
	case TEMPERATURE:
		{
			rnd := rand.Intn(10)
			s.value = strconv.Itoa(rnd + 20)
			//fmt.Println(s.Value)
		}
	case HUMIDITY:
		{
			rnd := rand.Intn(100-50) + 50
			s.value = strconv.Itoa(rnd + 100)
			//fmt.Println(s.Value)
		}
	case WINDSPEED:
		{
			rnd := rand.Intn(1000-500) + 500
			s.value = strconv.Itoa(rnd)
			//fmt.Println(s.Value)
		}
	case RAINFALL:
		{
			rnd := rand.Intn(100-10) + 10
			s.value = strconv.Itoa(rnd + 1000)
			//fmt.Println(s.Value)
		}
	}
}
