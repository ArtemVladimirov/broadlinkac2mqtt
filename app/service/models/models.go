package models

import "time"

type Device struct {
	Config     DeviceConfig
	Auth       DeviceAuth
	StatusRaw  DeviceStatusRaw
	StatusMqtt DeviceStatusMqtt
}

type DeviceConfig struct {
	Mac  string
	Ip   string
	Name string
	Port uint16
}

type DeviceAuth struct {
	LastMessageId int
	DevType       int
	Id            [4]byte
	Key           []byte
	Iv            []byte
}

type DeviceStatusMqtt struct {
	FanMode     string
	SwingMode   string
	Mode        string
	Temperature float32
}

type DeviceStatusRaw struct {
	UpdatedAt          time.Time
	Temperature        float32
	Power              byte
	FixationVertical   byte
	Mode               byte
	Sleep              byte
	Display            byte
	Mildew             byte
	Health             byte
	FixationHorizontal byte
	FanSpeed           byte
	IFeel              byte
	Mute               byte
	Turbo              byte
	Clean              byte
}

type CreateDeviceInput struct {
	Config DeviceConfig
}

type CreateDeviceReturn struct {
	Device Device
}

type AuthDeviceInput struct {
	Mac string
}

type SendCommandInput struct {
	Command byte
	Payload []byte
	Mac     string
}

type SendCommandReturn struct {
	Payload []byte
}

type GetDeviceAmbientTemperatureInput struct {
	Mac string
}

type GetDeviceStatesInput struct {
	Mac string
}

type PublishDiscoveryTopicInput struct {
	Device DeviceConfig
}

type UpdateFanModeInput struct {
	Mac     string
	FanMode string
}

type UpdateModeInput struct {
	Mac  string
	Mode string
}

type UpdateSwingModeInput struct {
	Mac       string
	SwingMode string
}

type UpdateTemperatureInput struct {
	Mac         string
	Temperature float32
}

type UpdateDeviceStatesInput struct {
	Mac         string
	FanMode     *string
	SwingMode   *string
	Mode        *string
	Temperature *float32
}

type CreateCommandPayloadReturn struct {
	Payload []byte
}

type UpdateDeviceAvailabilityInput struct {
	Mac          string
	Availability string
}

type StartDeviceMonitoringInput struct {
	Mac string
}

type GetStatesOnHomeAssistantRestartInput struct {
	Status string
}
