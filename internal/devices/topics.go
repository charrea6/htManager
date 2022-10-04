package devices

import (
	"strings"
)

type TopicValues map[string]any
type TopicsValues map[string]TopicValues

func (d *devices) handleTopicMessage(deviceId string, topic string, payload []byte) {
	topicsInfo, ok := d.topicInfo[deviceId]
	if !ok {
		return
	}
	if !topicsInfo.isValidPubTopic(topic) {
		return
	}
	topicsValues, ok := d.topicValues[deviceId]
	if !ok {
		topicsValues = make(TopicsValues)
		d.topicValues[deviceId] = topicsValues
	}
	topicsValues.setValue(topic, string(payload))
}

func (t *TopicsValues) setValue(topic string, value any) {
	entries := strings.Split(topic, "/")
	primaryTopic, ok := (*t)[entries[0]]
	if !ok {
		primaryTopic = make(TopicValues)
		(*t)[entries[0]] = primaryTopic
	}
	switch len(entries) {
	case 1:
		primaryTopic[""] = value
		break
	case 2:
		primaryTopic[entries[1]] = value
		break
	default:
		return
	}
}
