#!/bin/sh
trap 'kill %1; kill %2' INT
go run sensorMain.go -t &
go run sensorMain.go -r &
go run sensorMain.go -h &
go run sensorMain.go -w 
