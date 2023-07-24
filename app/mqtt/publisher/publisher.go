package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app/mqtt/models"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

type mqttPublisher struct {
	mqttConfig models.ConfigMqtt
	client     paho.Client
}

func NewMqttSender(mqttConfig models.ConfigMqtt, client paho.Client) *mqttPublisher {
	return &mqttPublisher{
		mqttConfig: mqttConfig,
		client:     client,
	}
}

func (m *mqttPublisher) PublishClimateDiscoveryTopic(ctx context.Context, logger *zerolog.Logger, input models.PublishClimateDiscoveryTopicInput) error {

	if m.mqttConfig.AutoDiscoveryTopic == nil {
		return nil
	}

	payload, err := json.Marshal(input.Topic)
	if err != nil {
		logger.Error().Err(err).Interface("input", input.Topic).Msg("Failed to marshal discovery topic")
		return err
	}

	topic := *m.mqttConfig.AutoDiscoveryTopic + "/" + models.DeviceClassClimate + "/" + input.Topic.UniqueId + "/config"

	token := m.client.Publish(topic, 0, m.mqttConfig.AutoDiscoveryTopicRetain, string(payload))
	select {
	case <-ctx.Done():
		return nil
	case <-token.Done():
		return nil
	}
}

func (m *mqttPublisher) PublishSwitchDiscoveryTopic(ctx context.Context, logger *zerolog.Logger, input models.PublishSwitchDiscoveryTopicInput) error {

	if m.mqttConfig.AutoDiscoveryTopic == nil {
		return nil
	}

	payload, err := json.Marshal(input.Topic)
	if err != nil {
		logger.Error().Err(err).Interface("input", input.Topic).Msg("Failed to marshal discovery topic")
		return err
	}

	topic := *m.mqttConfig.AutoDiscoveryTopic + "/" + models.DeviceClassSwitch + "/" + input.Topic.UniqueId + "/config"

	token := m.client.Publish(topic, 0, m.mqttConfig.AutoDiscoveryTopicRetain, string(payload))
	select {
	case <-ctx.Done():
		return nil
	case <-token.Done():
		return nil
	}
}

func (m *mqttPublisher) PublishAmbientTemp(ctx context.Context, logger *zerolog.Logger, input *models.PublishAmbientTempInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/current_temp/value"

	token := m.client.Publish(topic, 0, false, fmt.Sprintf("%.1f", input.Temperature))
	select {
	case <-ctx.Done():
		return nil
	case <-token.Done():
		return nil
	}
}

func (m *mqttPublisher) PublishTemperature(ctx context.Context, logger *zerolog.Logger, input *models.PublishTemperatureInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/temp/value"

	token := m.client.Publish(topic, 0, false, fmt.Sprintf("%.1f", input.Temperature))
	select {
	case <-ctx.Done():
		return nil
	case <-token.Done():
		return nil
	}
}

func (m *mqttPublisher) PublishMode(ctx context.Context, logger *zerolog.Logger, input *models.PublishModeInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/mode/value"

	token := m.client.Publish(topic, 0, false, input.Mode)
	select {
	case <-ctx.Done():
		return nil
	case <-token.Done():
		return nil
	}
}

func (m *mqttPublisher) PublishSwingMode(ctx context.Context, logger *zerolog.Logger, input *models.PublishSwingModeInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/swing_mode/value"

	token := m.client.Publish(topic, 0, false, input.SwingMode)
	select {
	case <-ctx.Done():
		return nil
	case <-token.Done():
		return nil
	}
}

func (m *mqttPublisher) PublishFanMode(ctx context.Context, logger *zerolog.Logger, input *models.PublishFanModeInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/fan_mode/value"

	token := m.client.Publish(topic, 0, false, input.FanMode)
	select {
	case <-ctx.Done():
		return nil
	case <-token.Done():
		return nil
	}
}

func (m *mqttPublisher) PublishAvailability(ctx context.Context, logger *zerolog.Logger, input *models.PublishAvailabilityInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/availability/value"

	token := m.client.Publish(topic, 0, false, input.Availability)
	select {
	case <-ctx.Done():
		return nil
	case <-token.Done():
		return nil
	}
}

func (m *mqttPublisher) PublishDisplaySwitch(ctx context.Context, logger *zerolog.Logger, input *models.PublishDisplaySwitchInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/display/switch/value"

	token := m.client.Publish(topic, 0, false, input.Status)
	select {
	case <-ctx.Done():
		return nil
	case <-token.Done():
		return nil
	}
}
