package updates

import (
	"errors"
	"htManager/internal/devices"
	"os"
	"reflect"
	"testing"
)

type FakeDirEntry struct {
	name     string
	fileMode os.FileMode
}

func (f *FakeDirEntry) Name() string {
	return f.name
}

func (f *FakeDirEntry) IsDir() bool {
	return f.fileMode.IsDir()
}

func (f *FakeDirEntry) Type() os.FileMode {
	return f.fileMode
}
func (f *FakeDirEntry) Info() (os.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func NewFakeDirEntryFile(name string) *FakeDirEntry {
	return &FakeDirEntry{
		name:     name,
		fileMode: 0,
	}
}

func NewFakeDirEntryDir(name string) *FakeDirEntry {
	return &FakeDirEntry{
		name:     name,
		fileMode: os.ModeDir,
	}
}

func Test_findMatches(t *testing.T) {
	type args struct {
		files      []os.DirEntry
		deviceInfo *devices.DeviceInfo
	}

	deviceInfoFlash1MB := &devices.DeviceInfo{
		Id:           "flash1MB",
		DeviceType:   "esp8266",
		Capabilities: []string{"flash1MB"},
	}

	deviceInfoFlash4MB := &devices.DeviceInfo{
		Id:           "flash4MB",
		DeviceType:   "esp8266",
		Capabilities: []string{"flash4MB"},
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: ">1MB Flash",
			args: args{
				files: []os.DirEntry{
					NewFakeDirEntryFile("homething.esp8266.v1.0.0.ota"),
					NewFakeDirEntryFile("homething.esp8266.v1.1.0.ota"),
					NewFakeDirEntryFile("homething.esp32.v1.2.0.ota"),
					NewFakeDirEntryFile("nomatch.txt"),
					NewFakeDirEntryDir("homething.esp8266.v2.0.0.ota"),
				},
				deviceInfo: deviceInfoFlash4MB,
			},
			want: []string{
				"v1.0.0",
				"v1.1.0",
			},
		},
		{
			name: "1MB Flash",
			args: args{
				files: []os.DirEntry{
					NewFakeDirEntryFile("homething.esp8266.app1.v1.0.0.ota"),
					NewFakeDirEntryFile("homething.esp8266.app2.v1.0.0.ota"),
					NewFakeDirEntryFile("homething.esp8266.app2.v1.1.0.ota"),
					NewFakeDirEntryFile("homething.esp8266.app1.v1.1.0.ota"),
					NewFakeDirEntryFile("homething.esp8266.app1.v2.0.0.ota"),
				},
				deviceInfo: deviceInfoFlash1MB,
			},
			want: []string{
				"v1.0.0",
				"v1.1.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findMatches(tt.args.files, tt.args.deviceInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}
