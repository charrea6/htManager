package devices

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

type ProfileEntry map[string]any
type ProfileEntries []ProfileEntry

type Profile struct {
	Version string                    `json:"version"`
	Profile map[string]ProfileEntries `json:"components"`
}

func decodeProfile(payload []byte) (string, error) {
	profile := Profile{}
	if err := json.Unmarshal(payload, &profile); err == nil {
		if profileBytes, err := yaml.Marshal(&profile.Profile); err == nil {
			return string(profileBytes), nil

		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func encodeProfile(profileStr string) ([]byte, error) {
	profile := Profile{}
	profile.Version = "1.0"
	profile.Profile = make(map[string]ProfileEntries)
	if err := yaml.Unmarshal([]byte(profileStr), profile.Profile); err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(profile)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
