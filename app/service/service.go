package service

import (
	"context"
	"encoding/hex"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app"
	models_mqtt "github.com/ArtemVladimirov/broadlinkac2mqtt/app/mqtt/models"
	models_repo "github.com/ArtemVladimirov/broadlinkac2mqtt/app/repository/models"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app/service/models"
	models_web "github.com/ArtemVladimirov/broadlinkac2mqtt/app/webClient/models"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/pkg/coder"
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

func (s *service) CreateDevice(ctx context.Context, logger *zerolog.Logger, input *models.CreateDeviceInput) error {
	rand.Seed(time.Now().UnixNano())

	key := []byte{0x09, 0x76, 0x28, 0x34, 0x3f, 0xe9, 0x9e, 0x23, 0x76, 0x5c, 0x15, 0x13, 0xac, 0xcf, 0x8b, 0x02}
	iv := []byte{0x56, 0x2e, 0x17, 0x99, 0x6d, 0x09, 0x3d, 0x28, 0xdd, 0xb3, 0xba, 0x69, 0x5a, 0x2e, 0x6f, 0x58}

	// Store device information in the repository
	upsertDeviceConfigInput := &models_repo.UpsertDeviceConfigInput{
		Config: input.Config,
	}
	err := s.cache.UpsertDeviceConfig(ctx, logger, upsertDeviceConfigInput)
	if err != nil {
		return err
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
		return err
	}

	return nil
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

	if uint16(response.Payload[0x22])|(uint16(response.Payload[0x23])<<8) != 0 {
		logger.Error().Err(models.ErrorInvalidResultPacket).Interface("input", sendCommandInput).Msg("Checksum is incorrect")
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

	response.Payload, err = coder.Decrypt(auth.Key, auth.Iv, response.Payload)
	if err != nil {
		logger.Error().Interface("input", input).Msg("failed to decrypt payload")
		return err
	}

	//Drop leading stuff as don't need
	response.Payload = response.Payload[2:]

	if len(response.Payload) < 40 {
		return models.ErrorInvalidResultPacketLength
	}

	ambientTemp := float32(response.Payload[15]-0b00100000) + (float32(response.Payload[31]) / 10)

	readAmbientTempInput := &models_repo.ReadAmbientTempInput{Mac: input.Mac}

	readAmbientTempReturn, err := s.cache.ReadAmbientTemp(ctx, logger, readAmbientTempInput)
	if err != nil {
		switch err {
		case models_repo.ErrorDeviceStatusAmbientTempNotFound:
			err = nil
		default:
			logger.Error().Interface("input", readAmbientTempInput).Str("device", input.Mac).Msg("failed to read the ambient temperature")
			return err
		}
	}

	if readAmbientTempReturn != nil {
		// Sometimes there is strange temperature
		if readAmbientTempReturn.Temperature-ambientTemp > 4 || ambientTemp-readAmbientTempReturn.Temperature > 4 {
			logger.Error().Interface("input", readAmbientTempInput).Str("device", input.Mac).Msg("failed to read the ambient temperature")
			return models.ErrorInvalidParameterTemperature
		}
	}

	logger.Debug().Float32("ambientTemp", ambientTemp).Str("device", input.Mac).Msg("Ambient temperature")

	if readAmbientTempReturn == nil || readAmbientTempReturn.Temperature != ambientTemp {
		// Sent  temperature to MQTT
		publishAmbientTempInput := &models_mqtt.PublishAmbientTempInput{
			Mac:         input.Mac,
			Temperature: ambientTemp,
		}

		err = s.mqtt.PublishAmbientTemp(ctx, logger, publishAmbientTempInput)
		if err != nil {
			logger.Error().Interface("input", publishAmbientTempInput).Msg("failed to publish ambient temperature")
			return err
		}

		// Save the new value in storage
		upsertAmbientTempInput := &models_repo.UpsertAmbientTempInput{Temperature: ambientTemp, Mac: input.Mac}

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

	if uint16(response.Payload[0x22])|(uint16(response.Payload[0x23])<<8) != 0 {
		logger.Error().Err(models.ErrorInvalidResultPacket).Interface("input", sendCommandInput).Msg("Checksum is incorrect")
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
		Temperature:        float32(8+(response.Payload[10]>>3)) + 0.5*float32(response.Payload[12]>>7),
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

	if raw.Temperature < 16.0 {
		logger.Error().Err(models.ErrorInvalidResultPacketLength).Str("device", input.Mac).Float32("temperature", raw.Temperature).Msg("Wrong temperature, skip package")
		return models.ErrorInvalidResultPacketLength
	}

	//////////////////////////////////////////////////////////////////
	//  Compare new statuses with old statuses and update  MQTT     //
	//////////////////////////////////////////////////////////////////

	readDeviceStatusRawInput := &models_repo.ReadDeviceStatusRawInput{
		Mac: input.Mac,
	}

	readDeviceStatusRawReturn, err := s.cache.ReadDeviceStatusRaw(ctx, logger, readDeviceStatusRawInput)
	if err != nil {
		switch err {
		case models_repo.ErrorDeviceStatusRawNotFound:
			err = nil
		default:
			logger.Error().Err(err).Interface("input", readDeviceStatusRawInput).Msg("Failed to read the device status")
			return err
		}
	}

	deviceStatusHass := raw.ConvertToDeviceStatusHass()
	logger.Debug().Interface("status", deviceStatusHass).Str("device", input.Mac).Msg("The converted current device status")

	g := new(errgroup.Group)
	g.Go(func() error {
		if readDeviceStatusRawReturn == nil ||
			readDeviceStatusRawReturn.Status.Temperature != raw.Temperature {

			publishTemperatureInput := &models_mqtt.PublishTemperatureInput{
				Mac:         input.Mac,
				Temperature: deviceStatusHass.Temperature,
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
		if readDeviceStatusRawReturn == nil ||
			readDeviceStatusRawReturn.Status.Mode != raw.Mode ||
			readDeviceStatusRawReturn.Status.Power != raw.Power {

			publishModeInput := &models_mqtt.PublishModeInput{
				Mac:  input.Mac,
				Mode: deviceStatusHass.Mode,
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
		if readDeviceStatusRawReturn == nil ||
			readDeviceStatusRawReturn.Status.FanSpeed != raw.FanSpeed ||
			readDeviceStatusRawReturn.Status.Mute != raw.Mute ||
			readDeviceStatusRawReturn.Status.Turbo != raw.Turbo {

			publishFanModeInput := &models_mqtt.PublishFanModeInput{
				Mac:     input.Mac,
				FanMode: deviceStatusHass.FanMode,
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
		if readDeviceStatusRawReturn == nil ||
			readDeviceStatusRawReturn.Status.FixationVertical != raw.FixationVertical {

			publishSwingModeInput := &models_mqtt.PublishSwingModeInput{
				Mac:       input.Mac,
				SwingMode: deviceStatusHass.SwingMode,
			}

			err = s.mqtt.PublishSwingMode(ctx, logger, publishSwingModeInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishSwingModeInput).Msg("Failed to publish the device swing mode")
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		if readDeviceStatusRawReturn == nil ||
			readDeviceStatusRawReturn.Status.Display != raw.Display {

			publishDisplaySwitchInput := &models_mqtt.PublishDisplaySwitchInput{
				Mac:    input.Mac,
				Status: deviceStatusHass.DisplaySwitch,
			}

			err = s.mqtt.PublishDisplaySwitch(ctx, logger, publishDisplaySwitchInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishDisplaySwitchInput).Msg("Failed to publish the display switch status")
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
		logger.Error().Err(err).Interface("input", input).Msg("failed to read device config")
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
		logger.Error().Err(err).Interface("input", input).Str("device", input.Mac).Msg("failed to send a command")
		return nil, err
	}

	return &models.SendCommandReturn{Payload: sendCommandReturn.Payload}, nil
}

func (s *service) PublishDiscoveryTopic(ctx context.Context, logger *zerolog.Logger, input *models.PublishDiscoveryTopicInput) error {

	prefix := s.topicPrefix + "/" + input.Device.Mac

	device := models_mqtt.DiscoveryTopicDevice{
		Model: "AirCon",
		Mf:    "broadlink",
		Sw:    "v1.4.1",
		Ids:   input.Device.Mac,
		Name:  input.Device.Name,
	}

	availability := models_mqtt.DiscoveryTopicAvailability{
		PayloadAvailable:    models.StatusOnline,
		PayloadNotAvailable: models.StatusOffline,
		Topic:               prefix + "/availability/value",
	}

	publishClimateDiscoveryTopicInput := models_mqtt.PublishClimateDiscoveryTopicInput{
		Topic: models_mqtt.ClimateDiscoveryTopic{
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
			Precision:               0.1,
			Device:                  device,
			UniqueId:                input.Device.Mac + "_ac",
			Availability:            availability,
			CurrentTemperatureTopic: prefix + "/current_temp/value",
			Name:                    input.Device.Name + " AC",
		},
	}
	err := s.mqtt.PublishClimateDiscoveryTopic(ctx, logger, publishClimateDiscoveryTopicInput)
	if err != nil {
		return err
	}

	publishSwitchScreenDiscoveryTopicInput := models_mqtt.PublishSwitchDiscoveryTopicInput{
		Topic: models_mqtt.SwitchDiscoveryTopic{
			Device:       device,
			Name:         input.Device.Name + " Screen",
			UniqueId:     input.Device.Mac + "_screen",
			StateTopic:   prefix + "/display/switch/value",
			CommandTopic: prefix + "/display/switch/set",
			Availability: availability,
		},
	}

	err = s.mqtt.PublishSwitchDiscoveryTopic(ctx, logger, publishSwitchScreenDiscoveryTopicInput)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateFanMode(ctx context.Context, logger *zerolog.Logger, input *models.UpdateFanModeInput) error {

	err := input.Validate()
	if err != nil {
		logger.Error().Err(err).Interface("input", input).Str("device", input.Mac).Msg("input data is not valid")
		return err
	}

	upsertMqttFanModeMessageInput := &models_repo.UpsertMqttFanModeMessageInput{
		Mac: input.Mac,
		FanMode: models_repo.MqttFanModeMessage{
			UpdatedAt: time.Now(),
			FanMode:   input.FanMode,
		},
	}

	err = s.cache.UpsertMqttFanModeMessage(ctx, logger, upsertMqttFanModeMessageInput)
	if err != nil {
		logger.Error().Interface("input", upsertMqttFanModeMessageInput).Str("device", input.Mac).Msg("failed to save mqtt message to cache storage")
		return err
	}

	publishFanModeInput := &models_mqtt.PublishFanModeInput{
		Mac:     input.Mac,
		FanMode: input.FanMode,
	}
	err = s.mqtt.PublishFanMode(ctx, logger, publishFanModeInput)
	if err != nil {
		logger.Error().Interface("input", publishFanModeInput).Str("device", input.Mac).Msg("failed to publish fan mode to mqtt")
		return err
	}

	return nil
}

func (s *service) UpdateMode(ctx context.Context, logger *zerolog.Logger, input *models.UpdateModeInput) error {

	err := input.Validate()
	if err != nil {
		logger.Error().Err(err).Interface("input", input).Str("device", input.Mac).Msg("input data is not valid")
		return err
	}

	upsertMqttModeMessageInput := &models_repo.UpsertMqttModeMessageInput{
		Mac: input.Mac,
		Mode: models_repo.MqttModeMessage{
			UpdatedAt: time.Now(),
			Mode:      input.Mode,
		},
	}

	err = s.cache.UpsertMqttModeMessage(ctx, logger, upsertMqttModeMessageInput)
	if err != nil {
		logger.Error().Interface("input", upsertMqttModeMessageInput).Str("device", input.Mac).Msg("failed to save mqtt message to cache storage")
		return err
	}

	publishModeInput := &models_mqtt.PublishModeInput{
		Mac:  input.Mac,
		Mode: input.Mode,
	}
	err = s.mqtt.PublishMode(ctx, logger, publishModeInput)
	if err != nil {
		logger.Error().Interface("input", publishModeInput).Str("device", input.Mac).Msg("failed to publish mode to mqtt")
		return err
	}

	return nil
}

func (s *service) UpdateSwingMode(ctx context.Context, logger *zerolog.Logger, input *models.UpdateSwingModeInput) error {

	err := input.Validate()
	if err != nil {
		logger.Error().Err(err).Interface("input", input).Str("device", input.Mac).Msg("input data is not valid")
		return err
	}

	upsertMqttSwingModeMessageInput := &models_repo.UpsertMqttSwingModeMessageInput{
		Mac: input.Mac,
		SwingMode: models_repo.MqttSwingModeMessage{
			UpdatedAt: time.Now(),
			SwingMode: input.SwingMode,
		},
	}

	err = s.cache.UpsertMqttSwingModeMessage(ctx, logger, upsertMqttSwingModeMessageInput)
	if err != nil {
		logger.Error().Interface("input", upsertMqttSwingModeMessageInput).Str("device", input.Mac).Msg("failed to save mqtt message to cache storage")
		return err
	}

	publishSwingModeInput := &models_mqtt.PublishSwingModeInput{
		Mac:       input.Mac,
		SwingMode: input.SwingMode,
	}
	err = s.mqtt.PublishSwingMode(ctx, logger, publishSwingModeInput)
	if err != nil {
		logger.Error().Interface("input", publishSwingModeInput).Str("device", input.Mac).Msg("failed to publish swing mode to mqtt")
		return err
	}

	return nil
}

func (s *service) UpdateTemperature(ctx context.Context, logger *zerolog.Logger, input *models.UpdateTemperatureInput) error {

	err := input.Validate()
	if err != nil {
		logger.Error().Err(err).Interface("input", input).Str("device", input.Mac).Msg("input data is not valid")
		return err
	}

	upsertMqttTemperatureMessageInput := &models_repo.UpsertMqttTemperatureMessageInput{
		Mac: input.Mac,
		Temperature: models_repo.MqttTemperatureMessage{
			UpdatedAt:   time.Now(),
			Temperature: input.Temperature,
		},
	}

	err = s.cache.UpsertMqttTemperatureMessage(ctx, logger, upsertMqttTemperatureMessageInput)
	if err != nil {
		logger.Error().Interface("input", upsertMqttTemperatureMessageInput).Str("device", input.Mac).Msg("failed to save mqtt message to cache storage")
		return err
	}

	return nil
}

func (s *service) UpdateDisplaySwitch(ctx context.Context, logger *zerolog.Logger, input *models.UpdateDisplaySwitchInput) error {

	err := input.Validate()
	if err != nil {
		logger.Error().Err(err).Interface("input", input).Str("device", input.Mac).Msg("input data is not valid")
		return err
	}

	isDisplayOn := false
	if input.Status == "ON" {
		isDisplayOn = true
	}

	upsertDisplaySwitchMessageInput := &models_repo.UpsertMqttDisplaySwitchMessageInput{
		Mac: input.Mac,
		DisplaySwitch: models_repo.MqttDisplaySwitchMessage{
			UpdatedAt:   time.Now(),
			IsDisplayOn: isDisplayOn,
		},
	}

	err = s.cache.UpsertMqttDisplaySwitchMessage(ctx, logger, upsertDisplaySwitchMessageInput)
	if err != nil {
		logger.Error().Err(err).Interface("input", upsertDisplaySwitchMessageInput).Str("device", input.Mac).Msg("failed to save mqtt message to cache storage")
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

	// SWING MODE
	var verticalFixation byte
	if input.SwingMode != nil {
		key, ok := models.VerticalFixationStatusesInvert[*input.SwingMode]
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

	// TEMPERATURE
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

	// FAN MODE
	var fanMode, turbo, mute byte
	if input.FanMode != nil {
		if *input.FanMode == "mute" {
			mute = models.StatusOn
		} else if *input.FanMode == "turbo" {
			turbo = models.StatusOn
		} else {
			key, ok := models.FanStatusesInvert[*input.FanMode]
			if !ok {
				logger.Error().Interface("input", input).Str("device", input.Mac).
					Str("fanMode", *input.FanMode).
					Err(models.ErrorInvalidParameterFanMode).
					Msg("Invalid parameter fan mode")

				return models.ErrorInvalidParameterFanMode
			} else {
				fanMode = byte(key)
				turbo = models.StatusOff
				mute = models.StatusOff
			}
		}
	} else {
		fanMode = readDeviceStatusRawReturn.Status.FanSpeed
		mute = readDeviceStatusRawReturn.Status.Mute
		turbo = readDeviceStatusRawReturn.Status.Turbo
	}

	// DISPLAY
	// Attention. Inverted logic
	// Byte 0 - turn ON, Byte 1 - turn OFF
	var displaySwitch byte = 1
	if input.IsDisplayOn != nil {
		if *input.IsDisplayOn {
			displaySwitch = 0
		}
	} else {
		displaySwitch = readDeviceStatusRawReturn.Status.Display
	}

	// MODE
	var mode, power byte
	if input.Mode != nil {
		switch strings.ToLower(*input.Mode) {
		case "cool":
			mode = byte(models.ModeStatusesInvert["cool"])
			power = models.StatusOn
		case "heat":
			mode = byte(models.ModeStatusesInvert["heat"])
			power = models.StatusOn
		case "auto":
			mode = byte(models.ModeStatusesInvert["auto"])
			power = models.StatusOn
		case "dry":
			mode = byte(models.ModeStatusesInvert["dry"])
			power = models.StatusOn
		case "fan_only":
			mode = byte(models.ModeStatusesInvert["fan_only"])
			power = models.StatusOn
		case "off":
			power = models.StatusOff
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
	payload[20] = 0b00000000 | displaySwitch<<4 | readDeviceStatusRawReturn.Status.Mildew<<3
	payload[21] = 0b00000000
	payload[22] = 0b00000000

	// Add checksum
	var payloadChecksum [32]byte
	payloadChecksum[0] = byte(len(payload) + 2)

	copy(payloadChecksum[2:], payload[:])

	var checksum int
	for i := 0; i < len(payload); i += 2 {
		checksum += int(payload[i])<<8 + int(append(payload[:], byte(0))[i+1])
	}
	checksum = (checksum >> 16) + (checksum & 0xFFFF)
	checksum = ^checksum & 0xFFFF

	payloadChecksum[len(payload)+2] = byte((checksum >> 8) & 0xFF)
	payloadChecksum[len(payload)+3] = byte(checksum & 0xFF)

	sendCommandInput := &models.SendCommandInput{
		Command: 0x6a,
		Payload: payloadChecksum[:],
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

	upsertDeviceAvailabilityInput := &models_repo.UpsertDeviceAvailabilityInput{
		Mac:          input.Mac,
		Availability: input.Availability,
	}

	err := s.cache.UpsertDeviceAvailability(ctx, logger, upsertDeviceAvailabilityInput)
	if err != nil {
		logger.Error().Interface("input", upsertDeviceAvailabilityInput).Str("device", input.Mac).Msg("failed to upsert device availability")
		return err
	}

	publishAvailabilityInput := &models_mqtt.PublishAvailabilityInput{
		Mac:          input.Mac,
		Availability: input.Availability,
	}

	err = s.mqtt.PublishAvailability(ctx, logger, publishAvailabilityInput)
	if err != nil {
		logger.Error().Interface("input", publishAvailabilityInput).Str("device", input.Mac).Msg("failed to create command payload")
		return err
	}

	return nil
}

func (s *service) StartDeviceMonitoring(ctx context.Context, logger *zerolog.Logger, input *models.StartDeviceMonitoringInput) error {

	var (
		modeUpdatedTime, swingModeUpdatedTime, fanModeUpdatedTime, temperatureUpdatedTime time.Time
		isDisplayOnUpdatedTime                                                            time.Time
		lastGetDeviceState, lastGetAmbientTemp                                            time.Time

		isDeviceAvailable bool
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if time.Now().Sub(lastGetAmbientTemp).Seconds() > 180 {
				err := s.GetDeviceAmbientTemperature(ctx, logger, &models.GetDeviceAmbientTemperatureInput{Mac: input.Mac})
				if err != nil {
					logger.Error().Str("device", input.Mac).Msg("failed to get ambient temperature")
					err = nil
					continue
				}
				lastGetAmbientTemp = time.Now()
			} else {
				var (
					forcedUpdateDeviceState  = false
					mode, swingMode, fanMode *string
					temperature              *float32
					isDisplayOn              *bool
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
						forcedUpdateDeviceState = true
						mode = &message.Mode.Mode
					}
				}

				if message.FanMode != nil {
					if message.FanMode.UpdatedAt != fanModeUpdatedTime {
						forcedUpdateDeviceState = true
						fanMode = &message.FanMode.FanMode
					}
				}

				if message.SwingMode != nil {
					if message.SwingMode.UpdatedAt != swingModeUpdatedTime {
						forcedUpdateDeviceState = true
						swingMode = &message.SwingMode.SwingMode
					}
				}

				if message.Temperature != nil {
					if message.Temperature.UpdatedAt != temperatureUpdatedTime {
						forcedUpdateDeviceState = true
						temperature = &message.Temperature.Temperature
					}
				}

				if message.IsDisplayOn != nil {
					if message.IsDisplayOn.UpdatedAt != isDisplayOnUpdatedTime {
						forcedUpdateDeviceState = true
						isDisplayOn = &message.IsDisplayOn.IsDisplayOn
					}
				}

				if forcedUpdateDeviceState || int(time.Now().Sub(lastGetDeviceState).Seconds()) > s.updateInterval {
					for {
						err = s.GetDeviceStates(ctx, logger, &models.GetDeviceStatesInput{Mac: input.Mac})
						if err != nil {
							logger.Error().Err(err).Interface("device", input.Mac).Msg("Failed to get AC States")

							// If we cannot receive data from the air conditioner within three intervals,
							// then we send the status that the air conditioner is unavailable
							if time.Now().Sub(lastGetDeviceState).Seconds() > float64(s.updateInterval)*3 && isDeviceAvailable {
								updateDeviceAvailabilityInput := &models.UpdateDeviceAvailabilityInput{
									Mac:          input.Mac,
									Availability: models.StatusOffline,
								}
								err = s.UpdateDeviceAvailability(ctx, logger, updateDeviceAvailabilityInput)
								if err != nil {
									logger.Error().Err(err).Str("device", input.Mac).Interface("input", updateDeviceAvailabilityInput).Msg("Failed to update device availability")
									err = nil
								}
								isDeviceAvailable = false
							}
							err = nil
							continue
						} else {
							lastGetDeviceState = time.Now()
							if !isDeviceAvailable {
								updateDeviceAvailabilityInput := &models.UpdateDeviceAvailabilityInput{
									Mac:          input.Mac,
									Availability: models.StatusOnline,
								}
								err = s.UpdateDeviceAvailability(ctx, logger, updateDeviceAvailabilityInput)
								if err != nil {
									logger.Error().Err(err).Str("device", input.Mac).Interface("input", updateDeviceAvailabilityInput).Msg("Failed to update device availability")
									err = nil
								}
								isDeviceAvailable = true
							}
							break
						}
					}
				}

				if forcedUpdateDeviceState && isDeviceAvailable {
					// A short pause before sending a new message to the air conditioner so that it does not hang
					time.Sleep(time.Millisecond * 500)

					updateDeviceStatesInput := &models.UpdateDeviceStatesInput{
						Mac:         input.Mac,
						FanMode:     fanMode,
						SwingMode:   swingMode,
						Mode:        mode,
						Temperature: temperature,
						IsDisplayOn: isDisplayOn,
					}
					err := s.UpdateDeviceStates(ctx, logger, updateDeviceStatesInput)
					if err != nil {
						logger.Error().Err(err).Str("device", input.Mac).Interface("input", updateDeviceStatesInput).Msg("Failed to update device states")
						err = nil
						continue
					}

					// Reset the time of the last update to get fresh data from the air conditioner
					lastGetDeviceState = time.UnixMicro(0)

					if message.Mode != nil {
						modeUpdatedTime = message.Mode.UpdatedAt
					}
					if message.FanMode != nil {
						fanModeUpdatedTime = message.FanMode.UpdatedAt
					}
					if message.SwingMode != nil {
						swingModeUpdatedTime = message.SwingMode.UpdatedAt
					}
					if message.Temperature != nil {
						temperatureUpdatedTime = message.Temperature.UpdatedAt
					}
					if message.IsDisplayOn != nil {
						isDisplayOnUpdatedTime = message.IsDisplayOn.UpdatedAt
					}
				}

				time.Sleep(time.Millisecond * 500)
			}
		}
	}
}

func (s *service) PublishStatesOnHomeAssistantRestart(ctx context.Context, logger *zerolog.Logger, input *models.PublishStatesOnHomeAssistantRestartInput) error {

	if input.Status != models.StatusOnline {
		return nil
	}

	readAuthedDevicesReturn, err := s.cache.ReadAuthedDevices(ctx, logger)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read authed devices")
		return err
	}

	for _, mac := range readAuthedDevicesReturn.Macs {

		/////////////////////////////////
		// Read all states and configs //
		/////////////////////////////////

		readDeviceStatusRawInput := &models_repo.ReadDeviceStatusRawInput{
			Mac: mac,
		}

		readDeviceStatusRawReturn, err := s.cache.ReadDeviceStatusRaw(ctx, logger, readDeviceStatusRawInput)
		if err != nil {
			logger.Error().Err(err).Interface("input", readDeviceStatusRawInput).Msg("Failed to read the device status")
			return err
		}

		hassStatus := readDeviceStatusRawReturn.Status.ConvertToDeviceStatusHass()

		readAmbientTempInput := &models_repo.ReadAmbientTempInput{Mac: mac}

		readAmbientTempReturn, err := s.cache.ReadAmbientTemp(ctx, logger, readAmbientTempInput)
		if err != nil {
			logger.Error().Interface("input", readAmbientTempInput).Str("device", mac).Msg("failed to read the ambient temperature")
			return err
		}

		readDeviceAvailabilityInput := &models_repo.ReadDeviceAvailabilityInput{Mac: mac}

		readDeviceAvailabilityReturn, err := s.cache.ReadDeviceAvailability(ctx, logger, readDeviceAvailabilityInput)
		if err != nil {
			logger.Error().Interface("input", readDeviceAvailabilityInput).Str("device", mac).Msg("failed to read the device availability")
			return err
		}

		readDeviceConfigInput := &models_repo.ReadDeviceConfigInput{
			Mac: mac,
		}

		readDeviceConfigReturn, err := s.cache.ReadDeviceConfig(ctx, logger, readDeviceConfigInput)
		if err != nil {
			logger.Error().Err(err).Interface("input", input).Msg("failed to read device config")
			return err
		}

		/////////////////////////////////
		// 		Publish all topics     //
		/////////////////////////////////

		err = s.PublishDiscoveryTopic(ctx, logger, &models.PublishDiscoveryTopicInput{Device: readDeviceConfigReturn.Config})
		if err != nil {
			logger.Error().Interface("device", mac).Msg("failed to publish the discovery topic")
			return err
		}

		time.Sleep(time.Millisecond * 500)

		g := new(errgroup.Group)

		g.Go(func() error {
			publishAvailabilityInput := &models_mqtt.PublishAvailabilityInput{
				Mac:          mac,
				Availability: readDeviceAvailabilityReturn.Availability,
			}

			err := s.mqtt.PublishAvailability(ctx, logger, publishAvailabilityInput)
			if err != nil {
				logger.Error().Interface("input", publishAvailabilityInput).Str("device", mac).Msg("failed to create command payload")
				return err
			}

			return nil
		})

		g.Go(func() error {

			// Send  temperature to MQTT
			publishAmbientTempInput := &models_mqtt.PublishAmbientTempInput{
				Mac:         mac,
				Temperature: readAmbientTempReturn.Temperature,
			}

			err = s.mqtt.PublishAmbientTemp(ctx, logger, publishAmbientTempInput)
			if err != nil {
				logger.Error().Interface("input", publishAmbientTempInput).Msg("failed to publish ambient temperature")
				return err
			}

			return nil
		})

		g.Go(func() error {

			publishTemperatureInput := &models_mqtt.PublishTemperatureInput{
				Mac:         mac,
				Temperature: readDeviceStatusRawReturn.Status.Temperature,
			}

			err = s.mqtt.PublishTemperature(ctx, logger, publishTemperatureInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishTemperatureInput).Msg("Failed to publish the device set temperature")
				return err
			}

			return nil
		})

		g.Go(func() error {

			publishModeInput := &models_mqtt.PublishModeInput{
				Mac:  mac,
				Mode: hassStatus.Mode,
			}

			err = s.mqtt.PublishMode(ctx, logger, publishModeInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishModeInput).Msg("Failed to publish the device mode")
				return err
			}

			return nil
		})

		g.Go(func() error {

			publishFanModeInput := &models_mqtt.PublishFanModeInput{
				Mac:     mac,
				FanMode: hassStatus.FanMode,
			}

			err = s.mqtt.PublishFanMode(ctx, logger, publishFanModeInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishFanModeInput).Msg("Failed to publish the device fan mode")
				return err
			}

			return nil
		})

		g.Go(func() error {

			publishSwingModeInput := &models_mqtt.PublishSwingModeInput{
				Mac:       mac,
				SwingMode: hassStatus.SwingMode,
			}

			err = s.mqtt.PublishSwingMode(ctx, logger, publishSwingModeInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishSwingModeInput).Msg("Failed to publish the device swing mode")
				return err
			}

			return nil
		})

		g.Go(func() error {

			publishDisplaySwitchInput := &models_mqtt.PublishDisplaySwitchInput{
				Mac:    mac,
				Status: hassStatus.DisplaySwitch,
			}

			err = s.mqtt.PublishDisplaySwitch(ctx, logger, publishDisplaySwitchInput)
			if err != nil {
				logger.Error().Err(err).Interface("input", publishDisplaySwitchInput).Msg("Failed to publish the display switch status")
				return err
			}

			return nil
		})

		// Wait for all HTTP fetches to complete.
		if err = g.Wait(); err != nil {
			return err
		}

	}

	return nil
}
