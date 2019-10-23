package MQTTChannel

import (
	"fmt"

	"../buffers"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// we are using the http port as additional ID
func Connect(clientId string) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker("tcp://hivemq:1883").SetClientID("station-server-" + clientId)
	opts.SetUsername("admin")
	opts.SetPassword("hivemq")

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return c
}

func Subscribe(client mqtt.Client, topic string) {
	if token := client.Subscribe(topic, 0, onSubscribe); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

}

func onSubscribe(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	value := string(message.Payload())

	fmt.Println(topic)
	fmt.Println(value)

	MQTTBuffer(topic, value)
}

func MQTTBuffer(topic, value string) {
	var bufferId string
	switch topic {
	case "Temperatur":
		bufferId = buffers.TEMP_BUFFER
	case "Niederschlag":
		bufferId = buffers.RAIN_BUFFER
	case "Windgeschw":
		bufferId = buffers.WIND_BUFFER
	case "Luftfeuchtigkeit":
		bufferId = buffers.HUM_BUFFER
	}

	buffers.AddToBuffer(bufferId, value)
}
