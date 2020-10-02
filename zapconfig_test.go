package zapconfig

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDefault(t *testing.T) {
	configurator := Configurator{}
	require.True(t, configurator.IsEmpty())
	require.Equal(t, configurator.String(), "{info false false false <nil>  {        <nil> <nil> <nil> <nil> <nil> } [] [] map[]}")

	config, err := configurator.Config()
	require.NoError(t, err)
	require.NotNil(t, config)
	require.True(t, configurator.IsEmpty())
	require.NotEqual(t, configurator.opts, config)

	tempfile, err := ioutil.TempFile("", "zapconfig")
	require.NoError(t, err)
	defer os.Remove(tempfile.Name())

	configurator.SetOutputPath(tempfile.Name())
	require.False(t, configurator.IsEmpty())
	config, err = configurator.Config()
	require.NoError(t, err)
	require.NotNil(t, config)
	require.False(t, configurator.IsEmpty())
	require.NotEqual(t, configurator.opts, config)

	logger, err := config.Build()
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Info("Hello World!")
	logger.Sync()

	b, err := ioutil.ReadFile(tempfile.Name())
	require.NoError(t, err)
	require.Contains(t, string(b), "zapconfig/zapconfig_test.go:")
	require.Contains(t, string(b), "\tHello World!")
}

func TestBuilder(t *testing.T) {
	var cases = []struct {
		name    string
		setupFn func(config *Configurator)
	}{
		{"no-action", func(config *Configurator) {}},
		{"setlevel-debug", func(config *Configurator) { config.SetLevel(zap.DebugLevel) }},
		{"setlevel-warn", func(config *Configurator) { config.SetLevel(zap.WarnLevel) }},
		{"setoutput-stdout", func(config *Configurator) { config.SetOutputPath("stdout") }},
		{"setoutput-stderr", func(config *Configurator) { config.SetOutputPath("stderr") }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := Configurator{}
			tc.setupFn(&config)
			logger, err := config.Build()
			require.NoError(t, err)
			require.NotNil(t, logger)
			// fmt.Println(config)
		})
	}
}
