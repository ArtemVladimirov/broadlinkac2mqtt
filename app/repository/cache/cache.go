package cache

import (
	"context"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app/repository/models"
	models_service "github.com/ArtemVladimirov/broadlinkac2mqtt/app/service/models"
	"github.com/rs/zerolog"
	"sync"
)

type cache struct {
	devices map[string]models.Device
	mutex   *sync.Mutex
}

func NewCache() *cache {
	return &cache{
		devices: make(map[string]models.Device),
		mutex:   new(sync.Mutex),
	}
}

func (c *cache) UpsertDeviceConfig(ctx context.Context, logger *zerolog.Logger, input *models.UpsertDeviceConfigInput) error {
	var device models.Device

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device = c.devices[input.Config.Mac]
	device.Config = input.Config
	c.devices[input.Config.Mac] = device
	return nil
}

func (c *cache) ReadDeviceConfig(ctx context.Context, logger *zerolog.Logger, input *models.ReadDeviceConfigInput) (*models.ReadDeviceConfigReturn, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceNotFound
	}

	return &models.ReadDeviceConfigReturn{Config: device.Config}, nil
}

func (c *cache) UpsertDeviceAuth(ctx context.Context, logger *zerolog.Logger, input *models.UpsertDeviceAuthInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return models.ErrorDeviceNotFound
	}

	device.Auth = &input.Auth
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) ReadDeviceAuth(ctx context.Context, logger *zerolog.Logger, input *models.ReadDeviceAuthInput) (*models.ReadDeviceAuthReturn, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceNotFound
	}

	if device.Auth == nil {
		message := "device not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceAuthNotFound
	}

	return &models.ReadDeviceAuthReturn{Auth: *device.Auth}, nil
}

func (c *cache) UpsertAmbientTemp(ctx context.Context, logger *zerolog.Logger, input *models.UpsertAmbientTempInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return models.ErrorDeviceNotFound
	}

	device.DeviceStatus.AmbientTemp = &input.Temperature
	c.devices[input.Mac] = device

	return nil
}

func (c *cache) ReadAmbientTemp(ctx context.Context, logger *zerolog.Logger, input *models.ReadAmbientTempInput) (*models.ReadAmbientTempReturn, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceNotFound
	}

	if device.DeviceStatus.AmbientTemp == nil {
		message := "device status ambient temp is not found in cache"
		logger.Debug().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceStatusAmbientTempNotFound
	}

	return &models.ReadAmbientTempReturn{Temperature: *device.DeviceStatus.AmbientTemp}, nil
}

func (c *cache) UpsertDeviceStatus(ctx context.Context, logger *zerolog.Logger, input *models.UpsertDeviceStatusInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return models.ErrorDeviceNotFound
	}

	device.DeviceStatus.Temperature = &input.Status.Temperature
	device.DeviceStatus.SwingMode = &input.Status.SwingMode
	device.DeviceStatus.Mode = &input.Status.Mode
	device.DeviceStatus.FanMode = &input.Status.FanMode

	c.devices[input.Mac] = device

	return nil
}

func (c *cache) ReadDeviceStatus(ctx context.Context, logger *zerolog.Logger, input *models.ReadDeviceStatusInput) (*models.ReadDeviceStatusReturn, error) {
	var status models_service.DeviceStatusMqtt

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceNotFound
	}

	if device.DeviceStatus.FanMode == nil ||
		device.DeviceStatus.Mode == nil ||
		device.DeviceStatus.SwingMode == nil ||
		device.DeviceStatus.Temperature == nil {
		message := "device status is not found in cache"
		logger.Debug().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceStatusNotFound
	}

	status.FanMode = *device.DeviceStatus.FanMode
	status.SwingMode = *device.DeviceStatus.SwingMode
	status.Temperature = *device.DeviceStatus.Temperature
	status.Mode = *device.DeviceStatus.Mode

	return &models.ReadDeviceStatusReturn{Status: status}, nil
}

func (c *cache) UpsertDeviceStatusRaw(ctx context.Context, logger *zerolog.Logger, input *models.UpsertDeviceStatusRawInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return models.ErrorDeviceNotFound
	}

	device.DeviceStatusRaw = &input.Status
	c.devices[input.Mac] = device

	return nil
}

func (c *cache) ReadDeviceStatusRaw(ctx context.Context, logger *zerolog.Logger, input *models.ReadDeviceStatusRawInput) (*models.ReadDeviceStatusRawReturn, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceNotFound
	}

	if device.DeviceStatusRaw == nil {
		message := "device status raw is not found in cache"
		logger.Debug().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceStatusRawNotFound
	}

	return &models.ReadDeviceStatusRawReturn{Status: *device.DeviceStatusRaw}, nil
}

func (c *cache) UpsertMqttModeMessage(ctx context.Context, logger *zerolog.Logger, input *models.UpsertMqttModeMessageInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.Mode = &input.Mode
	c.devices[input.Mac] = device

	return nil
}

func (c *cache) UpsertMqttSwingModeMessage(ctx context.Context, logger *zerolog.Logger, input *models.UpsertMqttSwingModeMessageInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.SwingMode = &input.SwingMode
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) UpsertMqttFanModeMessage(ctx context.Context, logger *zerolog.Logger, input *models.UpsertMqttFanModeMessageInput) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.FanMode = &input.FanMode
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) UpsertMqttTemperatureMessage(ctx context.Context, logger *zerolog.Logger, input *models.UpsertMqttTemperatureMessageInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return models.ErrorDeviceNotFound
	}

	device.MqttLastMessage.Temperature = &input.Temperature
	c.devices[input.Mac] = device
	return nil
}

func (c *cache) ReadMqttMessage(ctx context.Context, logger *zerolog.Logger, input *models.ReadMqttMessageInput) (*models.ReadMqttMessageReturn, error) {

	var state models.ReadMqttMessageReturn

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceNotFound
	}

	state.Temperature = device.MqttLastMessage.Temperature
	state.Mode = device.MqttLastMessage.Mode
	state.SwingMode = device.MqttLastMessage.SwingMode
	state.FanMode = device.MqttLastMessage.FanMode

	return &state, nil
}

func (c *cache) UpsertDeviceAvailability(ctx context.Context, logger *zerolog.Logger, input *models.UpsertDeviceAvailabilityInput) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return models.ErrorDeviceNotFound
	}

	device.DeviceStatus.Availability = &input.Availability
	c.devices[input.Mac] = device

	return nil
}

func (c *cache) ReadDeviceAvailability(ctx context.Context, logger *zerolog.Logger, input *models.ReadDeviceAvailabilityInput) (*models.ReadDeviceAvailabilityReturn, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	device, ok := c.devices[input.Mac]
	if !ok {
		message := "device is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceNotFound
	}

	if device.DeviceStatus.Availability == nil {
		message := "device status ambient temp is not found in cache"
		logger.Error().Interface("input", input).Msg(message)
		return nil, models.ErrorDeviceStatusAvailabilityNotFound
	}

	return &models.ReadDeviceAvailabilityReturn{Availability: *device.DeviceStatus.Availability}, nil
}

func (c *cache) ReadAuthedDevices(ctx context.Context, logger *zerolog.Logger) (*models.ReadAuthedDevicesReturn, error) {
	var macs []string

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for mac := range c.devices {
		macs = append(macs, mac)
	}

	return &models.ReadAuthedDevicesReturn{Macs: macs}, nil
}
