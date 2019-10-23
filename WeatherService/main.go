package main

import (
	"fmt"
	"os"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"

	"./db"
	ws "./gen-go/weatherService"
	"./handler"
)

func connectToDatabase() *db.Persistence {
	persistence := db.NewPersistence()
	attempts := 3
	for attempts != 0 {
		err := persistence.Connect()
		if nil != err {
			fmt.Println("Unable to connect...")
			fmt.Println(err.Error())
			attempts--
			time.Sleep(2 * time.Second)
		} else {
			// break from loop
			break
		}
	}
	// make sure that the table is created
	persistence.CreateTables()
	return persistence
}

var persistence *db.Persistence

// the main takes port as an argument
func main() {
	port := os.Args[1]

	if "" == port {
		fmt.Println("Need to specify port for the Weather Service!")
		return
	}

	// connect to db
	persistence = connectToDatabase()

	addr := "0.0.0.0:" + port

	// Declare the serialization protocol
	//protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	protocolFactory := thrift.NewTJSONProtocolFactory()

	// Declare the transport method to use
	transportFactory := thrift.NewTBufferedTransportFactory(8192)

	// Get a tansport
	transport, _ := thrift.NewTServerSocket(addr)

	// Implements the interface to our service
	handler := &handler.WeatherServiceHandler{Persistence: persistence}

	// Tell the Thrift processor which interface implementation to use
	processor := ws.NewWeatherProcessor(handler)

	// Start the server
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", addr)
	server.Serve()
}
