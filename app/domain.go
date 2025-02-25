package app

import (
	"context"
	modelsMqtt "github.com/ArtemVladimirov/broadlinkac2mqtt/app/mqtt/models"
	modelsCache "github.com/ArtemVladimirov/broadlinkac2mqtt/app/repository/models"
	modelsService "github.com/ArtemVladimirov/broadlinkac2mqtt/app/service/models"
	modelsWeb "github.com/ArtemVladimirov/broadlinkac2mqtt/app/webClient/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttSubscriber interface {
	UpdateFanModeCommandTopic(ctx context.Context) mqtt.MessageHandler
	UpdateSwingModeCommandTopic(ctx context.Context) mqtt.MessageHandler
	UpdateModeCommandTopic(ctx context.Context) mqtt.MessageHandler
	UpdateTemperatureCommandTopic(ctx context.Context) mqtt.MessageHandler
	UpdateDisplaySwitchCommandTopic(ctx context.Context) mqtt.MessageHandler

	GetStatesOnHomeAssistantRestart(ctx context.Context) mqtt.MessageHandler
}

type MqttPublisher interface {
	PublishClimateDiscoveryTopic(ctx context.Context, input modelsMqtt.PublishClimateDiscoveryTopicInput) error
	PublishSwitchDiscoveryTopic(ctx context.Context, input modelsMqtt.PublishSwitchDiscoveryTopicInput) error
	PublishAmbientTemp(ctx context.Context, input *modelsMqtt.PublishAmbientTempInput) error
	PublishTemperature(ctx context.Context, input *modelsMqtt.PublishTemperatureInput) error
	PublishMode(ctx context.Context, input *modelsMqtt.PublishModeInput) error
	PublishSwingMode(ctx context.Context, input *modelsMqtt.PublishSwingModeInput) error
	PublishFanMode(ctx context.Context, input *modelsMqtt.PublishFanModeInput) error
	PublishAvailability(ctx context.Context, input *modelsMqtt.PublishAvailabilityInput) error
	PublishDisplaySwitch(ctx context.Context, input *modelsMqtt.PublishDisplaySwitchInput) error
}

type Service interface {
	PublishDiscoveryTopic(ctx context.Context, input *modelsService.PublishDiscoveryTopicInput) error
	CreateDevice(ctx context.Context, input *modelsService.CreateDeviceInput) error
	AuthDevice(ctx context.Context, input *modelsService.AuthDeviceInput) error
	GetDeviceAmbientTemperature(ctx context.Context, input *modelsService.GetDeviceAmbientTemperatureInput) error
	GetDeviceStates(ctx context.Context, input *modelsService.GetDeviceStatesInput) error

	UpdateFanMode(ctx context.Context, input *modelsService.UpdateFanModeInput) error
	UpdateMode(ctx context.Context, input *modelsService.UpdateModeInput) error
	UpdateSwingMode(ctx context.Context, input *modelsService.UpdateSwingModeInput) error
	UpdateTemperature(ctx context.Context, input *modelsService.UpdateTemperatureInput) error
	UpdateDisplaySwitch(ctx context.Context, input *modelsService.UpdateDisplaySwitchInput) error

	UpdateDeviceAvailability(ctx context.Context, input *modelsService.UpdateDeviceAvailabilityInput) error

	StartDeviceMonitoring(ctx context.Context, input *modelsService.StartDeviceMonitoringInput) error

	PublishStatesOnHomeAssistantRestart(ctx context.Context, input *modelsService.PublishStatesOnHomeAssistantRestartInput) error
}

type WebClient interface {
	SendCommand(ctx context.Context, input *modelsWeb.SendCommandInput) (*modelsWeb.SendCommandReturn, error)
}

type Cache interface {
	UpsertDeviceConfig(ctx context.Context, input *modelsCache.UpsertDeviceConfigInput) error
	ReadDeviceConfig(ctx context.Context, input *modelsCache.ReadDeviceConfigInput) (*modelsCache.ReadDeviceConfigReturn, error)

	UpsertDeviceAuth(ctx context.Context, input *modelsCache.UpsertDeviceAuthInput) error
	ReadDeviceAuth(ctx context.Context, input *modelsCache.ReadDeviceAuthInput) (*modelsCache.ReadDeviceAuthReturn, error)

	UpsertAmbientTemp(ctx context.Context, input *modelsCache.UpsertAmbientTempInput) error
	ReadAmbientTemp(ctx context.Context, input *modelsCache.ReadAmbientTempInput) (*modelsCache.ReadAmbientTempReturn, error)

	UpsertDeviceStatusRaw(ctx context.Context, input *modelsCache.UpsertDeviceStatusRawInput) error
	ReadDeviceStatusRaw(ctx context.Context, input *modelsCache.ReadDeviceStatusRawInput) (*modelsCache.ReadDeviceStatusRawReturn, error)

	UpsertMqttModeMessage(ctx context.Context, input *modelsCache.UpsertMqttModeMessageInput) error
	UpsertMqttSwingModeMessage(ctx context.Context, input *modelsCache.UpsertMqttSwingModeMessageInput) error
	UpsertMqttFanModeMessage(ctx context.Context, input *modelsCache.UpsertMqttFanModeMessageInput) error
	UpsertMqttTemperatureMessage(ctx context.Context, input *modelsCache.UpsertMqttTemperatureMessageInput) error
	UpsertMqttDisplaySwitchMessage(ctx context.Context, input *modelsCache.UpsertMqttDisplaySwitchMessageInput) error

	ReadMqttMessage(ctx context.Context, input *modelsCache.ReadMqttMessageInput) (*modelsCache.ReadMqttMessageReturn, error)

	UpsertDeviceAvailability(ctx context.Context, input *modelsCache.UpsertDeviceAvailabilityInput) error
	ReadDeviceAvailability(ctx context.Context, input *modelsCache.ReadDeviceAvailabilityInput) (*modelsCache.ReadDeviceAvailabilityReturn, error)

	ReadAuthedDevices(ctx context.Context) (*modelsCache.ReadAuthedDevicesReturn, error)
}
