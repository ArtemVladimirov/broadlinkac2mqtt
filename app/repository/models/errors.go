package models

import "errors"

var (
	ErrorDeviceNotFound                   = errors.New("ErrorDeviceNotFound")
	ErrorDeviceAuthNotFound               = errors.New("ErrorDeviceAuthNotFound")
	ErrorDeviceStatusNotFound             = errors.New("ErrorDeviceStatusNotFound")
	ErrorDeviceStatusRawNotFound          = errors.New("ErrorDeviceStatusRawNotFound")
	ErrorDeviceStatusAvailabilityNotFound = errors.New("ErrorDeviceStatusAvailabilityNotFound")

	ErrorDeviceStatusAmbientTempNotFound = errors.New("ErrorDeviceStatusAmbientTempNotFound")
)
