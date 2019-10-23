package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	mqttC "./MQTTChannel"
	"./buffers"
	"./httpServerChannel"
	"./socket"
	wsc "./weatherServiceChannel"
)

func startUDP() {
	hostName := "0.0.0.0"
	portNum := "6000"
	service := hostName + ":" + portNum

	udpAddr, err := net.ResolveUDPAddr("udp", service)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// setup listener for incoming UDP connection
	listener, err := net.ListenUDP("udp", udpAddr)

	// TODO: get reference of UDP socket and close it externally
	//defer listener.Close()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	fmt.Println("UDP server up and listening on port 6000")

	reqChannel := make(chan string, 4)

	go func() {
		for {
			// wait for UDP client to connect
			socket.HandleUDPConnection(listener, reqChannel)
		}
	}()

	go buffers.CacheValuesBuffer(reqChannel)
}

// start the Weather Service Socket and return a reference
func startServiceSocket(numberOfServices int) *wsc.WeatherServiceChannel {
	fmt.Println("Connecting to the weather service...")

	// We connect to the weather service
	weatherServiceChannel := &wsc.WeatherServiceChannel{}

	go func() {
		for {
			// we try to connect to one of the available services
			for i := 0; i < numberOfServices; i++ {
				hostPort := "weather-service-" + strconv.Itoa(i) + ":909" + strconv.Itoa(i)
				if err := weatherServiceChannel.Connect(hostPort); err != nil {
					fmt.Println("Was unable to connect to " + hostPort)
					fmt.Println(err.Error())
				}
			}
			// same sleep as the sensors
			time.Sleep(time.Second)

			// parse the values to correct types
			temp, _ := strconv.ParseFloat(buffers.ExtractValueFromBuffer(buffers.TEMP_BUFFER), 64)
			hum, _ := strconv.ParseInt(buffers.ExtractValueFromBuffer(buffers.HUM_BUFFER), 10, 16)
			wind, _ := strconv.ParseInt(buffers.ExtractValueFromBuffer(buffers.WIND_BUFFER), 10, 16)
			rain, _ := strconv.ParseFloat(buffers.ExtractValueFromBuffer(buffers.RAIN_BUFFER), 64)

			weatherServiceChannel.SendReport(temp, int16(hum), int16(wind), rain)

			time.Sleep(time.Second)
			// close the socket to allow for the other connections
			weatherServiceChannel.Close()
		}
	}()
	return weatherServiceChannel
}

// The program takes port for the HTTP server
// and the number of available weather services
func main() {
	availableServices := os.Args[2]
	if "" == availableServices {
		fmt.Println("Specify the number of available services")
		return
	}
	nbrServices, _ := strconv.Atoi(availableServices)
	if nbrServices <= 0 {
		fmt.Println("Number of services can not be less than 1")
		return
	}

	HTTPPort := os.Args[1]
	if "" == HTTPPort {
		fmt.Println("You need to specify the HTTP port of the Weather Station")
		return
	}

	// Initilize the buffers
	buffers.CreateBuffer(buffers.TEMP_BUFFER)
	buffers.CreateBuffer(buffers.HUM_BUFFER)
	buffers.CreateBuffer(buffers.RAIN_BUFFER)
	buffers.CreateBuffer(buffers.WIND_BUFFER)

	//startUDP()

	client := mqttC.Connect(HTTPPort)
	mqttC.Subscribe(client, "Temperatur")
	mqttC.Subscribe(client, "Luftfeuchtigkeit")
	mqttC.Subscribe(client, "Windgeschw")
	mqttC.Subscribe(client, "Niederschlag")

	weatherServiceChannel := startServiceSocket(nbrServices)

	httpServerChannel.StartHTTPServer(HTTPPort, weatherServiceChannel)
}
