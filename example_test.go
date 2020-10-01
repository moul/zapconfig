package zapconfig_test

import "moul.io/zapconfig"

func Example() {
	logger := zapconfig.Configurator{}.MustBuildLogger()
	logger.Info("hello!")
}
