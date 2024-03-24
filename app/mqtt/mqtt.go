package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/ArtemVladimirov/broadlinkac2mqtt/config"
	paho "github.com/eclipse/paho.mqtt.golang"
)

func NewMqttConfig(logger *slog.Logger, cfg config.Mqtt) (*paho.ClientOptions, error) {
	//Configure MQTT Client
	uri, err := url.Parse(cfg.Broker)
	if err != nil {
		message := "URL address is incorrect"
		logger.Error(message)
		return nil, errors.New(message)
	}

	opts := paho.NewClientOptions().AddBroker(uri.String()).SetClientID(cfg.ClientId)

	if cfg.User != nil {
		opts.SetUsername(*cfg.User)
	}
	if cfg.Password != nil {
		opts.SetPassword(*cfg.Password)
	}

	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(false)
	opts.SetConnectionLostHandler(func(client paho.Client, err error) {
		logger.Error("MQTT connection lost", slog.Any("err", err))
	})
	opts.SetOnConnectHandler(func(client paho.Client) {
		logger.Info("Connected to MQTT")
	})

	if uri.Scheme == "mqtts" || uri.Scheme == "ssl" {
		tlsConfig := &tls.Config{}

		if cfg.CertificateClient != nil && cfg.KeyClient != nil {
			cert, err := tls.LoadX509KeyPair(*cfg.CertificateClient, *cfg.KeyClient)
			if err != nil {
				logger.Error("Failed to load the client key pair", slog.Any("err", err))
				return nil, err
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		if cfg.CertificateAuthority != nil {
			caCert, err := os.ReadFile(*cfg.CertificateAuthority)
			if err != nil {
				logger.Error("Failed to load the authority certificate", slog.Any("err", err))
				return nil, err
			}

			// Create a certificate pool and add the CA certificate
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			tlsConfig.RootCAs = caCertPool
		}

		tlsConfig.InsecureSkipVerify = cfg.SkipCertCnCheck

		opts.SetTLSConfig(tlsConfig)
	}

	return opts, err
}
