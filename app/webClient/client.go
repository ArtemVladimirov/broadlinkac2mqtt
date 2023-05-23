package webClient

import (
	"context"
	"github.com/ArtemVladimirov/broadlinkac2mqtt/app/webClient/models"
	"github.com/rs/zerolog"
	"net"
	"strconv"
	"time"
)

type webClient struct {
}

func NewWebClient() *webClient {
	return &webClient{}
}

func (w *webClient) SendCommand(ctx context.Context, logger *zerolog.Logger, input *models.SendCommandInput) (*models.SendCommandReturn, error) {

	conn, err := net.Dial("udp", input.Ip+":"+strconv.Itoa(int(input.Port)))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to dial address")
		return nil, err
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to set deadline")
		return nil, err
	}

	_, err = conn.Write(input.Payload)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to write the payload")
		return nil, err
	}

	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read the response")
		return nil, err
	}

	return &models.SendCommandReturn{Payload: response}, nil
}
