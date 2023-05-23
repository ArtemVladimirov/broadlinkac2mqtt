package subscriber

import (
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

func Routers(logger *zerolog.Logger, mac string, topicPrefix string, client mqtt.Client, handler app.MqttSubscriber) {

	if token := client.Subscribe(topicPrefix+"/"+mac+"/fan_mode/set", 0, handler.UpdateFanModeCommandTopic(logger)); token.Wait() && token.Error() != nil {
		logger.Error().Err(token.Error())
	}
	if token := client.Subscribe(topicPrefix+"/"+mac+"/swing_mode/set", 0, handler.UpdateSwingModeCommandTopic(logger)); token.Wait() && token.Error() != nil {
		logger.Error().Err(token.Error())
	}
	if token := client.Subscribe(topicPrefix+"/"+mac+"/mode/set", 0, handler.UpdateModeCommandTopic(logger)); token.Wait() && token.Error() != nil {
		logger.Error().Err(token.Error())
	}
	if token := client.Subscribe(topicPrefix+"/"+mac+"/temp/set", 0, handler.UpdateTemperatureCommandTopic(logger)); token.Wait() && token.Error() != nil {
		logger.Error().Err(token.Error())
	}

}
