package devices

import (
	"encoding/json"
	"fmt"
	"github.com/fxamacker/cbor/v2"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

type RawDeviceInfo struct {
	IpAddr       string `json:"ip"`
	Description  string `json:"description"`
	Version      string `json:"version"`
	Capabilities string `json:"capabilities"`
}

type DeviceDiagMemInfo struct {
	Free int `json:"free"`
	Low  int `json:"low"`
}

type DeviceDiagStackInfo struct {
	Name         string `json:"name"`
	StackMinLeft int    `json:"stackMinLeft"`
}

type DeviceDiag struct {
	LastSeen *time.Time            `json:"lastSeen,omitempty"`
	Uptime   int                   `json:"uptime"`
	MemInfo  DeviceDiagMemInfo     `json:"mem"`
	TaskInfo []DeviceDiagStackInfo `json:"tasks,omitempty"`
}

type ProfileEntry map[string]any
type ProfileEntries []ProfileEntry

type Profile struct {
	_       struct{} `cbor:",toarray"`
	Version int
	Profile map[string]ProfileEntries
}

func (d *devices) handleDeviceMessage(deviceId string, topic string, payload []byte) {
	fmt.Printf("Handling info topic %s for device %s\n", topic, deviceId)
	switch topic {
	case "info":
		d.handleDeviceMessageInfo(deviceId, payload)
		break
	case "diag":
		d.handleDeviceMessageDiag(deviceId, payload)
		break
	case "status":
		d.handleDeviceMessageStatus(deviceId, payload)
		break
	case "topics":
		d.handleDeviceMessageTopics(deviceId, payload)
		break
	case "profile":
		d.handleDeviceMessageProfile(deviceId, payload)
		break
	}
}

func (d *devices) handleDeviceMessageInfo(deviceId string, payload []byte) {
	info := RawDeviceInfo{}
	if json.Unmarshal(payload, &info) == nil {
		d.devicesInfo[deviceId] = info
	}
}

func (d *devices) handleDeviceMessageDiag(deviceId string, payload []byte) {
	diag := DeviceDiag{}
	if json.Unmarshal(payload, &diag) == nil {
		now := time.Now()
		diag.LastSeen = &now
		d.devicesDiag[deviceId] = diag
	}
}

func (d *devices) handleDeviceMessageStatus(deviceId string, payload []byte) {
	d.devicesStatus[deviceId] = string(payload)
}

func (d *devices) handleDeviceMessageTopics(deviceId string, payload []byte) {

}

func (d *devices) handleDeviceMessageProfile(deviceId string, payload []byte) {
	profile := Profile{}
	if err := cbor.Unmarshal(payload, &profile); err == nil {
		fmt.Printf("Profile: %v", profile)
		if profileBytes, err := yaml.Marshal(&profile.Profile); err == nil {
			d.devicesProfile[deviceId] = string(profileBytes)
		}
	} else {
		fmt.Printf("CBOR unmarshal failed %v", err)
	}
}

func (d *RawDeviceInfo) toDeviceInfo(deviceId string, lastSeen *time.Time) DeviceInfo {
	device := DeviceInfo{
		Id:           deviceId,
		Description:  d.Description,
		IPAddr:       d.IpAddr,
		Version:      d.Version,
		Capabilities: strings.Split(d.Capabilities, ","),
		LastSeen:     lastSeen,
	}
	return device
}
