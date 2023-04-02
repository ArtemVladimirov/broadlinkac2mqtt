package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	"os"
)

type (
	Config struct {
		Service Service   `yaml:"service"`
		Mqtt    Mqtt      `yaml:"mqtt"`
		Devices []Devices `yaml:"devices"`
	}

	Service struct {
		UpdateInterval int    `env-default:"10"    yaml:"update_interval"`
		LogLevel       string `env-default:"error" yaml:"log_level"`
	}

	Mqtt struct {
		Broker                   string  `env-required:"true"            yaml:"broker"`
		User                     *string `yaml:"user"`
		Password                 *string `yaml:"password"`
		ClientId                 string  `env-default:"broadlinkac"      yaml:"client_id"`
		TopicPrefix              string  `env-default:"airac"            yaml:"topic_prefix"`
		AutoDiscoveryTopic       *string `yaml:"auto_discovery_topic"`
		AutoDiscoveryTopicRetain bool    `env-default:"true"             yaml:"auto_discovery_topic_retain"`
		CertificateAuthority     *string `yaml:"certificate_authority"`
		SkipCertCnCheck          bool    `env-default:"true"            yaml:"skip_cert_cn_check"`
		CertificateClient        *string `yaml:"certificate_client"`
		KeyClient                *string `yaml:"key-client"`
	}

	Devices struct {
		Ip   string `env-required:"true" yaml:"ip"`
		Mac  string `env-required:"true" yaml:"mac"`
		Name string `env-required:"true" yaml:"name"`
		Port uint16 `env-required:"true" yaml:"port"`
	}
)

// NewConfig returns app config.
func NewConfig(logger *zerolog.Logger) (*Config, error) {
	logger.Debug().Msg("Start reading a config file")

	cfg := &Config{}

	const configFile = "./config/config.yml"
	if _, err := os.Stat(configFile); err == nil {
		err := cleanenv.ReadConfig(configFile, cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	} else {
		msg := "config file not found"
		logger.Fatal().Msg(msg)
		return nil, errors.New(msg)
	}

	return cfg, nil
}
