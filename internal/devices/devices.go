package devices

import (
	"encoding/json"
	"errors"
	"github.com/fxamacker/cbor/v2"
	"log"
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

type TopicInfo map[string]int

type RawTopicDescription struct {
	_   struct{} `cbor:",toarray"`
	Pub TopicInfo
	Sub TopicInfo
}

type TopicDescription struct {
	Pub TopicInfo `json:"pub"`
	Sub TopicInfo `json:"sub"`
}

type RawTopicsInfo struct {
	_                struct{} `cbor:",toarray"`
	TopicDescription map[int]RawTopicDescription
	Topics           map[string]int
}

type TopicsInfo struct {
	Topics map[string]TopicDescription `json:"topics"`
}

const InvalidTopicType = -1

var (
	InvalidPubTopicError        = errors.New("invalid pub topic")
	InvalidTypeForPubTopicError = errors.New("invalid topic type for pub topic")
)

func (d *devices) handleDeviceMessage(deviceId string, topic string, payload []byte) {
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
		d.info[deviceId] = info
		now := time.Now()
		d.sendUpdateMessage(deviceId, InfoUpdateMessage, info.toDeviceInfo(deviceId, &now))
	}
}

func (d *devices) handleDeviceMessageDiag(deviceId string, payload []byte) {
	diag := DeviceDiag{}
	if json.Unmarshal(payload, &diag) == nil {
		now := time.Now()
		diag.LastSeen = &now
		d.diag[deviceId] = diag

		d.sendUpdateMessage(deviceId, DiagUpdateMessage, diag)
	}
}

func (d *devices) handleDeviceMessageStatus(deviceId string, payload []byte) {
	status := string(payload)
	d.status[deviceId] = status
	d.sendUpdateMessage(deviceId, StatusUpdateMessage, status)
}

func (d *devices) handleDeviceMessageTopics(deviceId string, payload []byte) {
	rawTopicsInfo := RawTopicsInfo{}
	if err := cbor.Unmarshal(payload, &rawTopicsInfo); err == nil {
		topicsInfo := TopicsInfo{Topics: make(map[string]TopicDescription)}
		for name, descriptionId := range rawTopicsInfo.Topics {
			if rawTopicDescription, ok := rawTopicsInfo.TopicDescription[descriptionId]; ok {
				topicsInfo.Topics[name] = TopicDescription{
					Pub: rawTopicDescription.Pub,
					Sub: rawTopicDescription.Sub,
				}
			}
		}
		d.topicInfo[deviceId] = topicsInfo
	} else {
		log.Printf("Topics: CBOR unmarshal failed %v\n", err)
	}
}

func (d *devices) handleDeviceMessageProfile(deviceId string, payload []byte) {
	if profile, err := decodeProfile(payload); err == nil {
		d.profile[deviceId] = profile
	} else {
		log.Printf("Profile: CBOR unmarshal failed %v\n", err)
	}
}

func (d *devices) sendUpdateMessage(deviceId string, updateType string, data any) {
	event := DeviceUpdateEvent{
		Id:   deviceId,
		Type: updateType,
		Data: data,
	}
	log.Printf("Sending update message: %s for %s", updateType, deviceId)
	d.lock.Lock()
	for _, client := range d.updateClients {
		client.DeviceUpdated(event)
	}
	d.lock.Unlock()
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

func (t *TopicsInfo) getPubTopicType(topic string) int {
	entries := strings.Split(topic, "/")
	if topicInfo, ok := t.Topics[entries[0]]; ok {
		switch len(entries) {
		case 1:
			if topicType, ok := topicInfo.Pub[""]; ok {
				return topicType
			}
			break
		case 2:
			if topicType, ok := topicInfo.Pub[entries[1]]; ok {
				return topicType
			}
			break
		default:
			break
		}
	}
	return InvalidTopicType
}

func (t *TopicsInfo) isValidPubTopic(topic string) bool {
	return t.getPubTopicType(topic) != -1
}

func (t *TopicsInfo) convertPubTopicValue(topic string, data []byte) (any, error) {
	topicType := t.getPubTopicType(topic)
	if topicType == InvalidTopicType {
		return nil, InvalidPubTopicError
	}
	switch topicType {
	case 0:
		break
	}
	return nil, InvalidTypeForPubTopicError
}
