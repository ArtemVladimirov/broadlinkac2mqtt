package models

import "errors"

var (
	ErrorInvalidResultPacket       = errors.New("ErrorInvalidResultPacket")
	ErrorInvalidResultPacketLength = errors.New("ErrorInvalidResultPacketLength")

	ErrorInvalidParameterTemperature   = errors.New("ErrorInvalidParameterTemperature")
	ErrorInvalidParameterSwingMode     = errors.New("ErrorInvalidParameterSwingMode")
	ErrorInvalidParameterFanMode       = errors.New("ErrorInvalidParameterFanMode")
	ErrorInvalidParameterMode          = errors.New("ErrorInvalidParameterMode")
	ErrorInvalidParameterDisplayStatus = errors.New("ErrorInvalidParameterDisplayStatus")
)
