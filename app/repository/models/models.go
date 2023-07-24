package models

import (
	models_service "github.com/ArtemVladimirov/broadlinkac2mqtt/app/service/models"
	"time"
)

type Device struct {
	Config          models_service.DeviceConfig
	Auth            *models_service.DeviceAuth
	DeviceStatus    DeviceStatus
	DeviceStatusRaw *models_service.DeviceStatusRaw
	MqttLastMessage MqttStatus
}

type DeviceStatus struct {
	Availability *string
	AmbientTemp  *float32
}

type MqttStatus struct {
	FanMode       *MqttFanModeMessage
	SwingMode     *MqttSwingModeMessage
	Mode          *MqttModeMessage
	Temperature   *MqttTemperatureMessage
	DisplaySwitch *MqttDisplaySwitchMessage
}

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
	Temperature float32
}

type ReadAmbientTempInput struct {
	Mac string
}

type ReadAmbientTempReturn struct {
	Temperature float32
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

type MqttDisplaySwitchMessage struct {
	UpdatedAt   time.Time
	IsDisplayOn bool
}

type UpsertMqttDisplaySwitchMessageInput struct {
	Mac           string
	DisplaySwitch MqttDisplaySwitchMessage
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
	IsDisplayOn *MqttDisplaySwitchMessage
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
