package main

import (
	"context"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app/mqtt"
	workspaceMqttModels "github.com/ArtemVladimirov/broadlinkac2mqtt/app/mqtt/models"
	workspaceMqttSender "github.com/ArtemVladimirov/broadlinkac2mqtt/app/mqtt/publisher"
	workspaceMqttReceiver "github.com/ArtemVladimirov/broadlinkac2mqtt/app/mqtt/subscriber"
	workspaceCache "github.com/ArtemVladimirov/broadlinkac2mqtt/app/repository/cache"
	workspaceService "github.com/ArtemVladimirov/broadlinkac2mqtt/app/service"
	workspaceServiceModels "github.com/ArtemVladimirov/broadlinkac2mqtt/app/service/models"
	workspaceWebClient "github.com/ArtemVladimirov/broadlinkac2mqtt/app/webClient"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/config"
	paho "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/sync/errgroup"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type App struct {
	devices             []workspaceServiceModels.DeviceConfig
	autoDiscoveryTopic  *string
	topicPrefix         string
	logLevel            string
	wsBroadLinkReceiver app.WebClient
	wsMqttReceiver      app.MqttSubscriber
	wsService           app.Service
	client              paho.Client
}

func NewApp(logger *slog.Logger) (*App, error) {
	// Configuration
	cfg, err := config.NewConfig(logger)
	if err != nil {
		return nil, err
	}

	// MQTT
	mqttConfig := workspaceMqttModels.ConfigMqtt{
		Broker:                   cfg.Mqtt.Broker,
		User:                     cfg.Mqtt.User,
		Password:                 cfg.Mqtt.Password,
		ClientId:                 cfg.Mqtt.ClientId,
		TopicPrefix:              cfg.Mqtt.TopicPrefix,
		AutoDiscoveryTopic:       cfg.Mqtt.AutoDiscoveryTopic,
		AutoDiscoveryTopicRetain: cfg.Mqtt.AutoDiscoveryTopicRetain,
	}

	opts, _ := mqtt.NewMqttConfig(logger, cfg.Mqtt)
	client := paho.NewClient(opts)

	//Configure MQTT Sender Layer
	mqttSender := workspaceMqttSender.NewMqttSender(
		mqttConfig,
		client,
	)

	//Configure Service Layer
	service := workspaceService.NewService(
		cfg.Mqtt.TopicPrefix,
		cfg.Service.UpdateInterval,
		mqttSender,
		workspaceWebClient.NewWebClient(),
		workspaceCache.NewCache(),
	)
	//Configure MQTT Receiver Layer
	mqttReceiver := workspaceMqttReceiver.NewMqttReceiver(
		service,
		mqttConfig,
	)

	devices := make([]workspaceServiceModels.DeviceConfig, 0, len(cfg.Devices))
	for _, device := range cfg.Devices {
		if len(device.TemperatureUnit) == 0 {
			device.TemperatureUnit = "C"
		}

		dev := workspaceServiceModels.DeviceConfig{
			Ip:              device.Ip,
			Mac:             strings.ToLower(device.Mac),
			Name:            device.Name,
			Port:            device.Port,
			TemperatureUnit: strings.ToUpper(device.TemperatureUnit),
		}

		err = dev.Validate()
		if err != nil {
			logger.Error("device config is incorrect", slog.String("device", device.Mac), slog.Any("err", err))
			return nil, err
		}

		devices = append(devices, dev)
	}

	application := &App{
		wsMqttReceiver:     mqttReceiver,
		client:             client,
		devices:            devices,
		wsService:          service,
		topicPrefix:        cfg.Mqtt.TopicPrefix,
		autoDiscoveryTopic: cfg.Mqtt.AutoDiscoveryTopic,
		logLevel:           cfg.Service.LogLevel,
	}

	return application, nil
}

func (app *App) Run(ctx context.Context, logger *slog.Logger) error {
	// Run MQTT
	if token := app.client.Connect(); token.Wait() && token.Error() != nil {
		err := token.Error()
		if err != nil {
			logger.ErrorContext(ctx, "failed to connect mqtt",
				slog.Any("err", err))
			return err
		}
	}

	if app.autoDiscoveryTopic != nil {
		if token := app.client.Subscribe(*app.autoDiscoveryTopic+"/status", 0, app.wsMqttReceiver.GetStatesOnHomeAssistantRestart(ctx, logger)); token.Wait() && token.Error() != nil {
			err := token.Error()
			if err != nil {
				logger.ErrorContext(ctx, "failed to subscribe on LWT",
					slog.Any("err", err))

				return err
			}
		}
	}

	// Create Device
	for _, device := range app.devices {
		err := app.wsService.CreateDevice(ctx, logger, &workspaceServiceModels.CreateDeviceInput{
			Config: workspaceServiceModels.DeviceConfig{
				Mac:  device.Mac,
				Ip:   device.Ip,
				Name: device.Name,
				Port: device.Port,
			}})
		if err != nil {
			logger.ErrorContext(ctx, "failed to create the device",
				slog.Any("err", err))
			return err
		}
	}

	for _, device := range app.devices {
		device := device
		go func() {
			for {
				err := app.wsService.AuthDevice(ctx, logger, &workspaceServiceModels.AuthDeviceInput{Mac: device.Mac})
				if err == nil {
					break
				}
				logger.ErrorContext(ctx, "failed to Auth device "+device.Mac+". Reconnect in 3 seconds...",
					slog.Any("err", err))
				time.Sleep(time.Second * 3)
			}

			// Subscribe on MQTT handlers
			workspaceMqttReceiver.Routers(ctx, logger, device.Mac, app.topicPrefix, app.client, app.wsMqttReceiver)

			//Publish Discovery Topic
			if app.autoDiscoveryTopic != nil {
				err := app.wsService.PublishDiscoveryTopic(ctx, logger, &workspaceServiceModels.PublishDiscoveryTopicInput{Device: device})
				if err != nil {
					return
				}
			}

			err := app.wsService.StartDeviceMonitoring(ctx, logger, &workspaceServiceModels.StartDeviceMonitoringInput{Mac: device.Mac})
			if err != nil {
				return
			}
		}()
	}

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	killSignal := <-interrupt
	switch killSignal {
	case syscall.SIGKILL:
		logger.Info("Got SIGKILL...")
	case syscall.SIGQUIT:
		logger.Info("Got SIGQUIT...")
	case syscall.SIGTERM:
		logger.Info("Got SIGTERM...")
	case syscall.SIGINT:
		logger.Info("Got SIGINT...")
	default:
		logger.Info("Undefined killSignal...")
	}
	// Publish offline states for devices
	g := new(errgroup.Group)
	for _, device := range app.devices {
		device := device
		g.Go(func() error {
			err := app.wsService.UpdateDeviceAvailability(ctx, logger, &workspaceServiceModels.UpdateDeviceAvailabilityInput{
				Mac:          device.Mac,
				Availability: "offline",
			})
			if err != nil {
				logger.ErrorContext(ctx, "failed to update availability",
					slog.String("device", device.Mac),
					slog.Any("err", err))
				return err
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	// Disconnect MQTT
	app.client.Disconnect(100)

	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logLevel := &slog.LevelVar{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))

	application, err := NewApp(logger)
	if err != nil {
		logger.ErrorContext(ctx, "failed to get a new App", slog.Any("err", err))
		return
	}

	switch application.logLevel {
	case "error":
		logLevel.Set(slog.LevelError)
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "disabled":
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	case "info":
		logLevel.Set(slog.LevelInfo)
	default:
		logLevel.Set(slog.LevelError)
	}

	// Run
	err = application.Run(ctx, logger)
	if err != nil {
		logger.ErrorContext(ctx, "failed to get a new App", slog.Any("err", err))
		return
	}
}
