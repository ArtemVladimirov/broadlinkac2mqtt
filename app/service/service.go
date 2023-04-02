package service

import (
	"context"
	"encoding/hex"
	"github.com/ArtVladimirov/BroadlinkAC2Mqtt/app"
	models_mqtt "github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/mqtt/models"
	models_repo "github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/repository/models"
	"github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/service/models"
	models_web "github.com/ArtVladimirov/BroadlinkAC2Mqtt/app/webClient/models"
	"github.com/ArtVladimirov/BroadlinkAC2Mqtt/pkg/coder"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type service struct {
	updateInterval int
	topicPrefix    string
	mqtt           app.MqttPublisher
	webClient      app.WebClient
	cache          app.Cache
}

func NewService(topicPrefix string, updateInterval int, mqtt app.MqttPublisher, webClient app.WebClient, cache app.Cache) *service {
	return &service{
		topicPrefix:    topicPrefix,
		updateInterval: updateInterval,
		mqtt:           mqtt,
		webClient:      webClient,
		cache:          cache,
	}
}

func (s *service) CreateDevice(ctx context.Context, logger *zerolog.Logger, input *models.CreateDeviceInput) (*models.CreateDeviceReturn, error) {
	rand.Seed(time.Now().UnixNano())

	key := []byte{0x09, 0x76, 0x28, 0x34, 0x3f, 0xe9, 0x9e, 0x23, 0x76, 0x5c, 0x15, 0x13, 0xac, 0xcf, 0x8b, 0x02}
	iv := []byte{0x56, 0x2e, 0x17, 0x99, 0x6d, 0x09, 0x3d, 0x28, 0xdd, 0xb3, 0xba, 0x69, 0x5a, 0x2e, 0x6f, 0x58}

	// Store device information in the repository
	upsertDeviceConfigInput := &models_repo.UpsertDeviceConfigInput{
		Config: input.Config,
	}
	err := s.cache.UpsertDeviceConfig(ctx, logger, upsertDeviceConfigInput)
	if err != nil {
		return nil, err
	}

	auth := models.DeviceAuth{
		LastMessageId: rand.Intn(0xffff),
		DevType:       0x4E2a,
		Id:            [4]byte{0, 0, 0, 0},
		Key:           key,
		Iv:            iv,
	}

	// Store device information in the repository
	upsertDeviceAuthInput := &models_repo.UpsertDeviceAuthInput{
		Mac:  input.Config.Mac,
		Auth: auth,
	}
	err = s.cache.UpsertDeviceAuth(ctx, logger, upsertDeviceAuthInput)
	if err != nil {
		return nil, err
	}

	return &models.CreateDeviceReturn{}, nil
}

/*
AuthDevice

	Request

0000   34 ea 34 da da c8 e0 d5 5e 68 9e 3e 08 00 45 00
0010   00 a4 16 d4 00 00 80 11 00 00 c0 a8 01 24 c0 a8
0020   01 13 f9 a1 00 50 00 90 84 29 5a a5 aa 55 5a a5
0030   aa 55 00 00 00 00 00 00 00 00 00 00 00 00 00 00
0040   00 00 00 00 00 00 00 00 00 00 de f0 00 00 2a 4e
0050   65 00 63 f7 34 ea 34 da da c8 00 00 00 00 a1 c3
0060   00 00 45 34 52 e7 f9 2e da 95 83 44 93 08 35 ef
0070   9a 6d fb 69 2d c3 70 b9 04 43 ac 5c d6 3f bb 53
0080   ad fa 08 81 4c a7 f8 cf 41 71 00 32 8e 57 0c 3b
0090   86 c9 4d 05 70 84 49 a3 89 e2 9a e1 04 54 36 a0
00a0   5b dd dc 02 c1 61 af 13 25 e8 7e 19 b0 f7 d1 ce
00b0   06 8d

	Response

0000   e0 d5 5e 68 9e 3e 34 ea 34 da da c8 08 00 45 00
0010   00 74 56 1e 00 00 40 11 a0 d3 c0 a8 01 13 c0 a8
0020   01 24 00 50 f9 a1 00 60 18 82 5a a5 aa 55 5a a5
0030   aa 55 00 00 00 00 00 00 00 00 00 00 00 00 00 00
0040   00 00 00 00 00 00 00 00 00 00 28 dc 00 00 2a 4e
0050   e9 03 63 f7 34 ea 34 da da c8 00 00 00 00 c1 c7
0060   00 00 bb 6c bb bb 34 58 5c d4 42 b9 cf bb db 30
0070   3e ea 55 af e0 62 cd d6 38 16 4b 81 cc 38 40 84
0080   ef 9e
*/
func (s *service) AuthDevice(ctx context.Context, logger *zerolog.Logger, input *models.AuthDeviceInput) error {

	payload := [0x50]byte{}
	payload[0x04] = 0x31
	payload[0x05] = 0x31
	payload[0x06] = 0x31
	payload[0x07] = 0x31
	payload[0x08] = 0x31
	payload[0x09] = 0x31
	payload[0x0a] = 0x31
	payload[0x0b] = 0x31
	payload[0x0c] = 0x31
	payload[0x0d] = 0x31
	payload[0x0e] = 0x31
	payload[0x0f] = 0x31
	payload[0x10] = 0x31
	payload[0x11] = 0x31
	payload[0x12] = 0x31
	payload[0x1e] = 0x01
	payload[0x2d] = 0x01
	payload[0x30] = byte('T')
	payload[0x31] = byte('e')
	payload[0x32] = byte('s')
	payload[0x33] = byte('t')
	payload[0x34] = byte(' ')
	payload[0x35] = byte(' ')
	payload[0x36] = byte('1')

	sendCommandInput := &models.SendCommandInput{
		Command: 0x65,
		Payload: payload[:],
		Mac:     input.Mac,
	}

	response, err := s.SendCommand(ctx, logger, sendCommandInput)
	if err != nil {
		logger.Error().Err(err).Interface("input", input).Msg("failed to send command")
		return err
	}

	// Decode message
	if len(response.Payload) >= 0x38 {
		response.Payload = response.Payload[0x38:]
	} else {
		logger.Error().Interface("input", input).Msg("response is too short")
		return err
	}

	// Read the saved value in repo if no
	readDeviceAuthInput := &models_repo.ReadDeviceAuthInput{
		Mac: input.Mac,
	}
	readDeviceAuthReturn, err := s.cache.ReadDeviceAuth(ctx, logger, readDeviceAuthInput)
	if err != nil {
		logger.Error().Interface("input", input).Msg("device not found")
		return err
	}
	auth := readDeviceAuthReturn.Auth
	response.Payload, _ = coder.Decrypt(auth.Key, auth.Iv, response.Payload)

	auth = models.DeviceAuth{
		LastMessageId: auth.LastMessageId,
		DevType:       auth.DevType,
		Id:            [4]byte{response.Payload[0], response.Payload[1], response.Payload[2], response.Payload[3]},
		Key:           response.Payload[0x04:0x14],
		Iv:            auth.Iv,
	}

	// Update the last message id in the cache
	upsertDeviceAuthInput := &models_repo.UpsertDeviceAuthInput{
		Mac:  input.Mac,
		Auth: auth,
	}
	err = s.cache.UpsertDeviceAuth(ctx, logger, upsertDeviceAuthInput)
	if err != nil {
		return err
	}

	return nil
}

/*
GetDeviceAmbientTemperature

	Request

0000   34 ea 34 da da c8 e0 d5 5e 68 9e 3e 08 00 45 00
0010   00 64 16 d6 00 00 80 11 00 00 c0 a8 01 24 c0 a8
0020   01 13 e0 dc 00 50 00 50 83 e9// 5a a5 aa 55 5a a5
0030   aa 55 00 00 00 00 00 00 00 00 00 00 00 00 00 00
0040   00 00 00 00 00 00 00 00 00 00 a1 d1 00 00 2a 4e
0050   6a 00 90 7c 34 ea 34 da da c8 01 00 00 00 b9 c0
0060   00 00 3d 19 77 32 16 2c b4 f5 f9 e1 8a ca 7b 1b
0070   ff 13

	Response

0000   e0 d5 5e 68 9e 3e 34 ea 34 da da c8 08 00 45 00
0010   00 84 56 ab 00 00 40 11 a0 36 c0 a8 01 13 c0 a8
0020   01 24 00 50 e0 dc 00 70 08 12 5a a5 aa 55 5a a5
0030   aa 55 00 00 00 00 00 00 00 00 00 00 00 00 00 00
0040   00 00 00 00 00 00 00 00 00 00 40 e3 00 00 2a 4e
0050   ee 03 90 7c 34 ea 34 da da c8 01 00 00 00 cf c1
0060   00 00 2c 4f a6 c5 65 f7 8b 46 82 92 20 a3 6f bf
0070   65 24 a6 8a 04 97 eb 37 ef e6 a6 42 2a 4f 6b 8a
0080   ed 81 d1 67 c3 8d b2 69 c5 0a e4 e2 91 05 bc 52
0090   5e 60
*/
func (s *service) GetDeviceAmbientTemperature(ctx context.Context, logger *zerolog.Logger, input *models.GetDeviceAmbientTemperatureInput) error {
	//

	payload, err := hex.DecodeString("0C00BB0006800000020021011B7E0000")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to decode string")
	}

	sendCommandInput := &models.SendCommandInput{
		Command: 0x6a,
		Payload: payload,
		Mac:     input.Mac,
	}

	response, err := s.SendCommand(ctx, logger, sendCommandInput)
	if err != nil {
		logger.Error().Err(err).Interface("input", sendCommandInput).Msg("Failed to send a command")
		return err
	}

	if (response.Payload[0x22] | (response.Payload[0x23] << 8)) != 0 {
		logger.Error().Err(err).Interface("input", sendCommandInput).Msg("Checksum is incorrect")
		return models.ErrorInvalidResultPacket
	}

	// Decode message
	if len(response.Payload) >= 0x38 {
		response.Payload = response.Payload[0x38:]
	} else {
		logger.Error().Interface("input", input).Msg("response is too short")
		return models.ErrorInvalidResultPacketLength
	}

	// Read the saved value in repo if no
	readDeviceAuthInput := &models_repo.ReadDeviceAuthInput{
		Mac: input.Mac,
	}
	readDeviceAuthReturn, err := s.cache.ReadDeviceAuth(ctx, logger, readDeviceAuthInput)
	if err != nil {
		logger.Error().Interface("input", input).Msg("device not found")
		return err
	}

	auth := readDeviceAuthReturn.Auth

	response.Payload, _ = coder.Decrypt(auth.Key, auth.Iv, response.Payload)

	//Drop leading stuff as don't need
	response.Payload = response.Payload[2:]

	if len(response.Payload) < 40 {
		return models.ErrorInvalidResultPacketLength
	}

	ambientTemp := response.Payload[15] & 0b00011111

	readAmbientTempInput := &models_repo.ReadAmbientTempInput{Mac: input.Mac}

	readAmbientTempReturn, err := s.cache.ReadAmbientTemp(ctx, logger, readAmbientTempInput)
	if err != nil {
		logger.Error().Interface("input", readAmbientTempInput).Str("device", input.Mac).Msg("failed to read the ambient temperature")
		return err
	}

	logger.Debug().Int8("ambientTemp", int8(ambientTemp)).Str("device", input.Mac).Msg("Ambient temperature")

	if readAmbientTempReturn == nil || readAmbientTempReturn.Temperature != int8(ambientTemp) {
		// Sent  temperature to MQTT
		publishAmbientTempInput := &models_mqtt.PublishAmbientTempInput{
			Mac:         input.Mac,
			Temperature: int8(ambientTemp),
		}

		err = s.mqtt.PublishAmbientTemp(ctx, logger, publishAmbientTempInput)
		if err != nil {
			logger.Error().Interface("input", publishAmbientTempInput).Msg("failed to publish ambient temperature")
			return err
		}

		// Save the new value in storage
		upsertAmbientTempInput := &models_repo.UpsertAmbientTempInput{Temperature: int8(ambientTemp), Mac: input.Mac}

		err = s.cache.UpsertAmbientTemp(ctx, logger, upsertAmbientTempInput)
		if err != nil {
			logger.Error().Interface("input", upsertAmbientTempInput).Msg("failed to upsert the temperature")
			return err
		}
	}

	return nil

}

