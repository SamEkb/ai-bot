package bot

import (
	"errors"
	"time"
)

type Config struct {
	Token   string
	Timeout time.Duration
}

const defaultTimeout = 10

var ErrTokenRequired = errors.New("token does not exist")

func NewConfig(token string, timeout time.Duration) *Config {
	return &Config{
		Token:   token,
		Timeout: timeout,
	}
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return ErrTokenRequired
	}

	if c.Timeout <= 0 {
		c.Timeout = defaultTimeout
	}

	return nil
}
