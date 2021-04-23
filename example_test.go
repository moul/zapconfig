package zapconfig_test

import (
	"go.uber.org/zap/zapcore"

	"moul.io/zapconfig"
)

func Example() {
	logger := zapconfig.Configurator{}.MustBuild()
	logger.Info("hello!")
}

func Example_configuration() {
	logger := zapconfig.New().
		EnableStacktrace().
		SetLevel(zapcore.DebugLevel).
		SetOutputPath("stderr").
		SetOutputPaths([]string{"stderr", "stdout", "./path/to/log.txt"}).
		SetPreset("light-console").
		MustBuild()
	logger.Info("hello!")
}
