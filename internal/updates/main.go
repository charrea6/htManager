package updates

import (
	"htManager/internal/devices"
	"log"
	"os"
	"regexp"
)

type UpdateManager interface {
	AvailableUpdatesForDevice(deviceInfo *devices.DeviceInfo) []string
}

type updateManagerImpl struct {
	Path string
}

const flash1MB = "flash1MB"

func NewUpdateManager(path string) UpdateManager {
	return &updateManagerImpl{Path: path}
}

func (u *updateManagerImpl) AvailableUpdatesForDevice(deviceInfo *devices.DeviceInfo) []string {
	files, err := os.ReadDir(u.Path)
	if err != nil {
		log.Printf("Failed to read contents of dir %s, error %s", u.Path, err)
		return nil
	}
	return findMatches(files, deviceInfo)
}

func findMatches(files []os.DirEntry, deviceInfo *devices.DeviceInfo) []string {
	expStr := `homething\.` + deviceInfo.DeviceType + `\.`
	if deviceInfo.HasCapability(flash1MB) {
		expStr = expStr + `app([12])\.`
	}
	expStr += `(.+)\.ota`
	exp := regexp.MustCompile(expStr)

	matches := make([][]string, 0)

	for _, entry := range files {
		if !entry.Type().IsRegular() {
			continue
		}
		match := exp.FindStringSubmatch(entry.Name())
		if match != nil {
			matches = append(matches, match)
		}
	}
	result := make([]string, 0)
	if deviceInfo.HasCapability(flash1MB) {
		versionsToAppN := make(map[string]int)
		for _, match := range matches {
			appN := 2
			if match[1] == "1" {
				appN = 1
			}
			if currentAppN, ok := versionsToAppN[match[2]]; ok {
				currentAppN |= appN
				if currentAppN == 3 {
					result = append(result, match[2])
				}
			} else {
				versionsToAppN[match[2]] = appN
			}
		}
	} else {
		for _, match := range matches {
			result = append(result, match[1])
		}
	}
	return result
}
