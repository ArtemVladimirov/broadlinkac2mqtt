package cache

import (
	"context"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app/repository/models"
	"log/slog"
	"sync"
)

type cache struct {
	devices map[string]models.Device
	mutex   *sync.RWMutex
}

func NewCache() *cache {
	return &cache{
		devices: make(map[string]models.Device),
		mutex:   new(sync.RWMutex),
	}
}

func (c *cache) UpsertDeviceConfig(ctx context.Context, logger *slog.Logger, input *models.UpsertDeviceConfigInput) error {
	var device models.Device

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device = c.devices[input.Config.Mac]
	device.Config = input.Config
	c.devices[input.Config.Mac] = device
	return nil
}

func (c *cache) ReadDeviceConfig(ctx context.Context, logger *slog.Logger, input *models.ReadDeviceConfigInput) (*models.ReadDeviceConfigReturn, error) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	return &models.ReadDeviceConfigReturn{Config: device.Config}, nil
}

func (c *cache) UpsertDeviceAuth(ctx context.Context, logger *slog.Logger, input *models.UpsertDeviceAuthInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.Auth = &input.Auth
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) ReadDeviceAuth(ctx context.Context, logger *slog.Logger, input *models.ReadDeviceAuthInput) (*models.ReadDeviceAuthReturn, error) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	if device.Auth == nil {
		logger.ErrorContext(ctx, "device not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceAuthNotFound
	}

	return &models.ReadDeviceAuthReturn{Auth: *device.Auth}, nil
}

func (c *cache) UpsertAmbientTemp(ctx context.Context, logger *slog.Logger, input *models.UpsertAmbientTempInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.DeviceStatus.AmbientTemp = &input.Temperature
	c.devices[input.Mac] = device

	return nil
}

func (c *cache) ReadAmbientTemp(ctx context.Context, logger *slog.Logger, input *models.ReadAmbientTempInput) (*models.ReadAmbientTempReturn, error) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.ErrorContext(ctx, message, slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	if device.DeviceStatus.AmbientTemp == nil {
		logger.ErrorContext(ctx, "device status ambient temp is not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceStatusAmbientTempNotFound
	}

	return &models.ReadAmbientTempReturn{Temperature: *device.DeviceStatus.AmbientTemp}, nil
}

func (c *cache) UpsertDeviceStatusRaw(ctx context.Context, logger *slog.Logger, input *models.UpsertDeviceStatusRawInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.DeviceStatusRaw = &input.Status
	c.devices[input.Mac] = device

	return nil
}

func (c *cache) ReadDeviceStatusRaw(ctx context.Context, logger *slog.Logger, input *models.ReadDeviceStatusRawInput) (*models.ReadDeviceStatusRawReturn, error) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.ErrorContext(ctx, message, slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	if device.DeviceStatusRaw == nil {
		logger.ErrorContext(ctx, "device status raw is not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceStatusRawNotFound
	}

	return &models.ReadDeviceStatusRawReturn{Status: *device.DeviceStatusRaw}, nil
}

func (c *cache) UpsertMqttModeMessage(ctx context.Context, logger *slog.Logger, input *models.UpsertMqttModeMessageInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.Mode = &input.Mode
	c.devices[input.Mac] = device

	return nil
}

func (c *cache) UpsertMqttSwingModeMessage(ctx context.Context, logger *slog.Logger, input *models.UpsertMqttSwingModeMessageInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.SwingMode = &input.SwingMode
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) UpsertMqttFanModeMessage(ctx context.Context, logger *slog.Logger, input *models.UpsertMqttFanModeMessageInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.FanMode = &input.FanMode
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) UpsertMqttTemperatureMessage(ctx context.Context, logger *slog.Logger, input *models.UpsertMqttTemperatureMessageInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.Temperature = &input.Temperature
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) ReadMqttMessage(ctx context.Context, logger *slog.Logger, input *models.ReadMqttMessageInput) (*models.ReadMqttMessageReturn, error) {

	var state models.ReadMqttMessageReturn

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	state.Temperature = device.MqttLastMessage.Temperature
	state.Mode = device.MqttLastMessage.Mode
	state.SwingMode = device.MqttLastMessage.SwingMode
	state.FanMode = device.MqttLastMessage.FanMode
	state.IsDisplayOn = device.MqttLastMessage.DisplaySwitch

	return &state, nil
}

func (c *cache) UpsertDeviceAvailability(ctx context.Context, logger *slog.Logger, input *models.UpsertDeviceAvailabilityInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.DeviceStatus.Availability = &input.Availability
	c.devices[input.Mac] = device

	return nil
}

func (c *cache) ReadDeviceAvailability(ctx context.Context, logger *slog.Logger, input *models.ReadDeviceAvailabilityInput) (*models.ReadDeviceAvailabilityReturn, error) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.ErrorContext(ctx, message, slog.Any("input", input))
		return nil, models.ErrorDeviceNotFound
	}

	if device.DeviceStatus.Availability == nil {
		logger.ErrorContext(ctx, "device status ambient temp is not found in cache", slog.Any("input", input))
		return nil, models.ErrorDeviceStatusAvailabilityNotFound
	}

	return &models.ReadDeviceAvailabilityReturn{Availability: *device.DeviceStatus.Availability}, nil
}

func (c *cache) ReadAuthedDevices(ctx context.Context, logger *slog.Logger) (*models.ReadAuthedDevicesReturn, error) {
	var macs []string

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for mac := range c.devices {
		macs = append(macs, mac)
	}

	return &models.ReadAuthedDevicesReturn{Macs: macs}, nil
}

func (c *cache) UpsertMqttDisplaySwitchMessage(ctx context.Context, logger *slog.Logger, input *models.UpsertMqttDisplaySwitchMessageInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		logger.ErrorContext(ctx, "device is not found in cache", slog.Any("input", input))
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.DisplaySwitch = &input.DisplaySwitch
	c.devices[input.Mac] = device

	return nil
}
