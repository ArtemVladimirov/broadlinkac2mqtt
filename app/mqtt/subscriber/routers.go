package subscriber

import (
	"context"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

func Routers(ctx context.Context, logger *zerolog.Logger, mac string, topicPrefix string, client mqtt.Client, handler app.MqttSubscriber) {

	if token := client.Subscribe(topicPrefix+"/"+mac+"/fan_mode/set", 0, handler.UpdateFanModeCommandTopic(ctx, logger)); token.Wait() && token.Error() != nil {
		logger.Error().Err(token.Error())
	}
	if token := client.Subscribe(topicPrefix+"/"+mac+"/swing_mode/set", 0, handler.UpdateSwingModeCommandTopic(ctx, logger)); token.Wait() && token.Error() != nil {
		logger.Error().Err(token.Error())
	}
	if token := client.Subscribe(topicPrefix+"/"+mac+"/mode/set", 0, handler.UpdateModeCommandTopic(ctx, logger)); token.Wait() && token.Error() != nil {
		logger.Error().Err(token.Error())
	}
	if token := client.Subscribe(topicPrefix+"/"+mac+"/temp/set", 0, handler.UpdateTemperatureCommandTopic(ctx, logger)); token.Wait() && token.Error() != nil {
		logger.Error().Err(token.Error())
	}

}
