package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Service Service   `yaml:"service" json:"service"`
		Mqtt    Mqtt      `yaml:"mqtt" json:"mqtt"`
		Devices []Devices `yaml:"devices" json:"devices"`
	}

	Service struct {
		UpdateInterval int    `env-default:"10"    yaml:"update_interval" json:"update_interval"`
		LogLevel       string `env-default:"error" yaml:"log_level" json:"log_level"`
	}

	Mqtt struct {
		Broker                   string  `env-required:"true" yaml:"broker" json:"broker"`
		User                     *string `yaml:"user" json:"user"`
		Password                 *string `yaml:"password" json:"password"`
		ClientId                 string  `env-default:"broadlinkac" yaml:"client_id" json:"client_id"`
		TopicPrefix              string  `env-default:"airac" yaml:"topic_prefix" json:"topic_prefix"`
		AutoDiscoveryTopic       *string `yaml:"auto_discovery_topic" json:"auto_discovery_topic"`
		AutoDiscoveryTopicRetain bool    `env-default:"true" yaml:"auto_discovery_topic_retain" json:"auto_discovery_topic_retain"`
		CertificateAuthority     *string `yaml:"certificate_authority" json:"certificate_authority"`
		SkipCertCnCheck          bool    `env-default:"true" yaml:"skip_cert_cn_check" json:"skip_cert_cn_check"`
		CertificateClient        *string `yaml:"certificate_client" json:"certificate_client"`
		KeyClient                *string `yaml:"key-client" json:"key_client"`
	}

	Devices struct {
		Ip   string `env-required:"true" yaml:"ip" json:"ip"`
		Mac  string `env-required:"true" yaml:"mac" json:"mac"`
		Name string `env-required:"true" yaml:"name" json:"name"`
		Port uint16 `env-required:"true" yaml:"port" json:"port"`
		// TemperatureUnit defines the temperature unit of the device, C or F.
		// If this is not set, the temperature unit is Celsius.
		TemperatureUnit string `env-default:"C" yaml:"temperature_unit" json:"temperature_unit"` // BUG cleanenv env-default is not working
	}
)

// NewConfig returns app config.
func NewConfig(logger *slog.Logger) (*Config, error) {
	cfg := &Config{}

	files := [...]string{
		"./config/config.yml",
		"./config/config.json",
		"/data/options.json", // hassio
	}

	for i := range files {
		if _, err := os.Stat(files[i]); err == nil {
			err := cleanenv.ReadConfig(files[i], cfg)
			if err != nil {
				logger.Error("failed to read config", slog.Any("err", err))
				return nil, err
			}
			return cfg, nil
		}
	}

	return nil, errors.New("config file is not found")
}
