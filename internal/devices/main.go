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
	SetDeviceProfile(deviceId string, profile string) error
	GetDeviceTopics(deviceId string) *TopicsInfo
	GetDeviceTopicValues(deviceId string) *TopicsValues
	RebootDevice(deviceId string) error
	UpdateDevice(deviceId string, version string) error
}

type devices struct {
	client      mqtt.Client
	info        map[string]RawDeviceInfo
	diag        map[string]DeviceDiag
	status      map[string]string
	profile     map[string]string
	topicInfo   map[string]TopicsInfo
	topicValues map[string]TopicsValues
}

var deviceTopicRegExp = regexp.MustCompile("homething/([0-9a-f]+)/device/(.*)")
var topicsRegExp = regexp.MustCompile("homething/([0-9a-f]+)/(.*)")

func NewDevices(connection string) Devices {
	devices := &devices{
		info:        map[string]RawDeviceInfo{},
		diag:        map[string]DeviceDiag{},
		status:      map[string]string{},
		profile:     map[string]string{},
		topicInfo:   map[string]TopicsInfo{},
		topicValues: map[string]TopicsValues{},
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(connection)
	opts.SetClientID("htManager")
	opts.SetDefaultPublishHandler(devices.handleMessage)
	opts.SetAutoReconnect(true)
	opts.OnConnect = devices.handleConnect
	devices.client = mqtt.NewClient(opts)
	if token := devices.client.Connect(); token.Wait() && token.Error() != nil {
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
	deviceArray := make([]DeviceInfo, 0, len(d.info))
	for deviceId, rawDevice := range d.info {
		var lastSeen *time.Time
		if diag, ok := d.diag[deviceId]; ok {
			lastSeen = diag.LastSeen
		}
		deviceArray = append(deviceArray, rawDevice.toDeviceInfo(deviceId, lastSeen))
	}
	return deviceArray
}

func (d *devices) GetDeviceInfo(deviceId string) *DeviceInfo {
	if rawDevice, ok := d.info[deviceId]; ok {
		var lastSeen *time.Time
		if diag, ok := d.diag[deviceId]; ok {
			lastSeen = diag.LastSeen
		}

		device := rawDevice.toDeviceInfo(deviceId, lastSeen)
		return &device
	}
	return nil
}

func (d *devices) isDeviceKnown(deviceId string) bool {
	_, ok := d.info[deviceId]
	return ok
}

func (d *devices) GetDeviceDiag(deviceId string) *DeviceDiag {
	if d.isDeviceKnown(deviceId) {
		if diag, ok := d.diag[deviceId]; ok {
			return &diag
		}
	}
	return nil
}

func (d *devices) GetDeviceStatus(deviceId string) *string {
	if d.isDeviceKnown(deviceId) {
		if status, ok := d.status[deviceId]; ok {
			return &status
		}
	}
	return nil
}

func (d *devices) GetDeviceProfile(deviceId string) *string {
	if d.isDeviceKnown(deviceId) {
		if profile, ok := d.profile[deviceId]; ok {
			return &profile
		}
	}
	return nil
}

func (d *devices) SetDeviceProfile(deviceId string, profile string) error {
	profileBin, err := encodeProfile(profile)
	if err != nil {
		return fmt.Errorf("failed to encode profile: %s", err)
	}
	command := append([]byte("setprofile\x00"), profileBin...)
	t := d.client.Publish(fmt.Sprintf("homething/%s/device/ctrl", deviceId), 0, false, command)
	if !t.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("timeout waiting for response from broker")
	}
	if err := t.Error(); err != nil {
		return err
	}
	return nil
}

func (d *devices) GetDeviceTopics(deviceId string) *TopicsInfo {
	if d.isDeviceKnown(deviceId) {
		if topics, ok := d.topicInfo[deviceId]; ok {
			return &topics
		}
	}
	return nil
}

func (d *devices) GetDeviceTopicValues(deviceId string) *TopicsValues {
	if d.isDeviceKnown(deviceId) {
		if values, ok := d.topicValues[deviceId]; ok {
			return &values
		}
	}
	return nil
}

func (d *devices) RebootDevice(deviceId string) error {
	t := d.client.Publish(fmt.Sprintf("homething/%s/device/ctrl", deviceId), 0, false, []byte("restart"))
	if !t.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("timeout waiting for response from broker")
	}
	if err := t.Error(); err != nil {
		return err
	}
	return nil
}

func (d *devices) UpdateDevice(deviceId string, version string) error {
	return fmt.Errorf("not implemented")
}
