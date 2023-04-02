package models

import (
	models_service "github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/service/models"
	"time"
)

type ReadDeviceConfigInput struct {
	Mac string
}

type ReadDeviceConfigReturn struct {
	Config models_service.DeviceConfig
}

type UpsertDeviceConfigInput struct {
	Config models_service.DeviceConfig
}

type ReadDeviceAuthInput struct {
	Mac string
}

type ReadDeviceAuthReturn struct {
	Auth models_service.DeviceAuth
}

type UpsertDeviceAuthInput struct {
	Mac  string
	Auth models_service.DeviceAuth
}

type UpsertAmbientTempInput struct {
	Mac         string
	Temperature int8
}

type ReadAmbientTempInput struct {
	Mac string
}

type ReadAmbientTempReturn struct {
	Temperature int8
}

type ReadDeviceStatusInput struct {
	Mac string
}

type ReadDeviceStatusReturn struct {
	Status models_service.DeviceStatusMqtt
}

type UpsertDeviceStatusInput struct {
	Mac    string
	Status models_service.DeviceStatusMqtt
}

type ReadDeviceStatusRawInput struct {
	Mac string
}

type ReadDeviceStatusRawReturn struct {
	Status models_service.DeviceStatusRaw
}

type UpsertDeviceStatusRawInput struct {
	Mac    string
	Status models_service.DeviceStatusRaw
}

type MqttModeMessage struct {
	UpdatedAt time.Time
	Mode      string
}

type UpsertMqttModeMessageInput struct {
	Mac  string
	Mode MqttModeMessage
}

type MqttFanModeMessage struct {
	UpdatedAt time.Time
	FanMode   string
}

type UpsertMqttFanModeMessageInput struct {
	Mac     string
	FanMode MqttFanModeMessage
}

type MqttSwingModeMessage struct {
	UpdatedAt time.Time
	SwingMode string
}

type UpsertMqttSwingModeMessageInput struct {
	Mac       string
	SwingMode MqttSwingModeMessage
}

type MqttTemperatureMessage struct {
	UpdatedAt   time.Time
	Temperature float32
}

type UpsertMqttTemperatureMessageInput struct {
	Mac         string
	Temperature MqttTemperatureMessage
}

type ReadMqttMessageInput struct {
	Mac string
}

type ReadMqttMessageReturn struct {
	Temperature *MqttTemperatureMessage
	SwingMode   *MqttSwingModeMessage
	FanMode     *MqttFanModeMessage
	Mode        *MqttModeMessage
}

type UpsertDeviceAvailabilityInput struct {
	Mac          string
	Availability string
}

type ReadDeviceAvailabilityInput struct {
	Mac string
}

type ReadDeviceAvailabilityReturn struct {
	Availability string
}

type ReadAuthedDevicesReturn struct {
	Macs []string
}
