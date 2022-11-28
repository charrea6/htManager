package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"htManager/internal/devices"
	"log"
	"time"
)

type WebSocketInitMessage struct {
	Type string               `json:"type"`
	Data []devices.DeviceInfo `json:"data"`
}

type WebSocketConnection struct {
	ws             *websocket.Conn
	devices        devices.Devices
	selectedDevice string
}

type WebSocketClientRequest struct {
	Cmd string `json:"cmd"`
	Id  string `json:"id"`
}

type LastSeenUpdate struct {
	LastSeen *time.Time `json:"lastSeen"`
}

func (c *WebSocketConnection) handleConnection() {
	log.Println("Handle Connection starting...")
	initMsg := WebSocketInitMessage{Type: "init", Data: c.devices.GetDevices()}
	if bytes, err := json.Marshal(initMsg); err == nil {
		if err := c.ws.WriteMessage(websocket.TextMessage, bytes); err != nil {
			log.Printf("Failed to send init message: %s\n", err)
		}
	} else {
		log.Printf("Failed to marshal init message: %s\n", err)
	}
	log.Println("Init message sent")

	c.devices.RegisterUpdateNotificationClient(c)
	defer func() { c.devices.UnregisterUpdateNotificationClient(c) }()
	log.Println("Starting to receive ws messages...")
	for {
		//Read Message from client
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		request := WebSocketClientRequest{}
		if err := json.Unmarshal(message, &request); err != nil {
			log.Printf("Failed to unmarshal request: %s", err)
			continue
		}
		switch request.Cmd {
		case "selectDevice":
			c.selectedDevice = request.Id
			if diag := c.devices.GetDeviceDiag(request.Id); diag != nil {
				c.sendUpdateMessage(devices.DeviceUpdateEvent{
					Id:   request.Id,
					Type: devices.DiagUpdateMessage,
					Data: diag,
				})
			}
			break
		case "unselectDevice":
			if c.selectedDevice == request.Id {
				c.selectedDevice = ""
			}
			break
		default:
			log.Printf("Unknown request: %s", request.Cmd)
		}
	}
	log.Println("Finished ws receive")
}

func (c *WebSocketConnection) DeviceUpdated(event devices.DeviceUpdateEvent) {
	if event.Id == c.selectedDevice {
		if err := c.sendUpdateMessage(event); err != nil {
			log.Printf("Error while sending ws message: %s", err)
			return
		}
	}
	if event.Type == devices.DiagUpdateMessage {
		if diag, ok := event.Data.(devices.DeviceDiag); ok {
			lastSeen := LastSeenUpdate{LastSeen: diag.LastSeen}
			msg := devices.DeviceUpdateEvent{
				Id:   event.Id,
				Type: "lastSeen",
				Data: lastSeen,
			}
			if err := c.sendUpdateMessage(msg); err != nil {
				log.Printf("Error while sending ws message: %s", err)
				return
			}
		}
	}
}

func (c *WebSocketConnection) sendUpdateMessage(updateMessage devices.DeviceUpdateEvent) error {
	msg, err := json.Marshal(updateMessage)
	if err != nil {
		return err
	}
	if err := c.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("Error while sending ws message: %s", err)
		return err
	}
	return nil
}
