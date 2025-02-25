package models

import (
	"errors"
	"time"
)

type Device struct {
	Config DeviceConfig
	Auth   DeviceAuth
}

type DeviceConfig struct {
	Mac             string
	Ip              string
	Name            string
	Port            uint16
	TemperatureUnit string
}

func (input *DeviceConfig) Validate() error {
	if len(input.Mac) != 12 {
		return errors.New("mac address is wrong")
	}

	if input.TemperatureUnit != Celsius && input.TemperatureUnit != Fahrenheit {
		return errors.New("unknown temperature unit")
	}

	return nil
}

type DeviceAuth struct {
	LastMessageId int
	DevType       int
	Id            [4]byte
	Key           []byte
	Iv            []byte
}

type DeviceStatusHass struct {
	FanMode       string
	SwingMode     string
	Mode          string
	Temperature   float32
	DisplaySwitch string
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

func (raw DeviceStatusRaw) ConvertToDeviceStatusHass() (mqttStatus DeviceStatusHass) {
	var deviceStatusMqtt DeviceStatusHass

	// Temperature
	deviceStatusMqtt.Temperature = raw.Temperature

	// Modes
	if raw.Power == StatusOff {
		deviceStatusMqtt.Mode = "off"
	} else {
		status, ok := ModeStatuses[int(raw.Mode)]
		if ok {
			deviceStatusMqtt.Mode = status
		} else {
			deviceStatusMqtt.Mode = "error"
		}
	}

	// Fan Status
	fanStatus, ok := FanStatuses[int(raw.FanSpeed)]
	if ok {
		deviceStatusMqtt.FanMode = fanStatus
	} else {
		deviceStatusMqtt.FanMode = "error"
	}

	if raw.Mute == StatusOn {
		deviceStatusMqtt.FanMode = "mute"
	}

	if raw.Turbo == StatusOn {
		deviceStatusMqtt.FanMode = "turbo"
	}

	// Swing Modes
	verticalFixationStatus, ok := VerticalFixationStatuses[int(raw.FixationVertical)]
	if ok {
		deviceStatusMqtt.SwingMode = verticalFixationStatus
	}

	// Display Status
	// Attention. Inverted logic
	// Byte 0 - turn ON, Byte 1 - turn OFF
	if raw.Display == 1 {
		deviceStatusMqtt.DisplaySwitch = "OFF"
	} else {
		deviceStatusMqtt.DisplaySwitch = "ON"
	}

	return deviceStatusMqtt
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

func (input *UpdateFanModeInput) Validate() error {
	var fanModes = []string{"auto", "low", "medium", "high", "turbo", "mute"}

	for _, fanMode := range fanModes {
		if fanMode == input.FanMode {
			return nil
		}
	}

	return ErrorInvalidParameterFanMode
}

type UpdateModeInput struct {
	Mac  string
	Mode string
}

func (input UpdateModeInput) Validate() error {
	var modes = []string{"auto", "off", "cool", "heat", "dry", "fan_only"}

	for _, mode := range modes {
		if mode == input.Mode {
			return nil
		}
	}

	return ErrorInvalidParameterMode
}

type UpdateSwingModeInput struct {
	Mac       string
	SwingMode string
}

func (input *UpdateSwingModeInput) Validate() error {
	_, ok := VerticalFixationStatusesInvert[input.SwingMode]
	if !ok {
		return ErrorInvalidParameterSwingMode
	}

	return nil

}

type UpdateTemperatureInput struct {
	Mac         string
	Temperature float32
}

func (input UpdateTemperatureInput) Validate() error {
	if input.Temperature > 32 || input.Temperature < 16 {
		return ErrorInvalidParameterTemperature
	}
	return nil
}

type UpdateDeviceStatesInput struct {
	Mac         string
	FanMode     *string
	SwingMode   *string
	Mode        *string
	Temperature *float32
	IsDisplayOn *bool
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

type PublishStatesOnHomeAssistantRestartInput struct {
	Status string
}

type UpdateDisplaySwitchInput struct {
	Mac    string
	Status string
}

func (input *UpdateDisplaySwitchInput) Validate() error {
	if input.Status != "ON" && input.Status != "OFF" {
		return ErrorInvalidParameterDisplayStatus
	}
	return nil
}
