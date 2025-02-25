package subscriber

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/ArtemVladimirov/broadlinkac2mqtt/app"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app/mqtt/models"
	modelsservice "github.com/ArtemVladimirov/broadlinkac2mqtt/app/service/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttSubscriber struct {
	logger     *slog.Logger
	mqttConfig models.ConfigMqtt
	service    app.Service
}

func NewMqttReceiver(logger *slog.Logger, service app.Service, mqttConfig models.ConfigMqtt) app.MqttSubscriber {
	return &mqttSubscriber{
		logger:     logger,
		mqttConfig: mqttConfig,
		service:    service,
	}
}

func (m *mqttSubscriber) UpdateFanModeCommandTopic(ctx context.Context) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		mac := strings.TrimPrefix(strings.TrimSuffix(msg.Topic(), "/fan_mode/set"), m.mqttConfig.TopicPrefix+"/")

		m.logger.DebugContext(ctx, "new update fan mode message",
			slog.String("device", mac),
			slog.String("payload", string(msg.Payload())),
			slog.String("topic", msg.Topic()))

		updateFanModeInput := &modelsservice.UpdateFanModeInput{
			Mac:     mac,
			FanMode: string(msg.Payload()),
		}

		err := m.service.UpdateFanMode(ctx, updateFanModeInput)
		if err != nil {
			m.logger.ErrorContext(ctx, "failed to update fan mode", slog.Any("input", updateFanModeInput))
			return
		}
	}
}

func (m *mqttSubscriber) UpdateSwingModeCommandTopic(ctx context.Context) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		mac := strings.TrimPrefix(strings.TrimSuffix(msg.Topic(), "/swing_mode/set"), m.mqttConfig.TopicPrefix+"/")

		m.logger.DebugContext(ctx, "new update swing mode message",
			slog.String("device", mac),
			slog.String("payload", string(msg.Payload())),
			slog.String("topic", msg.Topic()))

		updateSwingModeInput := &modelsservice.UpdateSwingModeInput{
			Mac:       mac,
			SwingMode: string(msg.Payload()),
		}

		err := m.service.UpdateSwingMode(ctx, updateSwingModeInput)
		if err != nil {
			m.logger.ErrorContext(ctx, "failed to update swing mode", slog.Any("input", updateSwingModeInput))
			return
		}
	}
}

func (m *mqttSubscriber) UpdateModeCommandTopic(ctx context.Context) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		mac := strings.TrimPrefix(strings.TrimSuffix(msg.Topic(), "/mode/set"), m.mqttConfig.TopicPrefix+"/")

		m.logger.DebugContext(ctx, "new update mode message",
			slog.String("device", mac),
			slog.String("payload", string(msg.Payload())),
			slog.String("topic", msg.Topic()))

		updateModeInput := &modelsservice.UpdateModeInput{
			Mac:  mac,
			Mode: string(msg.Payload()),
		}

		err := m.service.UpdateMode(ctx, updateModeInput)
		if err != nil {
			m.logger.ErrorContext(ctx, "failed to update mode", slog.Any("input", updateModeInput))
			return
		}
	}
}

func (m *mqttSubscriber) UpdateTemperatureCommandTopic(ctx context.Context) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		mac := strings.TrimPrefix(strings.TrimSuffix(msg.Topic(), "/temp/set"), m.mqttConfig.TopicPrefix+"/")

		m.logger.DebugContext(ctx, "new update temperature mode message",
			slog.String("device", mac),
			slog.String("payload", string(msg.Payload())),
			slog.String("topic", msg.Topic()))

		temperature, err := strconv.ParseFloat(string(msg.Payload()), 32)
		if err != nil {
			m.logger.ErrorContext(ctx, "failed to parse temperature", slog.Any("err", err), slog.String("input", string(msg.Payload())))
			return
		}

		updateTemperatureInput := &modelsservice.UpdateTemperatureInput{
			Mac:         mac,
			Temperature: float32(temperature),
		}

		err = m.service.UpdateTemperature(ctx, updateTemperatureInput)
		if err != nil {
			m.logger.ErrorContext(ctx, "failed to update temperature", slog.Any("input", updateTemperatureInput))
			return
		}
	}
}

func (m *mqttSubscriber) GetStatesOnHomeAssistantRestart(ctx context.Context) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		m.logger.DebugContext(ctx, "new home assistant LWT message",
			slog.String("payload", string(msg.Payload())),
			slog.String("topic", msg.Topic()))

		getStatesOnHomeAssistantRestartInput := &modelsservice.PublishStatesOnHomeAssistantRestartInput{
			Status: string(msg.Payload()),
		}

		err := m.service.PublishStatesOnHomeAssistantRestart(ctx, getStatesOnHomeAssistantRestartInput)
		if err != nil {
			m.logger.ErrorContext(ctx, "failed to get states", slog.Any("input", getStatesOnHomeAssistantRestartInput))
			return
		}
	}
}

func (m *mqttSubscriber) UpdateDisplaySwitchCommandTopic(ctx context.Context) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		mac := strings.TrimPrefix(strings.TrimSuffix(msg.Topic(), "/display/switch/set"), m.mqttConfig.TopicPrefix+"/")

		m.logger.DebugContext(ctx, "new update display status message",
			slog.String("device", mac),
			slog.String("payload", string(msg.Payload())),
			slog.String("topic", msg.Topic()))

		updateDisplaySwitchInput := &modelsservice.UpdateDisplaySwitchInput{
			Mac:    mac,
			Status: string(msg.Payload()),
		}

		err := m.service.UpdateDisplaySwitch(ctx, updateDisplaySwitchInput)
		if err != nil {
			m.logger.ErrorContext(ctx, "failed to update display switch", slog.Any("input", updateDisplaySwitchInput))
			return
		}
	}
}
