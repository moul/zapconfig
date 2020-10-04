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
		checkFn (func(t *testing.T, config *zap.Config))
	}{
		{"no-action", nil, nil},
		{"setlevel-debug", func(config *Configurator) { config.SetLevel(zap.DebugLevel) }, nil},
		{"setlevel-warn", func(config *Configurator) { config.SetLevel(zap.WarnLevel) }, nil},
		// SetOutputPath
		{"setoutput-default",
			nil,
			func(t *testing.T, config *zap.Config) {
				require.Equal(t, config.OutputPaths, []string{"stderr"})
				require.Equal(t, config.ErrorOutputPaths, []string{"stderr"})
			},
		},
		{"setoutput-stdout",
			func(config *Configurator) { config.SetOutputPath("stdout") },
			func(t *testing.T, config *zap.Config) {
				require.Equal(t, config.OutputPaths, []string{"stdout"})
				require.Equal(t, config.ErrorOutputPaths, []string{"stdout"})
			},
		},
		{"setoutput-stderr",
			func(config *Configurator) { config.SetOutputPath("stderr") },
			func(t *testing.T, config *zap.Config) {
				require.Equal(t, config.OutputPaths, []string{"stderr"})
				require.Equal(t, config.ErrorOutputPaths, []string{"stderr"})
			},
		},
		// stacktrace
		{"default-stacktrace",
			nil,
			func(t *testing.T, config *zap.Config) { require.True(t, config.DisableStacktrace) },
		},
		{"EnableStacktrace",
			func(config *Configurator) { config.EnableStacktrace() },
			func(t *testing.T, config *zap.Config) { require.False(t, config.DisableStacktrace) },
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			configurator := Configurator{}
			if tc.setupFn != nil {
				tc.setupFn(&configurator)
			}
			config, err := configurator.Config()
			require.NoError(t, err)
			require.NotNil(t, config)

			if tc.checkFn != nil {
				tc.checkFn(t, &config)
			}

			logger, err := config.Build()
			require.NoError(t, err)
			require.NotNil(t, logger)
		})
	}
}
