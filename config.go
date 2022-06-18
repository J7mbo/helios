package helios

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

var InvalidConfigurationError = errors.New("invalid configuration provided")

type Config struct {
	configs []configuration
}

type configuration interface {
	validate() error
}

func NewConfig(configs ...configuration) (*Config, error) {
	for _, config := range configs {
		if err := config.validate(); err != nil {
			return nil, fmt.Errorf(
				"%w for %s: %s",
				InvalidConfigurationError,
				reflect.TypeOf(config).Elem().Name(),
				err.Error(),
			)
		}
	}

	return &Config{configs: configs}, nil
}

type PollInterval struct {
	time.Duration
}

func (c *PollInterval) validate() error {
	if c.Seconds() < 0.1 {
		return errors.New("poll time must be greater than or equal to 0.1 seconds")
	}

	return nil
}
