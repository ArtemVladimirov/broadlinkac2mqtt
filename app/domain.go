package app

import (
	"context"
	models_mqtt "github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/mqtt/models"
	models_cache "github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/repository/models"
	models_service "github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/service/models"
	models_web "github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/webClient/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

type MqttSubscriber interface {
	UpdateFanModeCommandTopic(logger *zerolog.Logger) mqtt.MessageHandler
	UpdateSwingModeCommandTopic(logger *zerolog.Logger) mqtt.MessageHandler
	UpdateModeCommandTopic(logger *zerolog.Logger) mqtt.MessageHandler
	UpdateTemperatureCommandTopic(logger *zerolog.Logger) mqtt.MessageHandler

	GetStatesOnHomeAssistantRestart(logger *zerolog.Logger) mqtt.MessageHandler
}

type MqttPublisher interface {
	PublishDiscoveryTopic(ctx context.Context, logger *zerolog.Logger, input models_mqtt.PublishDiscoveryTopicInput) error
	PublishAmbientTemp(ctx context.Context, logger *zerolog.Logger, input *models_mqtt.PublishAmbientTempInput) error
	PublishTemperature(ctx context.Context, logger *zerolog.Logger, input *models_mqtt.PublishTemperatureInput) error
	PublishMode(ctx context.Context, logger *zerolog.Logger, input *models_mqtt.PublishModeInput) error
	PublishSwingMode(ctx context.Context, logger *zerolog.Logger, input *models_mqtt.PublishSwingModeInput) error
	PublishFanMode(ctx context.Context, logger *zerolog.Logger, input *models_mqtt.PublishFanModeInput) error
	PublishAvailability(ctx context.Context, logger *zerolog.Logger, input *models_mqtt.PublishAvailabilityInput) error
}

type Service interface {
	PublishDiscoveryTopic(ctx context.Context, logger *zerolog.Logger, input *models_service.PublishDiscoveryTopicInput) error
	CreateDevice(ctx context.Context, logger *zerolog.Logger, input *models_service.CreateDeviceInput) (*models_service.CreateDeviceReturn, error)
	AuthDevice(ctx context.Context, logger *zerolog.Logger, input *models_service.AuthDeviceInput) error
	GetDeviceAmbientTemperature(ctx context.Context, logger *zerolog.Logger, input *models_service.GetDeviceAmbientTemperatureInput) error
	GetDeviceStates(ctx context.Context, logger *zerolog.Logger, input *models_service.GetDeviceStatesInput) error

	UpdateFanMode(ctx context.Context, logger *zerolog.Logger, input *models_service.UpdateFanModeInput) error
	UpdateMode(ctx context.Context, logger *zerolog.Logger, input *models_service.UpdateModeInput) error
	UpdateSwingMode(ctx context.Context, logger *zerolog.Logger, input *models_service.UpdateSwingModeInput) error
	UpdateTemperature(ctx context.Context, logger *zerolog.Logger, input *models_service.UpdateTemperatureInput) error

	UpdateDeviceAvailability(ctx context.Context, logger *zerolog.Logger, input *models_service.UpdateDeviceAvailabilityInput) error

	StartDeviceMonitoring(ctx context.Context, logger *zerolog.Logger, input *models_service.StartDeviceMonitoringInput) error

	GetStatesOnHomeAssistantRestart(ctx context.Context, logger *zerolog.Logger, input *models_service.GetStatesOnHomeAssistantRestartInput) error
}

type WebClient interface {
	SendCommand(ctx context.Context, logger *zerolog.Logger, input *models_web.SendCommandInput) (*models_web.SendCommandReturn, error)
}

type Cache interface {
	UpsertDeviceConfig(ctx context.Context, logger *zerolog.Logger, input *models_cache.UpsertDeviceConfigInput) error
	ReadDeviceConfig(ctx context.Context, logger *zerolog.Logger, input *models_cache.ReadDeviceConfigInput) (*models_cache.ReadDeviceConfigReturn, error)

	UpsertDeviceAuth(ctx context.Context, logger *zerolog.Logger, input *models_cache.UpsertDeviceAuthInput) error
	ReadDeviceAuth(ctx context.Context, logger *zerolog.Logger, input *models_cache.ReadDeviceAuthInput) (*models_cache.ReadDeviceAuthReturn, error)

	UpsertAmbientTemp(ctx context.Context, logger *zerolog.Logger, input *models_cache.UpsertAmbientTempInput) error
	ReadAmbientTemp(ctx context.Context, logger *zerolog.Logger, input *models_cache.ReadAmbientTempInput) (*models_cache.ReadAmbientTempReturn, error)

	UpsertDeviceStatus(ctx context.Context, logger *zerolog.Logger, input *models_cache.UpsertDeviceStatusInput) error
	ReadDeviceStatus(ctx context.Context, logger *zerolog.Logger, input *models_cache.ReadDeviceStatusInput) (*models_cache.ReadDeviceStatusReturn, error)

	UpsertDeviceStatusRaw(ctx context.Context, logger *zerolog.Logger, input *models_cache.UpsertDeviceStatusRawInput) error
	ReadDeviceStatusRaw(ctx context.Context, logger *zerolog.Logger, input *models_cache.ReadDeviceStatusRawInput) (*models_cache.ReadDeviceStatusRawReturn, error)

	UpsertMqttModeMessage(ctx context.Context, logger *zerolog.Logger, input *models_cache.UpsertMqttModeMessageInput) error
	UpsertMqttSwingModeMessage(ctx context.Context, logger *zerolog.Logger, input *models_cache.UpsertMqttSwingModeMessageInput) error
	UpsertMqttFanModeMessage(ctx context.Context, logger *zerolog.Logger, input *models_cache.UpsertMqttFanModeMessageInput) error
	UpsertMqttTemperatureMessage(ctx context.Context, logger *zerolog.Logger, input *models_cache.UpsertMqttTemperatureMessageInput) error

	ReadMqttMessage(ctx context.Context, logger *zerolog.Logger, input *models_cache.ReadMqttMessageInput) (*models_cache.ReadMqttMessageReturn, error)

	UpsertDeviceAvailability(ctx context.Context, logger *zerolog.Logger, input *models_cache.UpsertDeviceAvailabilityInput) error
	ReadDeviceAvailability(ctx context.Context, logger *zerolog.Logger, input *models_cache.ReadDeviceAvailabilityInput) (*models_cache.ReadDeviceAvailabilityReturn, error)

	ReadAuthedDevices(ctx context.Context, logger *zerolog.Logger) (*models_cache.ReadAuthedDevicesReturn, error)
}