func (s *service) GetDeviceStates(ctx context.Context, logger *zerolog.Logger, input *models.GetDeviceStatesInput) error {
	////////////////////////////////////////////////////////////
	//              SEND COMMAND TO GET STATES                //
	////////////////////////////////////////////////////////////

	payload, err := hex.DecodeString("0C00BB0006800000020011012B7E0000")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to decode string")
	}

	sendCommandInput := &models.SendCommandInput{
		Command: 0x6a,
		Payload: payload,
		Mac:     input.Mac,
	}

	response, err := s.SendCommand(ctx, logger, sendCommandInput)
	if err != nil {
		logger.Error().Err(err).Interface("input", sendCommandInput).Msg("Failed to send the command to get states")
		return err
	}

	////////////////////////////////////////////////////////////
	//                 DECODE RESPONSE                        //
	////////////////////////////////////////////////////////////

	if (response.Payload[0x22] | (response.Payload[0x23] << 8)) != 0 {
		logger.Error().Err(err).Interface("input", sendCommandInput).Msg("Checksum is incorrect")
		return models.ErrorInvalidResultPacket
	}

	// Read the saved value in repo if no
	readDeviceAuthInput := &models_repo.ReadDeviceAuthInput{
		Mac: input.Mac,
	}
	readDeviceAuthReturn, err := s.cache.ReadDeviceAuth(ctx, logger, readDeviceAuthInput)
	if err != nil {
		logger.Error().Interface("input", input).Msg("device not found")
		return err
	}
	auth := readDeviceAuthReturn.Auth

	// Decode message
	if len(response.Payload) >= 0x38 {
		response.Payload = response.Payload[0x38:]
	} else {
		logger.Error().Interface("input", input).Msg("response is too short")
		return models.ErrorInvalidResultPacketLength
	}

	response.Payload, err = coder.Decrypt(auth.Key, auth.Iv, response.Payload)
	if err != nil {
		logger.Error().Err(err).Interface("input", response.Payload).Msg("Failed to decrypt the response")
		return err
	}

	if response.Payload[4] != 0x07 {
		logger.Error().Err(err).Interface("input", sendCommandInput).Bytes("payload", response.Payload).Msg("It is not a result packet")
		return models.ErrorInvalidResultPacket
	}

	if response.Payload[0] != 0x19 {
		logger.Error().Err(err).Interface("input", sendCommandInput).Msg("The length of the packet is incorrect. Must be 25")
		return models.ErrorInvalidResultPacketLength
	}

	//Drop leading stuff as don't need
	response.Payload = response.Payload[2:]

	var raw = models.DeviceStatusRaw{
		UpdatedAt:          time.Now(),
		Temperature:        float32(8 + (response.Payload[10] >> 3) + byte(0.5*float32(response.Payload[12]>>7))),
		Power:              response.Payload[18] >> 5 & 0b00000001,
		FixationVertical:   response.Payload[10] & 0b00000111,
		Mode:               response.Payload[15] >> 5 & 0b00001111,
		Sleep:              response.Payload[15] >> 2 & 0b00000001,
		Display:            response.Payload[20] >> 4 & 0b00000001,
		Mildew:             response.Payload[20] >> 3 & 0b00000001,
		Health:             response.Payload[18] >> 1 & 0b00000001,
		FixationHorizontal: response.Payload[10] & 0b00000111,
		FanSpeed:           response.Payload[13] >> 5 & 0b00000111,
		IFeel:              response.Payload[15] >> 3 & 0b00000001,
		Mute:               response.Payload[14] >> 7 & 0b00000001,
		Turbo:              response.Payload[14] >> 6 & 0b00000001,
		Clean:              response.Payload[18] >> 2 & 0b00000001,
	}

	//////////////////////////////////////////////////////////////////
	// Comparison of Home Assistant statuses and BroadLink statuses //
	//////////////////////////////////////////////////////////////////
	var deviceStatusMqtt models.DeviceStatusMqtt

	// Temperature
	deviceStatusMqtt.Temperature = raw.Temperature

	// Modes
	//Modes: []string{"auto", "off", "cool", "heat", "dry", "fan_only"},
	if int(raw.Power) == onOffStatusesInvert["OFF"] {
		deviceStatusMqtt.Mode = "off"
	} else {
		status, ok := modeStatuses[int(raw.Mode)]
		if ok {
			deviceStatusMqtt.Mode = status
		} else {
			deviceStatusMqtt.Mode = "error"
		}
	}

	// Fan Status
	//FanModes:  "auto", "low", "medium", "high", "turbo", "mute"
	fanStatus, ok := fanStatuses[int(raw.FanSpeed)]
	if ok {
		deviceStatusMqtt.FanMode = fanStatus
	} else {
		deviceStatusMqtt.FanMode = "error"
	}

	if int(raw.Mute) == onOffStatusesInvert["ON"] {
		deviceStatusMqtt.FanMode = "mute"
	}

	if int(raw.Turbo) == onOffStatusesInvert["ON"] {
		deviceStatusMqtt.FanMode = "turbo"
	}

	// Swing Modes
	verticalFixationStatus, ok := verticalFixationStatuses[int(raw.FixationVertical)]
	if ok {
		deviceStatusMqtt.SwingMode = verticalFixationStatus
	} else {
		deviceStatusMqtt.SwingMode = ""
	}

	logger.Debug().Interface("status", deviceStatusMqtt).Str("device", input.Mac).Msg("The converted current device status")

	//////////////////////////////////////////////////////////////////
	//  Compare new statuses with old statuses and update  MQTT     //
	//////////////////////////////////////////////////////////////////
	readDeviceStatusInput := &models_repo.ReadDeviceStatusInput{
		Mac: input.Mac,
	}

	readDeviceStatusReturn, err := s.cache.ReadDeviceStatus(ctx, logger, readDeviceStatusInput)
	if err != nil {
		logger.Error().Err(err).Interface("input", readDeviceStatusInput).Msg("Failed to read the device status")
		return err
	}

	var updateTemperature, updateMode, updateSwingMode, updateFanMode bool

	if readDeviceStatusReturn == nil {
		updateTemperature, updateMode, updateSwingMode, updateFanMode = true, true, true, true
	} else {
		if deviceStatusMqtt.SwingMode != readDeviceStatusReturn.Status.SwingMode {
			updateSwingMode = true
		}
		if deviceStatusMqtt.Mode != readDeviceStatusReturn.Status.Mode {
			updateMode = true
		}
		if deviceStatusMqtt.FanMode != readDeviceStatusReturn.Status.FanMode {
			updateFanMode = true
		}
		if deviceStatusMqtt.Temperature != readDeviceStatusReturn.Status.Temperature {
			updateTemperature = true
		}
	}

	g := new(errgroup.Group)
	g.Go(func() error {
		if updateTemperature {
			publishTemperatureInput := &models_mqtt.PublishTemperatureInput{
				Mac:         input.Mac,
				Temperature: deviceStatusMqtt.Temperature,
			}

			err = s.mqtt.PublishTemperature(ctx, logger, publishTemperatureInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishTemperatureInput).Msg("Failed to publish the device set temperature")
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		if updateMode {
			publishModeInput := &models_mqtt.PublishModeInput{
				Mac:  input.Mac,
				Mode: deviceStatusMqtt.Mode,
			}

			err = s.mqtt.PublishMode(ctx, logger, publishModeInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishModeInput).Msg("Failed to publish the device mode")
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		if updateFanMode {
			publishFanModeInput := &models_mqtt.PublishFanModeInput{
				Mac:     input.Mac,
				FanMode: deviceStatusMqtt.FanMode,
			}

			err = s.mqtt.PublishFanMode(ctx, logger, publishFanModeInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishFanModeInput).Msg("Failed to publish the device fan mode")
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		if updateSwingMode {
			publishSwingModeInput := &models_mqtt.PublishSwingModeInput{
				Mac:       input.Mac,
				SwingMode: deviceStatusMqtt.SwingMode,
			}

			err = s.mqtt.PublishSwingMode(ctx, logger, publishSwingModeInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishSwingModeInput).Msg("Failed to publish the device swing mode")
				return err
			}
		}
		return nil
	})

	// Wait for all HTTP fetches to complete.
	if err = g.Wait(); err != nil {
		return err
	}

	//////////////////////////////////////////////////////////////////
	//  		Update device states in the database                //
	//////////////////////////////////////////////////////////////////

	if updateTemperature || updateMode || updateSwingMode || updateFanMode {
		upsertDeviceStatusInput := &models_repo.UpsertDeviceStatusInput{
			Mac:    input.Mac,
			Status: deviceStatusMqtt,
		}

		err = s.cache.UpsertDeviceStatus(ctx, logger, upsertDeviceStatusInput)
		if err != nil {
			logger.Error().Err(err).Interface("input", upsertDeviceStatusInput).Msg("Failed to upsert device status")
			return err
		}
	}

	upsertDeviceStatusRawInput := &models_repo.UpsertDeviceStatusRawInput{
		Mac:    input.Mac,
		Status: raw,
	}

	err = s.cache.UpsertDeviceStatusRaw(ctx, logger, upsertDeviceStatusRawInput)
	if err != nil {
		logger.Error().Err(err).Interface("input", upsertDeviceStatusRawInput).Msg("Failed to upsert the raw device status")
		return err
	}
	return nil
}

func (s *service) SendCommand(ctx context.Context, logger *zerolog.Logger, input *models.SendCommandInput) (*models.SendCommandReturn, error) {

	// Read the saved value in repo if no
	readDeviceAuthInput := &models_repo.ReadDeviceAuthInput{
		Mac: input.Mac,
	}
	readDeviceAuthReturn, err := s.cache.ReadDeviceAuth(ctx, logger, readDeviceAuthInput)
	if err != nil {
		logger.Error().Interface("input", input).Msg("device not found")
		return nil, err
	}

	auth := readDeviceAuthReturn.Auth

	auth.LastMessageId = (auth.LastMessageId + 1) & 0xffff

	var macByteSlice []byte

	for i := 0; i < len(input.Mac); i = i + 2 {
		val, err := strconv.ParseUint(input.Mac[i:i+2], 16, 8)
		if err != nil {
			logger.Error().Err(err).Interface("input", input).Msg("Mac address is not correct")
			return nil, err
		}
		macByteSlice = append(macByteSlice, byte(val))
	}

	var packet [0x38]byte

	// Insert body
	packet[0x00] = 0x5a
	packet[0x01] = 0xa5
	packet[0x02] = 0xaa
	packet[0x03] = 0x55
	packet[0x04] = 0x5a
	packet[0x05] = 0xa5
	packet[0x06] = 0xaa
	packet[0x07] = 0x55
	packet[0x24] = 0x2a
	packet[0x25] = 0x4e
	packet[0x26] = input.Command // command
	packet[0x28] = byte(auth.LastMessageId & 0xff)
	packet[0x29] = byte(auth.LastMessageId >> 8)
	packet[0x2a] = macByteSlice[0]
	packet[0x2b] = macByteSlice[1]
	packet[0x2c] = macByteSlice[2]
	packet[0x2d] = macByteSlice[3]
	packet[0x2e] = macByteSlice[4]
	packet[0x2f] = macByteSlice[5]
	packet[0x30] = auth.Id[0]
	packet[0x31] = auth.Id[1]
	packet[0x32] = auth.Id[2]
	packet[0x33] = auth.Id[3]

	checksum := 0xbeaf
	for i := range input.Payload {
		checksum += int(input.Payload[i])
		checksum = checksum & 0xffff
	}

	input.Payload, err = coder.Encrypt(auth.Key, auth.Iv, input.Payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to encrypt payload")
		return nil, err
	}

	packet[0x34] = byte(checksum & 0xff)
	packet[0x35] = byte(checksum >> 8)

	var packetSlice = packet[:]

	packetSlice = append(packetSlice, input.Payload...)

	// Create and insert Checksum
	checksum = 0xbeaf
	for i := range packetSlice {
		checksum += int(packetSlice[i])
		checksum = checksum & 0xffff
	}
	packetSlice[0x20] = byte(checksum & 0xff)
	packetSlice[0x21] = byte(checksum >> 8)

	// Update last message id in database
	upsertDeviceAuthInput := &models_repo.UpsertDeviceAuthInput{
		Mac:  input.Mac,
		Auth: auth,
	}
	err = s.cache.UpsertDeviceAuth(ctx, logger, upsertDeviceAuthInput)
	if err != nil {
		return nil, err
	}

	logger.Debug().Bytes("packet", packetSlice).Str("device", input.Mac)

	// Read config to get IP and Port
	readDeviceConfigInput := &models_repo.ReadDeviceConfigInput{
		Mac: input.Mac,
	}

	readDeviceConfigReturn, err := s.cache.ReadDeviceConfig(ctx, logger, readDeviceConfigInput)
	if err != nil {
		logger.Error().Interface("input", input).Msg("failed to read device config")
		return nil, err
	}

	// Send the packet
	sendCommandInput := &models_web.SendCommandInput{
		Payload: packetSlice,
		Ip:      readDeviceConfigReturn.Config.Ip,
		Port:    readDeviceConfigReturn.Config.Port,
	}

	sendCommandReturn, err := s.webClient.SendCommand(ctx, logger, sendCommandInput)
	if err != nil {
		logger.Error().Interface("input", input).Str("device", input.Mac).Msg("failed to send a command")
		return nil, err
	}

	return &models.SendCommandReturn{Payload: sendCommandReturn.Payload}, nil
}

func (s *service) PublishDiscoveryTopic(ctx context.Context, logger *zerolog.Logger, input *models.PublishDiscoveryTopicInput) error {

	prefix := s.topicPrefix + "/" + input.Device.Mac
	publishDiscoveryTopicInput := models_mqtt.PublishDiscoveryTopicInput{
		DiscoveryTopic: models_mqtt.DiscoveryTopic{
			FanModeCommandTopic:     prefix + "/fan_mode/set",
			FanModes:                []string{"auto", "low", "medium", "high", "turbo", "mute"},
			FanModeStateTopic:       prefix + "/fan_mode/value",
			ModeCommandTopic:        prefix + "/mode/set",
			ModeStateTopic:          prefix + "/mode/value",
			Modes:                   []string{"auto", "off", "cool", "heat", "dry", "fan_only"},
			SwingModeCommandTopic:   prefix + "/swing_mode/set",
			SwingModeStateTopic:     prefix + "/swing_mode/value",
			SwingModes:              []string{"top", "middle1", "middle2", "middle3", "bottom", "swing", "auto"},
			MinTemp:                 16,
			MaxTemp:                 32,
			TempStep:                0.5,
			TemperatureStateTopic:   prefix + "/temp/value",
			TemperatureCommandTopic: prefix + "/temp/set",
			Precision:               0.5,
			Device: models_mqtt.DiscoveryTopicDevice{
				Model: "AirCon",
				Mf:    "ArtVladimirov",
				Sw:    "v1.0.0",
				Ids:   input.Device.Mac,
				Name:  input.Device.Name,
			},
			UniqueId: input.Device.Mac,
			Availability: models_mqtt.DiscoveryTopicAvailability{
				PayloadAvailable:    "online",
				PayloadNotAvailable: "offline",
				Topic:               prefix + "/availability/value",
			},
			CurrentTemperatureTopic: prefix + "/current_temp/value",
			Name:                    input.Device.Name,
		},
	}
	err := s.mqtt.PublishDiscoveryTopic(ctx, logger, publishDiscoveryTopicInput)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateFanMode(ctx context.Context, logger *zerolog.Logger, input *models.UpdateFanModeInput) error {

	upsertMqttFanModeMessageInput := &models_repo.UpsertMqttFanModeMessageInput{
		Mac: input.Mac,
		FanMode: models_repo.MqttFanModeMessage{
			UpdatedAt: time.Now(),
			FanMode:   input.FanMode,
		},
	}

	err := s.cache.UpsertMqttFanModeMessage(ctx, logger, upsertMqttFanModeMessageInput)
	if err != nil {
		logger.Error().Interface("input", upsertMqttFanModeMessageInput).Str("device", input.Mac).Msg("failed to save mqtt message to cache storage")
		return err
	}

	return nil
}

func (s *service) UpdateMode(ctx context.Context, logger *zerolog.Logger, input *models.UpdateModeInput) error {

	upsertMqttModeMessageInput := &models_repo.UpsertMqttModeMessageInput{
		Mac: input.Mac,
		Mode: models_repo.MqttModeMessage{
			UpdatedAt: time.Now(),
			Mode:      input.Mode,
		},
	}

	err := s.cache.UpsertMqttModeMessage(ctx, logger, upsertMqttModeMessageInput)
	if err != nil {
		logger.Error().Interface("input", upsertMqttModeMessageInput).Str("device", input.Mac).Msg("failed to save mqtt message to cache storage")
		return err
	}

	return nil
}

func (s *service) UpdateSwingMode(ctx context.Context, logger *zerolog.Logger, input *models.UpdateSwingModeInput) error {

	upsertMqttSwingModeMessageInput := &models_repo.UpsertMqttSwingModeMessageInput{
		Mac: input.Mac,
		SwingMode: models_repo.MqttSwingModeMessage{
			UpdatedAt: time.Now(),
			SwingMode: input.SwingMode,
		},
	}

	err := s.cache.UpsertMqttSwingModeMessage(ctx, logger, upsertMqttSwingModeMessageInput)
	if err != nil {
		logger.Error().Interface("input", upsertMqttSwingModeMessageInput).Str("device", input.Mac).Msg("failed to save mqtt message to cache storage")
		return err
	}

	return nil
}

func (s *service) UpdateTemperature(ctx context.Context, logger *zerolog.Logger, input *models.UpdateTemperatureInput) error {

	upsertMqttTemperatureMessageInput := &models_repo.UpsertMqttTemperatureMessageInput{
		Mac: input.Mac,
		Temperature: models_repo.MqttTemperatureMessage{
			UpdatedAt:   time.Now(),
			Temperature: input.Temperature,
		},
	}

	err := s.cache.UpsertMqttTemperatureMessage(ctx, logger, upsertMqttTemperatureMessageInput)
	if err != nil {
		logger.Error().Interface("input", upsertMqttTemperatureMessageInput).Str("device", input.Mac).Msg("failed to save mqtt message to cache storage")
		return err
	}

	return nil
}

func (s *service) UpdateDeviceStates(ctx context.Context, logger *zerolog.Logger, input *models.UpdateDeviceStatesInput) error {

	readDeviceStatusRawInput := &models_repo.ReadDeviceStatusRawInput{
		Mac: input.Mac,
	}

	readDeviceStatusRawReturn, err := s.cache.ReadDeviceStatusRaw(ctx, logger, readDeviceStatusRawInput)
	if err != nil {
		logger.Error().Interface("input", input).Str("device", input.Mac).Msg("failed to read device raw status")
		return err
	}

	// Convert Home Assistant to BroadLink types
	var verticalFixation byte
	if input.SwingMode != nil {
		key, ok := verticalFixationStatusesInvert[*input.SwingMode]
		if !ok {
			logger.Error().Interface("input", input).Str("device", input.Mac).
				Str("swingMode", *input.SwingMode).
				Err(models.ErrorInvalidParameterSwingMode).
				Msg("Invalid parameter Swing mode")

			return models.ErrorInvalidParameterSwingMode
		} else {
			verticalFixation = byte(key)
		}
	} else {
		verticalFixation = readDeviceStatusRawReturn.Status.FixationVertical
	}

	var temperature, temperature05 int
	if input.Temperature != nil {
		if *input.Temperature > 32 || *input.Temperature < 16 {

			logger.Error().Interface("input", input).Str("device", input.Mac).
				Float32("temperature", *input.Temperature).
				Err(models.ErrorInvalidParameterTemperature).
				Msg("Invalid parameter temperature")

			return models.ErrorInvalidParameterTemperature
		}

		temperature = int(*input.Temperature) - 8

		if int(*input.Temperature*10)%(int(*input.Temperature)*10) == 5 {
			temperature05 = 1
		}

	} else {
		if readDeviceStatusRawReturn.Status.Temperature < 16 {
			temperature = 16 - 8
		} else if readDeviceStatusRawReturn.Status.Temperature > 32 {
			temperature = 32 - 8
		} else {
			temperature = int(readDeviceStatusRawReturn.Status.Temperature) - 8
			if readDeviceStatusRawReturn.Status.Temperature-float32(int(readDeviceStatusRawReturn.Status.Temperature)) != 0 {
				temperature05 = 1
			}
		}
	}
	var fanMode, turbo, mute byte
	if input.FanMode != nil {
		if *input.FanMode == "mute" {
			mute = byte(onOffStatusesInvert["ON"])
		} else if *input.FanMode == "turbo" {
			turbo = byte(onOffStatusesInvert["ON"])
		} else {
			key, ok := fanStatusesInvert[*input.FanMode]
			if !ok {
				logger.Error().Interface("input", input).Str("device", input.Mac).
					Str("fanMode", *input.FanMode).
					Err(models.ErrorInvalidParameterFanMode).
					Msg("Invalid parameter fan mode")

				return models.ErrorInvalidParameterFanMode
			} else {
				fanMode = byte(key)
				turbo = byte(onOffStatusesInvert["OFF"])
				mute = byte(onOffStatusesInvert["OFF"])
			}
		}
	} else {
		fanMode = readDeviceStatusRawReturn.Status.FanSpeed
		mute = readDeviceStatusRawReturn.Status.Mute
		turbo = readDeviceStatusRawReturn.Status.Turbo
	}

	// Mode
	var mode, power byte

	if input.Mode != nil {
		switch strings.ToLower(*input.Mode) {
		case "cool":
			mode = byte(modeStatusesInvert["cool"])
			power = byte(onOffStatusesInvert["ON"])
		case "heat":
			mode = byte(modeStatusesInvert["heat"])
			power = byte(onOffStatusesInvert["ON"])
		case "auto":
			mode = byte(modeStatusesInvert["auto"])
			power = byte(onOffStatusesInvert["ON"])
		case "dry":
			mode = byte(modeStatusesInvert["dry"])
			power = byte(onOffStatusesInvert["ON"])
		case "fan_only":
			mode = byte(modeStatusesInvert["fan_only"])
			power = byte(onOffStatusesInvert["ON"])
		case "off":
			power = byte(onOffStatusesInvert["OFF"])
		default:
			logger.Error().Interface("input", input).Str("device", input.Mac).
				Str("mode", *input.Mode).
				Err(models.ErrorInvalidParameterMode).
				Msg("Invalid parameter mode")
		}
	} else {
		power = readDeviceStatusRawReturn.Status.Power
		mode = readDeviceStatusRawReturn.Status.Mode
	}

	// Insert values in payload
	var payload [23]byte
	payload[0] = 0xbb
	payload[1] = 0x00
	payload[2] = 0x06
	payload[3] = 0x80
	payload[4] = 0x00
	payload[5] = 0x00
	payload[6] = 0x0f
	payload[7] = 0x00
	payload[8] = 0x01
	payload[9] = 0x01
	payload[10] = 0b00000000 | byte(temperature)<<3 | verticalFixation
	payload[11] = 0b00000000 | readDeviceStatusRawReturn.Status.FixationHorizontal<<5
	payload[12] = 0b00001111 | byte(temperature05)<<7
	payload[13] = 0b00000000 | fanMode<<5
	payload[14] = 0b00000000 | turbo<<6 | mute<<7
	payload[15] = 0b00000000 | mode<<5 | readDeviceStatusRawReturn.Status.Sleep<<2
	payload[16] = 0b00000000
	payload[17] = 0x00
	payload[18] = 0b00000000 | power<<5 | readDeviceStatusRawReturn.Status.Health<<1 | readDeviceStatusRawReturn.Status.Clean<<2
	payload[19] = 0x00
	payload[20] = 0b00000000 | readDeviceStatusRawReturn.Status.Display<<4 | readDeviceStatusRawReturn.Status.Mildew<<3
	payload[21] = 0b00000000
	payload[22] = 0b00000000

	// Additional preparations
	var requestPayload [32]byte
	requestPayload[0] = byte(len(payload) + 2)

	copy(requestPayload[2:], payload[:])

	var checksum int
	for i := 0; i < len(payload); i += 2 {
		checksum += int(payload[i])<<8 + int(append(payload[:], byte(0))[i+1])
	}
	checksum = (checksum >> 16) + (checksum & 0xFFFF)
	checksum = ^checksum & 0xFFFF

	requestPayload[len(payload)+2] = byte((checksum >> 8) & 0xFF)
	requestPayload[len(payload)+3] = byte(checksum & 0xFF)

	// Send command
	sendCommandInput := &models.SendCommandInput{
		Command: 0x6a,
		Payload: requestPayload[:],
		Mac:     input.Mac,
	}

	_, err = s.SendCommand(ctx, logger, sendCommandInput)
	if err != nil {
		logger.Error().Err(err).Interface("input", input).Str("device", input.Mac).Msg("failed to send a set command")
		return err
	}

	return nil
}

func (s *service) UpdateDeviceAvailability(ctx context.Context, logger *zerolog.Logger, input *models.UpdateDeviceAvailabilityInput) error {

	publishAvailabilityInput := &models_mqtt.PublishAvailabilityInput{
		Mac:          input.Mac,
		Availability: input.Availability,
	}

	err := s.mqtt.PublishAvailability(ctx, logger, publishAvailabilityInput)
	if err != nil {
		logger.Error().Interface("input", publishAvailabilityInput).Str("device", input.Mac).Msg("failed to create command payload")
		return err
	}

	return nil
}

func (s *service) StartDeviceMonitoring(ctx context.Context, logger *zerolog.Logger, input *models.StartDeviceMonitoringInput) error {

	// Update ambient temperature once in 3 minutes
	go func() {
		for {
			err := s.GetDeviceAmbientTemperature(ctx, logger, &models.GetDeviceAmbientTemperatureInput{Mac: input.Mac})
			if err != nil {
				logger.Error().Str("device", input.Mac).Msg("failed to get ambient temperature")
			}
			time.Sleep(time.Minute * 3)
		}
	}()

	var (
		modeUpdatedTime, swingModeUpdatedTime, fanModeUpdatedTime, temperatureUpdatedTime time.Time
		lastUpdate                                                                        time.Time
		isDeviceAvailable                                                                 bool
	)

	for {
		var (
			updateDeviceState        = false
			mode, swingMode, fanMode *string
			temperature              *float32
		)

		readMqttMessageInput := &models_repo.ReadMqttMessageInput{
			Mac: input.Mac,
		}
		message, err := s.cache.ReadMqttMessage(ctx, logger, readMqttMessageInput)
		if err != nil {
			return err
		}

		if message.Mode != nil {
			if message.Mode.UpdatedAt != modeUpdatedTime {
				updateDeviceState = true
				mode = &message.Mode.Mode
				modeUpdatedTime = message.Mode.UpdatedAt

				publishModeInput := &models_mqtt.PublishModeInput{
					Mac:  input.Mac,
					Mode: *mode,
				}
				err := s.mqtt.PublishMode(ctx, logger, publishModeInput)
				if err != nil {
					logger.Error().Interface("input", publishModeInput).Str("device", input.Mac).Msg("failed to publish mode to mqtt")
				}
			}
		}

		if message.FanMode != nil {
			if message.FanMode.UpdatedAt != fanModeUpdatedTime {
				updateDeviceState = true
				fanMode = &message.FanMode.FanMode
				fanModeUpdatedTime = message.FanMode.UpdatedAt

				publishFanModeInput := &models_mqtt.PublishFanModeInput{
					Mac:     input.Mac,
					FanMode: *fanMode,
				}
				err := s.mqtt.PublishFanMode(ctx, logger, publishFanModeInput)
				if err != nil {
					logger.Error().Interface("input", publishFanModeInput).Str("device", input.Mac).Msg("failed to publish fan mode to mqtt")
				}
			}
		}

		if message.SwingMode != nil {
			if message.SwingMode.UpdatedAt != swingModeUpdatedTime {
				updateDeviceState = true
				swingMode = &message.SwingMode.SwingMode
				swingModeUpdatedTime = message.FanMode.UpdatedAt

				publishSwingModeInput := &models_mqtt.PublishSwingModeInput{
					Mac:       input.Mac,
					SwingMode: *swingMode,
				}
				err := s.mqtt.PublishSwingMode(ctx, logger, publishSwingModeInput)
				if err != nil {
					logger.Error().Interface("input", publishSwingModeInput).Str("device", input.Mac).Msg("failed to publish swing mode to mqtt")
				}
			}
		}

		if message.Temperature != nil {
			if message.Temperature.UpdatedAt != temperatureUpdatedTime {
				updateDeviceState = true
				temperature = &message.Temperature.Temperature
				temperatureUpdatedTime = message.Temperature.UpdatedAt

				publishTemperatureModeInput := &models_mqtt.PublishTemperatureInput{
					Mac:         input.Mac,
					Temperature: *temperature,
				}
				err := s.mqtt.PublishTemperature(ctx, logger, publishTemperatureModeInput)
				if err != nil {
					logger.Error().Interface("input", publishTemperatureModeInput).Str("device", input.Mac).Msg("failed to publish temperature to mqtt")
				}
			}
		}

		if (updateDeviceState && time.Now().Sub(lastUpdate).Seconds() > 2) || int(time.Now().Sub(lastUpdate).Seconds()) > s.updateInterval {
			for {
				err = s.GetDeviceStates(ctx, logger, &models.GetDeviceStatesInput{Mac: input.Mac})
				if err != nil {
					logger.Error().Err(err).Interface("device", input.Mac).Msg("Failed to get AC States")
					if int(time.Now().Sub(lastUpdate).Seconds()) > s.updateInterval*3 && isDeviceAvailable {
						isDeviceAvailable = false
						updateDeviceAvailabilityInput := &models.UpdateDeviceAvailabilityInput{
							Mac:          input.Mac,
							Availability: "offline",
						}
						err = s.UpdateDeviceAvailability(ctx, logger, updateDeviceAvailabilityInput)
						if err != nil {
							logger.Error().Err(err).Str("device", input.Mac).Interface("input", updateDeviceAvailabilityInput).Msg("Failed to update device availability")
						}
					}
				} else {
					lastUpdate = time.Now()
					if !isDeviceAvailable {
						isDeviceAvailable = true
						updateDeviceAvailabilityInput := &models.UpdateDeviceAvailabilityInput{
							Mac:          input.Mac,
							Availability: "online",
						}
						err = s.UpdateDeviceAvailability(ctx, logger, updateDeviceAvailabilityInput)
						if err != nil {
							logger.Error().Err(err).Str("device", input.Mac).Interface("input", updateDeviceAvailabilityInput).Msg("Failed to update device availability")
						}
					}
					time.Sleep(time.Millisecond * 500)
					break
				}
			}
		}

		if updateDeviceState && isDeviceAvailable {
			updateDeviceStatesInput := &models.UpdateDeviceStatesInput{
				Mac:         input.Mac,
				FanMode:     fanMode,
				SwingMode:   swingMode,
				Mode:        mode,
				Temperature: temperature,
			}

			err := s.UpdateDeviceStates(ctx, logger, updateDeviceStatesInput)
			if err != nil {
				logger.Error().Err(err).Str("device", input.Mac).Interface("input", updateDeviceStatesInput).Msg("Failed to update device states")
			} else {
				lastUpdate = time.UnixMicro(0)
			}
		}

		time.Sleep(time.Millisecond * 500)
	}
}
