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
	opts       zap.Config
	stacktrace bool
}

// SetOutputPath sets zap.Config.OutputPaths and c.Config.ErrorOutputPaths with the given path.
func (c *Configurator) SetOutputPath(dest string) *Configurator {
	c.opts.OutputPaths = []string{dest}
	c.opts.ErrorOutputPaths = []string{dest}
	return c
}

// SetOutputPaths sets zap.Config.OutputPaths and c.Config.ErrorOutputPaths with the given paths.
func (c *Configurator) SetOutputPaths(dests []string) *Configurator {
	c.opts.OutputPaths = dests
	c.opts.ErrorOutputPaths = dests
	return c
}

// EnableStacktrace forces stacktraces to be enabled.
func (c *Configurator) EnableStacktrace() *Configurator {
	c.stacktrace = true
	return c
}

// SetLevel sets the minimal logging level.
func (c *Configurator) SetLevel(level zapcore.Level) *Configurator {
	c.opts.Level = zap.NewAtomicLevelAt(level)
	return c
}

// IsEmpty checks whether the Configurator isn't touched (default value) or if it was modified.
func (c Configurator) IsEmpty() bool {
	return reflect.DeepEqual(c, Configurator{})
}

// Config builds a zap.Config.
func (c Configurator) Config() (zap.Config, error) {
	copy := c.opts

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
	copy.DisableStacktrace = !c.stacktrace
	// FIXME: based on a flag or guess with env vars
	copy.Development = true

	return copy, nil
}

// String implements Stringer.
func (c Configurator) String() string {
	copy := c.opts
	if copy.Level == (zap.AtomicLevel{}) {
		copy.Level = zap.NewAtomicLevel()
	}
	return fmt.Sprintf("%v", copy)
}

// BuildLogger returns a configured *zap.Logger.
func (c Configurator) Build() (*zap.Logger, error) {
	config, err := c.Config()
	if err != nil {
		return nil, err
	}
	return config.Build()
}

func (c Configurator) MustBuild() *zap.Logger {
	logger, err := c.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
