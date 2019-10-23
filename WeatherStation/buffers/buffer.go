package buffers

import (
	"errors"
	"log"
	"os"
	"strings"
)

// This package is responsible of Sensor buffers.
// It will contain the creation and handling of existing buffers

const MAX_VALUE_BUFFER = 5

// the following are common buffer types
// which will be used every where in the project
const (
	TEMP_BUFFER = "TEMP_BUFFER"
	HUM_BUFFER  = "HUM_BUFFER"
	WIND_BUFFER = "WIND_BUFFER"
	RAIN_BUFFER = "RAIN_BUFFER"
)

type SensorBuffer struct {
	ValuesBuffer [MAX_VALUE_BUFFER]string
	CurrentIndex int
}

// Map the id of the buffer to create buffer
var Buffers map[string]*SensorBuffer = make(map[string]*SensorBuffer)

func GetBuffer(bufferId string) (*SensorBuffer, error) {
	if buffer, ok := Buffers[bufferId]; ok {
		return buffer, nil
	} else {
		return nil, errors.New("The wanted Buffer does not exist")
	}
}

func AddToBuffer(bufferId, value string) {
	sensorBuffer, _ := GetBuffer(bufferId)
	sensorBuffer.ValuesBuffer[sensorBuffer.CurrentIndex] = value
	sensorBuffer.CurrentIndex++
	// make sure that we go back to 0 once we reach the limit of buffer
	if sensorBuffer.CurrentIndex >= MAX_VALUE_BUFFER {
		sensorBuffer.CurrentIndex = 0
	}
}

func CreateBuffer(bufferId string) (*SensorBuffer, error) {
	if _, ok := Buffers[bufferId]; ok {
		return nil, errors.New("Buffer exists with the same ID")
	} else {
		var newBuffer *SensorBuffer
		newBuffer = new(SensorBuffer)
		newBuffer.CurrentIndex = 0
		Buffers[bufferId] = newBuffer
		return newBuffer, nil
	}
}

func ExtractValueFromBuffer(bufferId string) string {
	if buffer, ok := Buffers[bufferId]; ok {
		// make sure we have the right index
		currentIndex := buffer.CurrentIndex - 1
		if currentIndex == -1 {
			currentIndex = 0
		}

		// make sure that the value does contains ":"
		if strings.Contains(buffer.ValuesBuffer[currentIndex], ":") {
			return strings.Split(buffer.ValuesBuffer[currentIndex], ":")[1]
		} else {
			return buffer.ValuesBuffer[currentIndex]
		}
	} else {
		return ""
	}
}

func cacheSensorBufferInLog(sensorBuffer *SensorBuffer) {
	if sensorBuffer.CurrentIndex == MAX_VALUE_BUFFER {
		// TODO: open file and append values
		/*fmt.Println("-----------------------------")
		if strings.HasPrefix(sensorBuffer.ValuesBuffer[1], "Temperatur") {
			fmt.Println("Temperatur")
		}
		fmt.Println("-----------------------------")*/

		f, err := os.OpenFile("history-values.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)

		if err != nil {
			log.Println(err)
		}

		for _, value := range sensorBuffer.ValuesBuffer {
			//fmt.Printf("%s  \n", value)

			if _, err := f.WriteString(value + "\n"); err != nil {
				log.Println(err)
			}
		}
		f.Close()

		sensorBuffer.CurrentIndex = 0
	}
}

func CacheValuesBuffer(reqChannel chan string) {
	for {
		sensorValue := <-reqChannel
		var bufferId string
		if strings.HasPrefix(sensorValue, "Temperatur") {
			bufferId = TEMP_BUFFER
		}
		if strings.HasPrefix(sensorValue, "Luftfeuchtigkeit") {
			bufferId = HUM_BUFFER
		}
		if strings.HasPrefix(sensorValue, "Windgeschw") {
			bufferId = WIND_BUFFER
		}
		if strings.HasPrefix(sensorValue, "Niederschlag") {
			bufferId = RAIN_BUFFER
		}

		buffer, _ := GetBuffer(bufferId)
		buffer.ValuesBuffer[buffer.CurrentIndex] = sensorValue
		buffer.CurrentIndex++
		cacheSensorBufferInLog(buffer)
	}
}
