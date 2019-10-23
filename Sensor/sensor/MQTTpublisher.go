package sensor

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//Connect with the MQTT broker
//provided a URL and clientID

func (s *Sensor) Connect(clientId string) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://hivemq:1883").SetClientID(clientId)
	opts.SetUsername("admin")
	opts.SetPassword("hivemq")

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to the Broker")

	// Set MQTT client to sensor
	s.client = c
}

func (s *Sensor) MqttPublish() {
	//fmt.Printf("%s \n", s.value)
	if token := s.client.Publish(string(s.sType), 1, false, []byte(s.value)); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}
