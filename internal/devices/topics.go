package devices

import "fmt"

func (d *devices) handleTopicMessage(deviceId string, topic string, payload []byte) {
	fmt.Printf("Handling state topic %s for device %s\n", topic, deviceId)
}
