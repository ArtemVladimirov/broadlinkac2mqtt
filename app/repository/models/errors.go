package models

import "errors"

var (
	ErrorDeviceNotFound                   = errors.New("ErrorDeviceNotFound")
	ErrorDeviceAuthNotFound               = errors.New("ErrorDeviceAuthNotFound")
	ErrorDeviceStatusRawNotFound          = errors.New("ErrorDeviceStatusRawNotFound")
	ErrorDeviceStatusAvailabilityNotFound = errors.New("ErrorDeviceStatusAvailabilityNotFound")

	ErrorDeviceStatusAmbientTempNotFound = errors.New("ErrorDeviceStatusAmbientTempNotFound")
)
