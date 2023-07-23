package subscriber

import (
	"context"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app/mqtt/models"
	models_service "github.com/ArtemVladimirov/broadlinkac2mqtt/app/service/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
)

type mqttSubscriber struct {
	mqttConfig models.ConfigMqtt
	service    app.Service
}

func NewMqttReceiver(service app.Service, mqttConfig models.ConfigMqtt) *mqttSubscriber {
	return &mqttSubscriber{
		mqttConfig: mqttConfig,
		service:    service,
	}
}

func (m *mqttSubscriber) UpdateFanModeCommandTopic(ctx context.Context, logger *zerolog.Logger) mqtt.MessageHandler {
	return mqtt.MessageHandler(func(c mqtt.Client, msg mqtt.Message) {

		mac := strings.TrimPrefix(strings.TrimSuffix(msg.Topic(), "/fan_mode/set"), m.mqttConfig.TopicPrefix+"/")

		logger.Debug().Str("device", mac).Str("payload", string(msg.Payload())).Str("topic", msg.Topic()).Msg("new update fan mode message")

		updateFanModeInput := &models_service.UpdateFanModeInput{
			Mac:     mac,
			FanMode: string(msg.Payload()),
		}

		err := m.service.UpdateFanMode(ctx, logger, updateFanModeInput)
		if err != nil {
			logger.Error().Err(err).Str("device", mac).Interface("input", updateFanModeInput).Msg("failed to update fan mode")
			return
		}
	})
}

func (m *mqttSubscriber) UpdateSwingModeCommandTopic(ctx context.Context, logger *zerolog.Logger) mqtt.MessageHandler {
	return mqtt.MessageHandler(func(c mqtt.Client, msg mqtt.Message) {
		mac := strings.TrimPrefix(strings.TrimSuffix(msg.Topic(), "/swing_mode/set"), m.mqttConfig.TopicPrefix+"/")

		logger.Debug().Str("device", mac).Str("payload", string(msg.Payload())).Str("topic", msg.Topic()).Msg("new update swing mode message")

		updateSwingModeInput := &models_service.UpdateSwingModeInput{
			Mac:       mac,
			SwingMode: string(msg.Payload()),
		}

		err := m.service.UpdateSwingMode(ctx, logger, updateSwingModeInput)
		if err != nil {
			logger.Error().Err(err).Str("device", mac).Interface("input", updateSwingModeInput).Msg("failed to update swing mode")
			return
		}
	})
}

func (m *mqttSubscriber) UpdateModeCommandTopic(ctx context.Context, logger *zerolog.Logger) mqtt.MessageHandler {
	return mqtt.MessageHandler(func(c mqtt.Client, msg mqtt.Message) {
		mac := strings.TrimPrefix(strings.TrimSuffix(msg.Topic(), "/mode/set"), m.mqttConfig.TopicPrefix+"/")

		logger.Debug().Str("device", mac).Str("payload", string(msg.Payload())).Str("topic", msg.Topic()).Msg("new update mode message")

		updateModeInput := &models_service.UpdateModeInput{
			Mac:  mac,
			Mode: string(msg.Payload()),
		}

		err := m.service.UpdateMode(ctx, logger, updateModeInput)
		if err != nil {
			logger.Error().Err(err).Str("device", mac).Interface("input", updateModeInput).Msg("failed to update mode")
			return
		}
	})
}

func (m *mqttSubscriber) UpdateTemperatureCommandTopic(ctx context.Context, logger *zerolog.Logger) mqtt.MessageHandler {
	return mqtt.MessageHandler(func(c mqtt.Client, msg mqtt.Message) {
		mac := strings.TrimPrefix(strings.TrimSuffix(msg.Topic(), "/temp/set"), m.mqttConfig.TopicPrefix+"/")

		logger.Debug().Str("device", mac).Str("payload", string(msg.Payload())).Str("topic", msg.Topic()).Msg("new update temperature mode message")

		temperature, err := strconv.ParseFloat(string(msg.Payload()), 32)
		if err != nil {
			logger.Error().Err(err).Str("device", mac).Str("command", msg.Topic()).Msg("failed to parse temperature")
			return
		}

		updateTemperatureInput := &models_service.UpdateTemperatureInput{
			Mac:         mac,
			Temperature: float32(temperature),
		}

		err = m.service.UpdateTemperature(ctx, logger, updateTemperatureInput)
		if err != nil {
			logger.Error().Err(err).Str("device", mac).Interface("input", updateTemperatureInput).Msg("failed to update temperature")
			return
		}
	})
}

func (m *mqttSubscriber) GetStatesOnHomeAssistantRestart(ctx context.Context, logger *zerolog.Logger) mqtt.MessageHandler {
	return mqtt.MessageHandler(func(c mqtt.Client, msg mqtt.Message) {

		logger.Debug().Str("payload", string(msg.Payload())).Str("topic", msg.Topic()).Msg("new home assistant LWT message")

		getStatesOnHomeAssistantRestartInput := &models_service.GetStatesOnHomeAssistantRestartInput{
			Status: string(msg.Payload()),
		}

		err := m.service.GetStatesOnHomeAssistantRestart(ctx, logger, getStatesOnHomeAssistantRestartInput)
		if err != nil {
			logger.Error().Err(err).Interface("input", getStatesOnHomeAssistantRestartInput).Msg("failed to get states")
			return
		}
	})
}
