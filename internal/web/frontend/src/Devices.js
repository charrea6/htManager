export class DeviceList {
    constructor() {
        this.devices = [];
        this.selectedDevice = null;
        this.deviceListUpdated = null;
        this.deviceUpdated = null;
        this.connected = false;
        this.pending = null;
        this.connectWS();
    }

    selectDevice(deviceId, deviceUpdated) {
        if (!this.connected) {
            this.pending = () => { this.selectDevice(deviceId, deviceUpdated)};
            return;
        }
        this.selectedDevice = deviceId;
        this.deviceUpdated = deviceUpdated;
        for (const info of this.devices) {
            if (info.id === deviceId) {
                deviceUpdated('info', info);
                break;
            }
        }
        this.ws.send(JSON.stringify({'cmd': 'selectDevice', 'id': deviceId}));
    }

    unselectDevice(deviceId) {
        if (!this.connected) {
            this.pending = () => { this.unselectDevice(deviceId)};
            return;
        }

        this.deviceUpdated = null;
        this.selectedDevice = null;
        this.ws.send(JSON.stringify({'cmd': 'unselectDevice', 'id': deviceId}));
    }

    getDeviceInfo(deviceId) {
        for (const d of this.devices) {
            if (d.id === deviceId) {
                return d;
            }
        }
        return null;
    }

    connectWS() {
        let loc = window.location, protocol;
        if (loc.protocol === "https:") {
            protocol = "wss:";
        } else {
            protocol = "ws:";
        }
        this.ws = new WebSocket(`${protocol}//${loc.host}/api/ws`);
        this.ws.onopen = () => {
            this.connected = true;
            if (this.pending != null) {
                this.pending();
                this.pending = null;
            } else if (this.selectedDevice !== null) {
                this.ws.send(JSON.stringify({'cmd': 'selectDevice', 'id': this.selectedDevice}));
            }
        }
        this.ws.onclose = () => {
            this.connected = false;
            this.connectWS();
        };
        this.ws.onerror = () => {
            this.ws.close();
        }
        this.ws.onmessage = (event) => { this.processWSMessage(event)};
    }

    processWSMessage(event) {
        let msg = JSON.parse(event.data);
        switch (msg.type) {
            case 'init':
                this.handleInit(msg.data);
                break
            case 'lastSeen':
                this.handleLastSeen(msg);
                break;
            case 'diag':
            case 'info':
            case 'status':
            case 'topics':
            case 'values':
            case 'value':
                this.handleDeviceUpdate(msg);
                break;
            default:
                console.log(`Unknown message ${msg.type}`);
                break;
        }
    }

    handleInit(devices) {
        this.devices = devices;
        if (this.deviceListUpdated != null) {
            this.deviceListUpdated(devices);
        }
        if (this.deviceUpdated != null) {
            for (const d of devices) {
                if (d.id === this.selectedDevice){
                    this.deviceUpdated('info', d);
                    break;
                }
            }
        }
    }

    handleLastSeen(msg) {
        let devices = []
        this.devices.forEach((value) => {
            if (value.id === msg.id) {
                value = { ...value}
                value.lastSeen = msg.data.lastSeen;
            }
            devices.push(value)
        });
        this.devices = devices;
        if (this.deviceListUpdated != null) {
            this.deviceListUpdated(this.devices);
        }
    }

    handleDeviceUpdate(msg) {
        if (msg.id === this.selectedDevice) {
            this.deviceUpdated(msg.type, msg.data);
        }
        if (msg.type === 'info') {
            let devices = [];
            this.devices.forEach((value) => {
                if (value.id === msg.id) {
                    value = msg.data;
                }
                devices.push(value)
            });
            this.devices = devices;
            if (this.deviceListUpdated != null) {
                this.deviceListUpdated(this.devices);
            }
        }
        if (msg.type === 'removed') {
            let devices = [];
            this.devices.forEach((value) => {
                if (value.id !== msg.id) {
                    devices.push(value)
                }
            });
            this.devices = devices;
            if (this.deviceListUpdated != null) {
                this.deviceListUpdated(this.devices);
            }
        }
    }
}
