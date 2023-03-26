package publisher

import (
	"context"
	"encoding/json"
	"github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/mqtt/models"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"strconv"
)

const (
	deviceClass string = "climate"
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

func (m *mqttPublisher) PublishDiscoveryTopic(ctx context.Context, logger *zerolog.Logger, input models.PublishDiscoveryTopicInput) error {

	if m.mqttConfig.AutoDiscovery == false {
		return nil
	}

	payload, err := json.Marshal(input.DiscoveryTopic)
	if err != nil {
		logger.Error().Err(err).Interface("input", input.DiscoveryTopic).Msg("Failed to marshal discovery topic")
	}

	topic := m.mqttConfig.AutoDiscoveryTopic + "/" + deviceClass + "/" + input.DiscoveryTopic.UniqueId + "/config"
	m.client.Publish(topic, 0, m.mqttConfig.AutoDiscoveryTopicRetain, string(payload))

	return nil
}

func (m *mqttPublisher) PublishAmbientTemp(ctx context.Context, logger *zerolog.Logger, input *models.PublishAmbientTempInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/temp/value"

	m.client.Publish(topic, 0, false, strconv.Itoa(int(input.Temperature)))
	return nil
}

func (m *mqttPublisher) PublishTemperature(ctx context.Context, logger *zerolog.Logger, input *models.PublishTemperatureInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/current_temp/value"

	m.client.Publish(topic, 0, false, strconv.Itoa(int(input.Temperature)))
	return nil
}

func (m *mqttPublisher) PublishMode(ctx context.Context, logger *zerolog.Logger, input *models.PublishModeInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/mode/value"

	m.client.Publish(topic, 0, false, input.Mode)
	return nil
}

func (m *mqttPublisher) PublishSwingMode(ctx context.Context, logger *zerolog.Logger, input *models.PublishSwingModeInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/swing_mode/value"

	m.client.Publish(topic, 0, false, input.SwingMode)
	return nil
}

func (m *mqttPublisher) PublishFanMode(ctx context.Context, logger *zerolog.Logger, input *models.PublishFanModeInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/fan_mode/value"

	m.client.Publish(topic, 0, false, input.FanMode)
	return nil
}

func (m *mqttPublisher) PublishAvailability(ctx context.Context, logger *zerolog.Logger, input *models.PublishAvailabilityInput) error {

	topic := m.mqttConfig.TopicPrefix + "/" + input.Mac + "/availability/value"

	m.client.Publish(topic, 0, false, input.Availability)
	return nil
}
