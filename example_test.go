package zapconfig_test

import "moul.io/zapconfig"

func Example() {
	logger := zapconfig.Configurator{}.MustBuild()
	logger.Info("hello!")
}
