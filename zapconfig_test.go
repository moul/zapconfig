package zapconfig

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	configurator := Configurator{}
	require.True(t, configurator.IsEmpty())
	require.Equal(t, configurator.String(), "{}")

	config, err := configurator.BuildConfig()
	require.NoError(t, err)
	require.NotNil(t, config)
	require.True(t, configurator.IsEmpty())
	require.NotEqual(t, configurator.Config, config)

	tempfile, err := ioutil.TempFile("", "zapconfig")
	require.NoError(t, err)
	defer os.Remove(tempfile.Name())

	configurator.SetOutputPath(tempfile.Name())
	require.False(t, configurator.IsEmpty())
	config, err = configurator.BuildConfig()
	require.NoError(t, err)
	require.NotNil(t, config)
	require.False(t, configurator.IsEmpty())
	require.NotEqual(t, configurator.Config, config)

	logger, err := config.Build()
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Info("Hello World!")
	logger.Sync()

	b, err := ioutil.ReadFile(tempfile.Name())
	require.NoError(t, err)
	require.Contains(t, string(b), "zapconfig/zapconfig_test.go:38\tHello World!")
}
