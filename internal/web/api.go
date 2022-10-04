package web

import (
	"github.com/gin-gonic/gin"
	"htManager/internal/devices"
	"net/http"
)

type DeviceStatusResponse struct {
	Status string `json:"status"`
}

type DeviceProfileResponse struct {
	Profile string `json:"profile"`
}

type DeviceTopicValues struct {
	Values *devices.TopicsValues `json:"values"`
}

func initAPI(group *gin.RouterGroup, devices devices.Devices) {
	group.GET("/devices", func(context *gin.Context) {
		context.JSON(http.StatusOK, devices.GetDevices())
	})

	group.GET("/devices/:deviceId/info", func(context *gin.Context) {
		deviceId := context.Param("deviceId")
		if info := devices.GetDeviceInfo(deviceId); info == nil {
			context.Status(http.StatusNotFound)
		} else {
			context.JSON(http.StatusOK, info)
		}
	})

	group.GET("/devices/:deviceId/diag", func(context *gin.Context) {
		deviceId := context.Param("deviceId")
		if diag := devices.GetDeviceDiag(deviceId); diag == nil {
			context.Status(http.StatusNotFound)
		} else {
			context.JSON(http.StatusOK, diag)
		}
	})

	group.GET("/devices/:deviceId/status", func(context *gin.Context) {
		deviceId := context.Param("deviceId")
		if status := devices.GetDeviceStatus(deviceId); status == nil {
			context.Status(http.StatusNotFound)
		} else {
			context.JSON(http.StatusOK, DeviceStatusResponse{Status: *status})
		}
	})

	group.GET("/devices/:deviceId/profile", func(context *gin.Context) {
		deviceId := context.Param("deviceId")
		if profile := devices.GetDeviceProfile(deviceId); profile == nil {
			context.Status(http.StatusNotFound)
		} else {
			context.JSON(http.StatusOK, DeviceProfileResponse{Profile: *profile})
		}
	})

	group.GET("/devices/:deviceId/topics", func(context *gin.Context) {
		deviceId := context.Param("deviceId")
		if topics := devices.GetDeviceTopics(deviceId); topics == nil {
			context.Status(http.StatusNotFound)
		} else {
			context.JSON(http.StatusOK, topics)
		}
	})
	group.GET("/devices/:deviceId/topics/values", func(context *gin.Context) {
		deviceId := context.Param("deviceId")
		if values := devices.GetDeviceTopicValues(deviceId); values == nil {
			context.Status(http.StatusNotFound)
		} else {
			context.JSON(http.StatusOK, DeviceTopicValues{
				Values: values,
			})
		}
	})
}
