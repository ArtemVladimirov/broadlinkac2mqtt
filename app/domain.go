package app

import (
	"context"
	models_mqtt "github.com/ArtemVladimirov/broadlinkac2mqtt/app/mqtt/models"
	models_cache "github.com/ArtemVladimirov/broadlinkac2mqtt/app/repository/models"
	models_service "github.com/ArtemVladimirov/broadlinkac2mqtt/app/service/models"
	models_web "github.com/ArtemVladimirov/broadlinkac2mqtt/app/webClient/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
)

type MqttSubscriber interface {
	UpdateFanModeCommandTopic(ctx context.Context, logger *slog.Logger) mqtt.MessageHandler
	UpdateSwingModeCommandTopic(ctx context.Context, logger *slog.Logger) mqtt.MessageHandler
	UpdateModeCommandTopic(ctx context.Context, logger *slog.Logger) mqtt.MessageHandler
	UpdateTemperatureCommandTopic(ctx context.Context, logger *slog.Logger) mqtt.MessageHandler
	UpdateDisplaySwitchCommandTopic(ctx context.Context, logger *slog.Logger) mqtt.MessageHandler

	GetStatesOnHomeAssistantRestart(ctx context.Context, logger *slog.Logger) mqtt.MessageHandler
}

type MqttPublisher interface {
	PublishClimateDiscoveryTopic(ctx context.Context, logger *slog.Logger, input models_mqtt.PublishClimateDiscoveryTopicInput) error
	PublishSwitchDiscoveryTopic(ctx context.Context, logger *slog.Logger, input models_mqtt.PublishSwitchDiscoveryTopicInput) error
	PublishAmbientTemp(ctx context.Context, logger *slog.Logger, input *models_mqtt.PublishAmbientTempInput) error
	PublishTemperature(ctx context.Context, logger *slog.Logger, input *models_mqtt.PublishTemperatureInput) error
	PublishMode(ctx context.Context, logger *slog.Logger, input *models_mqtt.PublishModeInput) error
	PublishSwingMode(ctx context.Context, logger *slog.Logger, input *models_mqtt.PublishSwingModeInput) error
	PublishFanMode(ctx context.Context, logger *slog.Logger, input *models_mqtt.PublishFanModeInput) error
	PublishAvailability(ctx context.Context, logger *slog.Logger, input *models_mqtt.PublishAvailabilityInput) error
	PublishDisplaySwitch(ctx context.Context, logger *slog.Logger, input *models_mqtt.PublishDisplaySwitchInput) error
}

type Service interface {
	PublishDiscoveryTopic(ctx context.Context, logger *slog.Logger, input *models_service.PublishDiscoveryTopicInput) error
	CreateDevice(ctx context.Context, logger *slog.Logger, input *models_service.CreateDeviceInput) error
	AuthDevice(ctx context.Context, logger *slog.Logger, input *models_service.AuthDeviceInput) error
	GetDeviceAmbientTemperature(ctx context.Context, logger *slog.Logger, input *models_service.GetDeviceAmbientTemperatureInput) error
	GetDeviceStates(ctx context.Context, logger *slog.Logger, input *models_service.GetDeviceStatesInput) error

	UpdateFanMode(ctx context.Context, logger *slog.Logger, input *models_service.UpdateFanModeInput) error
	UpdateMode(ctx context.Context, logger *slog.Logger, input *models_service.UpdateModeInput) error
	UpdateSwingMode(ctx context.Context, logger *slog.Logger, input *models_service.UpdateSwingModeInput) error
	UpdateTemperature(ctx context.Context, logger *slog.Logger, input *models_service.UpdateTemperatureInput) error
	UpdateDisplaySwitch(ctx context.Context, logger *slog.Logger, input *models_service.UpdateDisplaySwitchInput) error

	UpdateDeviceAvailability(ctx context.Context, logger *slog.Logger, input *models_service.UpdateDeviceAvailabilityInput) error

	StartDeviceMonitoring(ctx context.Context, logger *slog.Logger, input *models_service.StartDeviceMonitoringInput) error

	PublishStatesOnHomeAssistantRestart(ctx context.Context, logger *slog.Logger, input *models_service.PublishStatesOnHomeAssistantRestartInput) error
}

type WebClient interface {
	SendCommand(ctx context.Context, logger *slog.Logger, input *models_web.SendCommandInput) (*models_web.SendCommandReturn, error)
}

type Cache interface {
	UpsertDeviceConfig(ctx context.Context, logger *slog.Logger, input *models_cache.UpsertDeviceConfigInput) error
	ReadDeviceConfig(ctx context.Context, logger *slog.Logger, input *models_cache.ReadDeviceConfigInput) (*models_cache.ReadDeviceConfigReturn, error)

	UpsertDeviceAuth(ctx context.Context, logger *slog.Logger, input *models_cache.UpsertDeviceAuthInput) error
	ReadDeviceAuth(ctx context.Context, logger *slog.Logger, input *models_cache.ReadDeviceAuthInput) (*models_cache.ReadDeviceAuthReturn, error)

	UpsertAmbientTemp(ctx context.Context, logger *slog.Logger, input *models_cache.UpsertAmbientTempInput) error
	ReadAmbientTemp(ctx context.Context, logger *slog.Logger, input *models_cache.ReadAmbientTempInput) (*models_cache.ReadAmbientTempReturn, error)

	UpsertDeviceStatusRaw(ctx context.Context, logger *slog.Logger, input *models_cache.UpsertDeviceStatusRawInput) error
	ReadDeviceStatusRaw(ctx context.Context, logger *slog.Logger, input *models_cache.ReadDeviceStatusRawInput) (*models_cache.ReadDeviceStatusRawReturn, error)

	UpsertMqttModeMessage(ctx context.Context, logger *slog.Logger, input *models_cache.UpsertMqttModeMessageInput) error
	UpsertMqttSwingModeMessage(ctx context.Context, logger *slog.Logger, input *models_cache.UpsertMqttSwingModeMessageInput) error
	UpsertMqttFanModeMessage(ctx context.Context, logger *slog.Logger, input *models_cache.UpsertMqttFanModeMessageInput) error
	UpsertMqttTemperatureMessage(ctx context.Context, logger *slog.Logger, input *models_cache.UpsertMqttTemperatureMessageInput) error
	UpsertMqttDisplaySwitchMessage(ctx context.Context, logger *slog.Logger, input *models_cache.UpsertMqttDisplaySwitchMessageInput) error

	ReadMqttMessage(ctx context.Context, logger *slog.Logger, input *models_cache.ReadMqttMessageInput) (*models_cache.ReadMqttMessageReturn, error)

	UpsertDeviceAvailability(ctx context.Context, logger *slog.Logger, input *models_cache.UpsertDeviceAvailabilityInput) error
	ReadDeviceAvailability(ctx context.Context, logger *slog.Logger, input *models_cache.ReadDeviceAvailabilityInput) (*models_cache.ReadDeviceAvailabilityReturn, error)

	ReadAuthedDevices(ctx context.Context, logger *slog.Logger) (*models_cache.ReadAuthedDevicesReturn, error)
}
