package subscriber

import (
	"context"
	"log/slog"

	"github.com/ArtemVladimirov/broadlinkac2mqtt/app"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Routers(ctx context.Context, logger *slog.Logger, mac string, topicPrefix string, client mqtt.Client, handler app.MqttSubscriber) {
	prefix := topicPrefix + "/" + mac

	if token := client.Subscribe(prefix+"/fan_mode/set", 0, handler.UpdateFanModeCommandTopic(ctx, logger)); token.Wait() && token.Error() != nil {
		logger.ErrorContext(ctx, "failed to subscribe on topic", slog.Any("err", token.Error()))
	}
	if token := client.Subscribe(prefix+"/swing_mode/set", 0, handler.UpdateSwingModeCommandTopic(ctx, logger)); token.Wait() && token.Error() != nil {
		logger.ErrorContext(ctx, "failed to subscribe on topic", slog.Any("err", token.Error()))
	}
	if token := client.Subscribe(prefix+"/mode/set", 0, handler.UpdateModeCommandTopic(ctx, logger)); token.Wait() && token.Error() != nil {
		logger.ErrorContext(ctx, "failed to subscribe on topic", slog.Any("err", token.Error()))
	}
	if token := client.Subscribe(prefix+"/temp/set", 0, handler.UpdateTemperatureCommandTopic(ctx, logger)); token.Wait() && token.Error() != nil {
		logger.ErrorContext(ctx, "failed to subscribe on topic", slog.Any("err", token.Error()))
	}
	if token := client.Subscribe(prefix+"/display/switch/set", 0, handler.UpdateDisplaySwitchCommandTopic(ctx, logger)); token.Wait() && token.Error() != nil {
		logger.ErrorContext(ctx, "failed to subscribe on topic", slog.Any("err", token.Error()))
	}
}
