package openai

import "errors"

type Config struct {
	Token string
}

func NewConfig(token string) *Config {
	return &Config{Token: token}
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return errors.New("openai token is not valid")
	}

	return nil
}
