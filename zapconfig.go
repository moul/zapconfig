package zapconfig

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Configurator is the main object of this package.
// Leaving it empty will generate an opinionated sane default zap.Config.
type Configurator struct {
	zap.Config
}

// SetOutputPath sets zap.Config.OutputPaths and c.Config.ErrorOutputPaths with the given path.
func (c *Configurator) SetOutputPath(dest string) {
	c.OutputPaths = []string{dest}
	c.ErrorOutputPaths = []string{dest}
}

// SetOutputPaths sets zap.Config.OutputPaths and c.Config.ErrorOutputPaths with the given paths.
func (c *Configurator) SetOutputPaths(dests []string) {
	c.OutputPaths = dests
	c.ErrorOutputPaths = dests
}

// IsEmpty checks whether the Configurator isn't touched (default value) or if it was modified.
func (c Configurator) IsEmpty() bool {
	return reflect.DeepEqual(c, Configurator{})
}

// BuildConfig builds a zap.Config.
func (c Configurator) BuildConfig() (zap.Config, error) {
	copy := c.Config

	if copy.OutputPaths == nil {
		copy.OutputPaths = []string{"stderr"}
	}
	if copy.ErrorOutputPaths == nil {
		copy.ErrorOutputPaths = []string{"stderr"}
	}
	if copy.Level == (zap.AtomicLevel{}) {
		copy.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	if copy.Encoding == "" {
		copy.Encoding = "console"
	}
	if reflect.DeepEqual(copy.EncoderConfig, zapcore.EncoderConfig{}) {
		copy.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		copy.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	// FIXME: easy way to disable this
	copy.DisableStacktrace = true
	// FIXME: based on a flag or guess with env vars
	copy.Development = true

	return copy, nil
}

// String implements Stringer.
func (c Configurator) String() string {
	if c.IsEmpty() {
		return "{}"
	}
	return fmt.Sprintf("%v", c.Config)
}

// MustBuildConfig is equivalent to BuildConfig but will panic in case of error instead of returning it.
func (c Configurator) MustBuildConfig() zap.Config {
	config, err := c.BuildConfig()
	if err != nil {
		panic(err)
	}
	return config
}

// BuildLogger returns a configured *zap.Logger.
func (c Configurator) BuildLogger() (*zap.Logger, error) {
	config, err := c.BuildConfig()
	if err != nil {
		return nil, err
	}
	return config.Build()
}

func (c Configurator) MustBuildLogger() *zap.Logger {
	logger, err := c.BuildLogger()
	if err != nil {
		panic(err)
	}
	return logger
}
