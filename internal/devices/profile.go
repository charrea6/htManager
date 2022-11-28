package devices

import (
	"github.com/fxamacker/cbor/v2"
	"gopkg.in/yaml.v2"
)

type ProfileEntry map[string]any
type ProfileEntries []ProfileEntry

type Profile struct {
	_       struct{} `cbor:",toarray"`
	Version int
	Profile map[string]ProfileEntries
}

func decodeProfile(payload []byte) (string, error) {
	profile := Profile{}
	if err := cbor.Unmarshal(payload, &profile); err == nil {
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
	profile.Version = 2
	profile.Profile = make(map[string]ProfileEntries)
	if err := yaml.Unmarshal([]byte(profileStr), profile.Profile); err != nil {
		return nil, err
	}
	bytes, err := cbor.Marshal(profile)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
