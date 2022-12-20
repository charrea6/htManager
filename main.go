package main

import (
	"flag"
	"fmt"
	"htManager/internal/devices"
	"htManager/internal/updates"
	"htManager/internal/web"
)

var mqttHost string
var mqttPort int
var updatesPath string

func main() {
	flag.StringVar(&mqttHost, "host", "localhost", "hostname of the MQTT server to connect to.")
	flag.StringVar(&updatesPath, "updates-path", ".", "Location of homething OTA files.")
	flag.IntVar(&mqttPort, "port", 1883, "Port number of the MQTT server to connect to.")
	flag.Parse()
	devicesManager := devices.NewDevices(fmt.Sprintf("tcp://%s:%d", mqttHost, mqttPort))
	updateManager := updates.NewUpdateManager(updatesPath)
	web.InitWebServer(devicesManager, updateManager)
}
