package zapconfig

import (
	"fmt"
	"os"
	"reflect"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Configurator is the main object of this package.
// Leaving it empty will generate an opinionated sane default zap.Config.
type Configurator struct {
	opts       zap.Config
	stacktrace bool
	configErr  error
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

const (
	consoleEncoding = "console"
	jsonEncoding    = "json"
)

// AvailablePresets is the list of preset supported by `SetPreset`.
var AvailablePresets = []string{"console", "json", "light-console", "light-json"}

// SetPreset configures various things based on just a keyword.
func (c *Configurator) SetPreset(name string) *Configurator {
	switch name {
	case "console":
		c.opts.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		c.opts.Encoding = consoleEncoding
		c.opts.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		c.opts.EncoderConfig.EncodeLevel = stableWidthCapitalColorLevelEncoder
		c.opts.EncoderConfig.EncodeName = stableWidthNameEncoder
		c.opts.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	case "light-console":
		c.opts.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		c.opts.Encoding = consoleEncoding
		c.opts.EncoderConfig.TimeKey = ""
		c.opts.EncoderConfig.EncodeLevel = stableWidthCapitalColorLevelEncoder
		c.opts.EncoderConfig.EncodeName = stableWidthNameEncoder
	case "json":
		c.opts.EncoderConfig = zap.NewProductionEncoderConfig()
		c.opts.Encoding = jsonEncoding
	case "light-json":
		c.opts.EncoderConfig = zap.NewProductionEncoderConfig()
		c.opts.Encoding = jsonEncoding
		c.opts.EncoderConfig.TimeKey = ""
	default:
		c.configErr = multierr.Append(c.configErr, fmt.Errorf("unknown preset: %q", name)) // nolint:goerr113
	}
	return c
}

// IsEmpty checks whether the Configurator isn't touched (default value) or if it was modified.
func (c Configurator) IsEmpty() bool {
	return reflect.DeepEqual(c, Configurator{})
}

// Config builds a zap.Config.
func (c Configurator) Config() (zap.Config, error) {
	if c.configErr != nil {
		return zap.Config{}, c.configErr
	}

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
	switch os.Getenv("ENVIRONMENT") {
	case "production", "prod":
		copy.Development = false
	case "development", "develop", "dev":
		copy.Development = true
	}

	if c.configErr != nil {
		return zap.Config{}, c.configErr
	}

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
