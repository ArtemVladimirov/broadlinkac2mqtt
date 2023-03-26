package cache

import (
	"context"
	"github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/repository/models"
	models_service "github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/service/models"
	"github.com/rs/zerolog"
)

type cache struct {
	deviceConfig map[string]models_service.DeviceConfig
	deviceAuth   map[string]models_service.DeviceAuth
	// Storage for converted states
	deviceStatusMqtt map[string]models_service.DeviceStatusMqtt
	// Storages for last received states from ac
	deviceStatusRaw map[string]models_service.DeviceStatusRaw
	ambientTemp     map[string]int8
	// Storages for last mqtt message
	mqttModeMessages        map[string]models.MqttModeMessage
	mqttFanModeMessages     map[string]models.MqttFanModeMessage
	mqttSwingModeMessages   map[string]models.MqttSwingModeMessage
	mqttTemperatureMessages map[string]models.MqttTemperatureMessage
}

func NewCache() *cache {
	deviceConfig := make(map[string]models_service.DeviceConfig)
	deviceAuth := make(map[string]models_service.DeviceAuth)
	deviceStatusMqtt := make(map[string]models_service.DeviceStatusMqtt)
	deviceStatusRaw := make(map[string]models_service.DeviceStatusRaw)
	ambientTemp := make(map[string]int8)
	mqttModeMessages := make(map[string]models.MqttModeMessage)
	mqttFanModeMessages := make(map[string]models.MqttFanModeMessage)
	mqttSwingModeMessages := make(map[string]models.MqttSwingModeMessage)
	mqttTemperatureMessages := make(map[string]models.MqttTemperatureMessage)

	return &cache{
		deviceAuth:              deviceAuth,
		deviceConfig:            deviceConfig,
		deviceStatusMqtt:        deviceStatusMqtt,
		deviceStatusRaw:         deviceStatusRaw,
		ambientTemp:             ambientTemp,
		mqttModeMessages:        mqttModeMessages,
		mqttFanModeMessages:     mqttFanModeMessages,
		mqttSwingModeMessages:   mqttSwingModeMessages,
		mqttTemperatureMessages: mqttTemperatureMessages,
	}
}

func (c *cache) UpsertDeviceConfig(ctx context.Context, logger *zerolog.Logger, input *models.UpsertDeviceConfigInput) error {
	c.deviceConfig[input.Config.Mac] = input.Config
	return nil
}

func (c *cache) ReadDeviceConfig(ctx context.Context, logger *zerolog.Logger, input *models.ReadDeviceConfigInput) (*models.ReadDeviceConfigReturn, error) {
	config, ok := c.deviceConfig[input.Mac]
	if !ok {
		message := "device config not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceNotFound
	}
	return &models.ReadDeviceConfigReturn{Config: config}, nil
}

func (c *cache) UpsertDeviceAuth(ctx context.Context, logger *zerolog.Logger, input *models.UpsertDeviceAuthInput) error {
	c.deviceAuth[input.Mac] = input.Auth
	return nil
}

func (c *cache) ReadDeviceAuth(ctx context.Context, logger *zerolog.Logger, input *models.ReadDeviceAuthInput) (*models.ReadDeviceAuthReturn, error) {
	auth, ok := c.deviceAuth[input.Mac]
	if !ok {
		message := "device auth not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceNotFound
	}
	return &models.ReadDeviceAuthReturn{Auth: auth}, nil
}

func (c *cache) UpsertAmbientTemp(ctx context.Context, logger *zerolog.Logger, input *models.UpsertAmbientTempInput) error {
	c.ambientTemp[input.Mac] = input.Temperature
	return nil
}

func (c *cache) ReadAmbientTemp(ctx context.Context, logger *zerolog.Logger, input *models.ReadAmbientTempInput) (*models.ReadAmbientTempReturn, error) {
	temperature, ok := c.ambientTemp[input.Mac]
	if !ok {
		message := "ambient temperature not found in cache"
		logger.Debug().Interface("input", input).Msg(message)
		return nil, nil
	}
	return &models.ReadAmbientTempReturn{Temperature: temperature}, nil
}

func (c *cache) UpsertDeviceStatus(ctx context.Context, logger *zerolog.Logger, input *models.UpsertDeviceStatusInput) error {
	c.deviceStatusMqtt[input.Mac] = input.Status
	return nil
}

func (c *cache) ReadDeviceStatus(ctx context.Context, logger *zerolog.Logger, input *models.ReadDeviceStatusInput) (*models.ReadDeviceStatusReturn, error) {
	status, ok := c.deviceStatusMqtt[input.Mac]
	if !ok {
		message := "device not found in cache"
		logger.Debug().Interface("input", input).Msg(message)
		return nil, nil
	}
	return &models.ReadDeviceStatusReturn{Status: status}, nil
}

func (c *cache) UpsertDeviceStatusRaw(ctx context.Context, logger *zerolog.Logger, input *models.UpsertDeviceStatusRawInput) error {
	c.deviceStatusRaw[input.Mac] = input.Status
	return nil
}

func (c *cache) ReadDeviceStatusRaw(ctx context.Context, logger *zerolog.Logger, input *models.ReadDeviceStatusRawInput) (*models.ReadDeviceStatusRawReturn, error) {
	status, ok := c.deviceStatusRaw[input.Mac]
	if !ok {
		message := "device not found in cache"
		logger.Debug().Interface("input", input).Msg(message)
		return nil, nil
	}
	return &models.ReadDeviceStatusRawReturn{Status: status}, nil
}

func (c *cache) UpsertMqttModeMessage(ctx context.Context, logger *zerolog.Logger, input *models.UpsertMqttModeMessageInput) error {
	c.mqttModeMessages[input.Mac] = input.Mode
	return nil
}

func (c *cache) UpsertMqttSwingModeMessage(ctx context.Context, logger *zerolog.Logger, input *models.UpsertMqttSwingModeMessageInput) error {
	c.mqttSwingModeMessages[input.Mac] = input.SwingMode
	return nil
}

func (c *cache) UpsertMqttFanModeMessage(ctx context.Context, logger *zerolog.Logger, input *models.UpsertMqttFanModeMessageInput) error {
	c.mqttFanModeMessages[input.Mac] = input.FanMode
	return nil
}

func (c *cache) UpsertMqttTemperatureMessage(ctx context.Context, logger *zerolog.Logger, input *models.UpsertMqttTemperatureMessageInput) error {
	c.mqttTemperatureMessages[input.Mac] = input.Temperature
	return nil
}

func (c *cache) ReadMqttMessage(ctx context.Context, logger *zerolog.Logger, input *models.ReadMqttMessageInput) (*models.ReadMqttMessageReturn, error) {

	var state models.ReadMqttMessageReturn

	mode, ok := c.mqttModeMessages[input.Mac]
	if ok {
		state.Mode = &mode
	}

	swingMode, ok := c.mqttSwingModeMessages[input.Mac]
	if ok {
		state.SwingMode = &swingMode
	}

	fanMode, ok := c.mqttFanModeMessages[input.Mac]
	if ok {
		state.FanMode = &fanMode
	}

	temperature, ok := c.mqttTemperatureMessages[input.Mac]

	if ok {
		state.Temperature = &temperature
	}

	return &state, nil
}
