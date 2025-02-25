package webClient

import (
	"context"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/ArtemVladimirov/broadlinkac2mqtt/app/webClient/models"
)

type webClient struct {
	logger *slog.Logger
}

func NewWebClient(logger *slog.Logger) *webClient {
	return &webClient{
		logger: logger,
	}
}

func (w *webClient) SendCommand(ctx context.Context, input *models.SendCommandInput) (*models.SendCommandReturn, error) {
	conn, err := net.Dial("udp", input.Ip+":"+strconv.Itoa(int(input.Port)))
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to dial address", slog.Any("err", err))
		return nil, err
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			w.logger.ErrorContext(ctx, "Failed to close client connection", slog.Any("err", err))
		}
	}(conn)

	err = conn.SetDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to set deadline", slog.Any("err", err))
		return nil, err
	}

	_, err = conn.Write(input.Payload)
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to write the payload", slog.Any("err", err))
		return nil, err
	}

	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to read the response", slog.Any("err", err))
		return nil, err
	}

	return &models.SendCommandReturn{Payload: response}, nil
}
