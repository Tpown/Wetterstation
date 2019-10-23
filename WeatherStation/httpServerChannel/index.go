package httpServerChannel

import (
	"encoding/json"
	"fmt"

	"../buffers"
	"../http"
	wsc "../weatherServiceChannel"
)

func StartHTTPServer(port string, weatherServiceChannel *wsc.WeatherServiceChannel) {
	var addr = "0.0.0.0:" + port
	var maxQueuedConnecions = 1024

	server := http.NewServer(addr, maxQueuedConnecions)

	// map url to SensorBuffer
	urlsToBuffers := make(map[string]*buffers.SensorBuffer)
	tempBuffer, _ := buffers.GetBuffer(buffers.TEMP_BUFFER)
	windBuffer, _ := buffers.GetBuffer(buffers.WIND_BUFFER)
	rainBuffer, _ := buffers.GetBuffer(buffers.RAIN_BUFFER)
	humBuffer, _ := buffers.GetBuffer(buffers.HUM_BUFFER)

	urlsToBuffers["/temp"] = tempBuffer
	urlsToBuffers["/wind"] = windBuffer
	urlsToBuffers["/rain"] = rainBuffer
	urlsToBuffers["/humidity"] = humBuffer

	server.SetGetHandler(handleGetRequest(urlsToBuffers, weatherServiceChannel))
	fmt.Println("Server running on port: " + port)

	err := server.Run()
	if err != nil {
		fmt.Printf("Something went wrong..\n")
	}
}

func handleGetRequest(urlsToBuffers map[string]*buffers.SensorBuffer, weatherServiceChannel *wsc.WeatherServiceChannel) func(req http.Request) http.Response {
	return func(req http.Request) http.Response {

		if sensorBuffer, ok := urlsToBuffers[req.Path()]; ok {
			fmt.Println(req.Method())
			fmt.Println(req.HTTPVersion())
			fmt.Println(req.Path())
			res := http.NewResponse()

			currentValueIndex := sensorBuffer.CurrentIndex - 1
			if currentValueIndex < 0 {
				currentValueIndex = 0
			}

			res.AddHeader("Content-Type", "application/json")
			res.SetBody(sensorBuffer.ValuesBuffer[currentValueIndex])
			res.SetStatusCode(200)

			return res
		}

		if req.Path() == "/last-report" {
			res := http.NewResponse()
			report, _ := weatherServiceChannel.GetLastReport()

			value, _ := json.Marshal(report)

			res.AddHeader("Content-Type", "application/json")
			res.SetBody(string(value))
			res.SetStatusCode(200)

			return res
		}

		return http.NewErrorResponse(404, "Unsupported URL")
	}
}
