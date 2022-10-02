package devices

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"regexp"
	"time"
)

type DeviceInfo struct {
	Id           string     `json:"id"`
	LastSeen     *time.Time `json:"lastSeen,omitempty"`
	Description  string     `json:"description"`
	IPAddr       string     `json:"ip_addr"`
	Version      string     `json:"version"`
	Capabilities []string   `json:"capabilities"`
}

type Devices interface {
	GetDevices() []DeviceInfo
	GetDeviceInfo(deviceId string) *DeviceInfo
	GetDeviceDiag(deviceId string) *DeviceDiag
	GetDeviceStatus(deviceId string) *string
	GetDeviceProfile(deviceId string) *string
}

type devices struct {
	devicesInfo    map[string]RawDeviceInfo
	devicesDiag    map[string]DeviceDiag
	devicesStatus  map[string]string
	devicesProfile map[string]string
}

var deviceTopicRegExp = regexp.MustCompile("homething/([0-9a-f]+)/device/(.*)")
var topicsRegExp = regexp.MustCompile("homething/([0-9a-f]+)/(.*)")

func NewDevices(connection string) Devices {
	devices := &devices{
		devicesInfo:    make(map[string]RawDeviceInfo),
		devicesDiag:    make(map[string]DeviceDiag),
		devicesStatus:  make(map[string]string),
		devicesProfile: make(map[string]string),
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(connection)
	opts.SetClientID("htManager")
	opts.SetDefaultPublishHandler(devices.handleMessage)
	opts.SetAutoReconnect(true)
	opts.OnConnect = devices.handleConnect
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return devices
}

func (d *devices) handleConnect(client mqtt.Client) {
	client.Subscribe("homething/#", 0, d.handleMessage)
}

func (d *devices) handleMessage(client mqtt.Client, msg mqtt.Message) {
	if matches := deviceTopicRegExp.FindStringSubmatch(msg.Topic()); len(matches) > 0 {
		d.handleDeviceMessage(matches[1], matches[2], msg.Payload())
	} else if matches := topicsRegExp.FindStringSubmatch(msg.Topic()); len(matches) > 0 {
		d.handleTopicMessage(matches[1], matches[2], msg.Payload())
	} else {
		fmt.Printf("Unmatched topic %s", msg.Topic())
	}
}

func (d *devices) GetDevices() []DeviceInfo {
	deviceArray := make([]DeviceInfo, 0, len(d.devicesInfo))
	for deviceId, rawDevice := range d.devicesInfo {
		var lastSeen *time.Time
		if diag, ok := d.devicesDiag[deviceId]; ok {
			lastSeen = diag.LastSeen
		}
		deviceArray = append(deviceArray, rawDevice.toDeviceInfo(deviceId, lastSeen))
	}
	return deviceArray
}

func (d *devices) GetDeviceInfo(deviceId string) *DeviceInfo {
	if rawDevice, ok := d.devicesInfo[deviceId]; ok {
		var lastSeen *time.Time
		if diag, ok := d.devicesDiag[deviceId]; ok {
			lastSeen = diag.LastSeen
		}

		device := rawDevice.toDeviceInfo(deviceId, lastSeen)
		return &device
	}
	return nil
}

func (d *devices) isDeviceKnown(deviceId string) bool {
	_, ok := d.devicesInfo[deviceId]
	return ok
}

func (d *devices) GetDeviceDiag(deviceId string) *DeviceDiag {
	if d.isDeviceKnown(deviceId) {
		if diag, ok := d.devicesDiag[deviceId]; ok {
			return &diag
		}
	}
	return nil
}

func (d *devices) GetDeviceStatus(deviceId string) *string {
	if d.isDeviceKnown(deviceId) {
		if status, ok := d.devicesStatus[deviceId]; ok {
			return &status
		}
	}
	return nil
}

func (d *devices) GetDeviceProfile(deviceId string) *string {
	if d.isDeviceKnown(deviceId) {
		if profile, ok := d.devicesProfile[deviceId]; ok {
			return &profile
		}
	}
	return nil
}
