package devices

import (
	"strings"
)

type TopicValues map[string]any
type TopicsValues map[string]TopicValues

type ValueUpdateEvent struct {
	TopicPath []string `json:"topic_path"`
	Value     string   `json:"value"`
}

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
	entries := topicsValues.setValue(topic, string(payload))
	if entries != nil {
		d.sendUpdateMessage(deviceId, ValueUpdateMessage, ValueUpdateEvent{
			TopicPath: entries,
			Value:     string(payload),
		})
	}
}

func (t *TopicsValues) setValue(topic string, value any) []string {
	entries := strings.Split(topic, "/")
	primaryTopic, ok := (*t)[entries[0]]
	if !ok {
		primaryTopic = make(TopicValues)
		(*t)[entries[0]] = primaryTopic
	}
	switch len(entries) {
	case 1:
		primaryTopic[""] = value
		entries = append(entries, "")
		break
	case 2:
		primaryTopic[entries[1]] = value
		break
	default:
		return nil
	}
	return entries
}
