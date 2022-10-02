package main

import (
	"flag"
	"fmt"
	"htManager/internal/devices"
	"htManager/internal/web"
)

var mqttHost string
var mqttPort int

func main() {
	flag.StringVar(&mqttHost, "host", "localhost", "hostname of the MQTT server to connect to.")
	flag.IntVar(&mqttPort, "port", 1883, "Port number of the MQTT server to connect to.")
	flag.Parse()
	devicesManager := devices.NewDevices(fmt.Sprintf("tcp://%s:%d", mqttHost, mqttPort))
	web.InitWebServer(devicesManager)
}
