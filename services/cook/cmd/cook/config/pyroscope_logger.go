package config

import (
	"fmt"

	"github.com/grafana/pyroscope-go"
	"github.com/rs/zerolog/log"
)

var _ pyroscope.Logger = &PyroscopeLogger{}

// PyroscopeLogger provides an adapter for pyroscope.Logger
type PyroscopeLogger struct{}

func (*PyroscopeLogger) Debugf(a string, b ...interface{}) {
	log.Debug().Msgf(fmt.Sprintf("pyroscope: %s", a), b...)
}

func (*PyroscopeLogger) Errorf(a string, b ...interface{}) {
	log.Error().Msgf(fmt.Sprintf("pyroscope: %s", a), b...)
}

func (*PyroscopeLogger) Infof(a string, b ...interface{}) {
	log.Info().Msgf(fmt.Sprintf("pyroscope: %s", a), b...)
}
