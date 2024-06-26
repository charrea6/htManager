package devices

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"regexp"
	"sync"
	"time"
)

const (
	InfoUpdateMessage    = "info"
	DiagUpdateMessage    = "diag"
	TopicsUpdateMessage  = "topics"
	ValueUpdateMessage   = "value"
	StatusUpdateMessage  = "status"
	DeviceRemovedMessage = "removed"
)

type DeviceInfo struct {
	Id           string     `json:"id"`
	LastSeen     *time.Time `json:"lastSeen,omitempty"`
	Description  string     `json:"description"`
	IPAddr       string     `json:"ip_addr"`
	Version      string     `json:"version"`
	DeviceType   string     `json:"deviceType"`
	Memory       uint       `json:"memory"`
	Capabilities []string   `json:"capabilities"`
}

type DeviceUpdateEvent struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Data any    `json:"data"`
}

type UpdateNotificationClient interface {
	DeviceUpdated(event DeviceUpdateEvent)
}

type Devices interface {
	GetDevices() []DeviceInfo
	RemoveDevice(deviceId string) error
	GetDeviceInfo(deviceId string) *DeviceInfo
	GetDeviceDiag(deviceId string) *DeviceDiag
	GetDeviceStatus(deviceId string) *string
	GetDeviceProfile(deviceId string) *string
	SetDeviceProfile(deviceId string, profile string) error
	GetDeviceTopics(deviceId string) *TopicsInfo
	GetDeviceTopicValues(deviceId string) *TopicsValues
	RebootDevice(deviceId string) error
	UpdateDevice(deviceId string, version string) error
	RegisterUpdateNotificationClient(client UpdateNotificationClient)
	UnregisterUpdateNotificationClient(client UpdateNotificationClient)
}

type devices struct {
	client        mqtt.Client
	info          map[string]RawDeviceInfo
	diag          map[string]DeviceDiag
	status        map[string]string
	profile       map[string]string
	topicInfo     map[string]TopicsInfo
	topicValues   map[string]TopicsValues
	lock          sync.Mutex
	updateClients []UpdateNotificationClient
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
	opts.SetClientID(fmt.Sprintf("htManager-%d", rand.Int()))
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
	t := d.client.Publish(fmt.Sprintf("homething/%s/device/ctrl", deviceId), 0, false, []byte("update "+version))
	if !t.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("timeout waiting for response from broker")
	}
	if err := t.Error(); err != nil {
		return err
	}
	return nil
}

func (d *devices) RegisterUpdateNotificationClient(client UpdateNotificationClient) {
	d.lock.Lock()
	d.updateClients = append(d.updateClients, client)
	d.lock.Unlock()
}

func (d *devices) UnregisterUpdateNotificationClient(client UpdateNotificationClient) {
	d.lock.Lock()
	for idx, value := range d.updateClients {
		if value == client {
			d.updateClients[idx] = d.updateClients[len(d.updateClients)-1]
			d.updateClients = d.updateClients[:len(d.updateClients)-1]
		}
	}
	d.lock.Unlock()
}

func (d *devices) RemoveDevice(deviceId string) error {
	if !d.isDeviceKnown(deviceId) {
		return fmt.Errorf("device %s not found", deviceId)
	}
	delete(d.info, deviceId)
	delete(d.diag, deviceId)
	delete(d.status, deviceId)
	delete(d.profile, deviceId)
	delete(d.topicInfo, deviceId)
	topicValues := d.topicValues[deviceId]
	delete(d.topicValues, deviceId)
	for primaryTopic, topicValues := range topicValues {
		for topic, _ := range topicValues {
			var topicPath string
			if topic == "" {
				topicPath = fmt.Sprintf("homething/%s/%s", deviceId, primaryTopic)
			} else {
				topicPath = fmt.Sprintf("homething/%s/%s/%s", deviceId, primaryTopic, topic)
			}

			t := d.client.Publish(topicPath, 0, true, []byte{})
			if !t.WaitTimeout(10 * time.Second) {
				return fmt.Errorf("timeout waiting for response from broker")
			}
			if err := t.Error(); err != nil {
				return err
			}
		}
	}
	for _, topic := range []string{"diag", "status", "profile", "topics", "info"} {
		topicPath := fmt.Sprintf("homething/%s/device/%s", deviceId, topic)
		t := d.client.Publish(topicPath, 0, true, []byte{})
		if !t.WaitTimeout(10 * time.Second) {
			return fmt.Errorf("timeout waiting for response from broker")
		}
		if err := t.Error(); err != nil {
			return err
		}
	}
	d.cleanupHomeAssistant(deviceId)
	d.sendUpdateMessage(deviceId, DeviceRemovedMessage, nil)
	return nil
}

func (d *devices) cleanupHomeAssistant(deviceId string) {
	topic := fmt.Sprintf("homeassistant/+/%s/#", deviceId)
	d.client.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		if message.Retained() && len(message.Payload()) > 0 {
			d.client.Publish(message.Topic(), 0, true, []byte{})
		}
	})
	time.AfterFunc(10*time.Second, func() {
		d.client.Unsubscribe(topic)
	})
}
