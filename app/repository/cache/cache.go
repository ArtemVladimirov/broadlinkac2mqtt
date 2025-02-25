package cache

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ArtemVladimirov/broadlinkac2mqtt/app/repository/models"
)

type cache struct {
	devices map[string]models.Device
	mutex   *sync.RWMutex
	logger  *slog.Logger
}

func NewCache(logger *slog.Logger) *cache {
	return &cache{
		devices: make(map[string]models.Device),
		mutex:   new(sync.RWMutex),
		logger:  logger,
	}
}

func (c *cache) UpsertDeviceConfig(ctx context.Context, input *models.UpsertDeviceConfigInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device := c.devices[input.Config.Mac]
	device.Config = input.Config
	c.devices[input.Config.Mac] = device
	return nil
}

func (c *cache) ReadDeviceConfig(ctx context.Context, input *models.ReadDeviceConfigInput) (*models.ReadDeviceConfigReturn, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	return &models.ReadDeviceConfigReturn{Config: device.Config}, nil
}

func (c *cache) UpsertDeviceAuth(ctx context.Context, input *models.UpsertDeviceAuthInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		return models.ErrorDeviceNotFound
	}

	device.Auth = &input.Auth
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) ReadDeviceAuth(ctx context.Context, input *models.ReadDeviceAuthInput) (*models.ReadDeviceAuthReturn, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	if device.Auth == nil {
		return nil, models.ErrorDeviceAuthNotFound
	}

	return &models.ReadDeviceAuthReturn{Auth: *device.Auth}, nil
}

func (c *cache) UpsertAmbientTemp(ctx context.Context, input *models.UpsertAmbientTempInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.DeviceStatus.AmbientTemp = &input.Temperature
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) ReadAmbientTemp(ctx context.Context, input *models.ReadAmbientTempInput) (*models.ReadAmbientTempReturn, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		c.logger.ErrorContext(ctx, message, slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	if device.DeviceStatus.AmbientTemp == nil {
		return nil, models.ErrorDeviceStatusAmbientTempNotFound
	}

	return &models.ReadAmbientTempReturn{Temperature: *device.DeviceStatus.AmbientTemp}, nil
}

func (c *cache) UpsertDeviceStatusRaw(ctx context.Context, input *models.UpsertDeviceStatusRawInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.DeviceStatusRaw = &input.Status
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) ReadDeviceStatusRaw(ctx context.Context, input *models.ReadDeviceStatusRawInput) (*models.ReadDeviceStatusRawReturn, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		c.logger.ErrorContext(ctx, message, slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	if device.DeviceStatusRaw == nil {
		return nil, models.ErrorDeviceStatusRawNotFound
	}

	return &models.ReadDeviceStatusRawReturn{Status: *device.DeviceStatusRaw}, nil
}

func (c *cache) UpsertMqttModeMessage(ctx context.Context, input *models.UpsertMqttModeMessageInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.Mode = &input.Mode
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) UpsertMqttSwingModeMessage(ctx context.Context, input *models.UpsertMqttSwingModeMessageInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.SwingMode = &input.SwingMode
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) UpsertMqttFanModeMessage(ctx context.Context, input *models.UpsertMqttFanModeMessageInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.FanMode = &input.FanMode
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) UpsertMqttTemperatureMessage(ctx context.Context, input *models.UpsertMqttTemperatureMessageInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.Temperature = &input.Temperature
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) ReadMqttMessage(ctx context.Context, input *models.ReadMqttMessageInput) (*models.ReadMqttMessageReturn, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	return &models.ReadMqttMessageReturn{
		Temperature: device.MqttLastMessage.Temperature,
		SwingMode:   device.MqttLastMessage.SwingMode,
		FanMode:     device.MqttLastMessage.FanMode,
		Mode:        device.MqttLastMessage.Mode,
		IsDisplayOn: device.MqttLastMessage.DisplaySwitch,
	}, nil
}

func (c *cache) UpsertDeviceAvailability(ctx context.Context, input *models.UpsertDeviceAvailabilityInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.DeviceStatus.Availability = &input.Availability
	c.devices[input.Mac] = device

	return nil
}

func (c *cache) ReadDeviceAvailability(ctx context.Context, input *models.ReadDeviceAvailabilityInput) (*models.ReadDeviceAvailabilityReturn, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		c.logger.ErrorContext(ctx, message, slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	if device.DeviceStatus.Availability == nil {
		c.logger.ErrorContext(ctx, "device status ambient temp is not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceStatusAvailabilityNotFound
	}

	return &models.ReadDeviceAvailabilityReturn{Availability: *device.DeviceStatus.Availability}, nil
}

func (c *cache) ReadAuthedDevices(ctx context.Context) (*models.ReadAuthedDevicesReturn, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	macs := make([]string, 0, len(c.devices))
	for mac := range c.devices {
		macs = append(macs, mac)
	}

	return &models.ReadAuthedDevicesReturn{Macs: macs}, nil
}

func (c *cache) UpsertMqttDisplaySwitchMessage(ctx context.Context, input *models.UpsertMqttDisplaySwitchMessageInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		c.logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.DisplaySwitch = &input.DisplaySwitch
	c.devices[input.Mac] = device
	return nil
}
