package settings

import (
	"github.com/caarlos0/env/v6"
	"github.com/rotisserie/eris"
)

type values struct {
	NatsURL       string `env:"NATS_URL,required"`
}

var Values *values

func Load() error {
	Values = &values{}

	if err := env.Parse(Values); err != nil {
		return eris.Wrap(err, "failed to parse values")
	}

	return nil
}
