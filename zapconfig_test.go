package zapconfig

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"moul.io/u"
)

// logFoobar should stay at the top of this file to avoid changing line numbers in generated outputs.
func logFoobar(t *testing.T, logger *zap.Logger) {
	logger.Debug("hello world!", zap.String("foo", "bar"))
	logger.Named("foobar").Info("hello world!", zap.Int("baz", 42))
}

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
		checkFn func(t *testing.T, config *zap.Config, output []string, err error)
	}{
		{"default",
			func(config *Configurator) {
				// noop
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.Equal(t, config.OutputPaths, []string{"stderr"})
				require.Equal(t, config.ErrorOutputPaths, []string{"stderr"})
				require.True(t, config.DisableStacktrace)
				require.Len(t, output, 2)
				require.Contains(t, output[0], "\t\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}")
				require.Contains(t, output[1], "\t\x1b[34mINFO\x1b[0m\tfoobar\tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
			},
		},
		{"setlevel-debug",
			func(config *Configurator) {
				config.SetLevel(zap.DebugLevel)
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.Len(t, output, 2)
				require.Contains(t, output[0], "\t\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}")
				require.Contains(t, output[1], "\t\x1b[34mINFO\x1b[0m\tfoobar\tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
			},
		},
		{"setlevel-info",
			func(config *Configurator) {
				config.SetLevel(zap.InfoLevel)
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.Len(t, output, 1)
				require.Contains(t, output[0], "\t\x1b[34mINFO\x1b[0m\tfoobar\tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
			},
		},
		{"setlevel-warn",
			func(config *Configurator) {
				config.SetLevel(zap.WarnLevel)
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.Len(t, output, 0)
			},
		},
		{"setoutput-stdout",
			func(config *Configurator) {
				config.SetOutputPath("stdout")
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.Equal(t, config.OutputPaths, []string{"stdout"})
				require.Equal(t, config.ErrorOutputPaths, []string{"stdout"})
				require.Contains(t, output[0], "\t\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}")
				require.Contains(t, output[1], "\t\x1b[34mINFO\x1b[0m\tfoobar\tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
			},
		},
		{"setoutput-stderr",
			func(config *Configurator) {
				config.SetOutputPath("stderr")
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.Equal(t, config.OutputPaths, []string{"stderr"})
				require.Equal(t, config.ErrorOutputPaths, []string{"stderr"})
				require.Contains(t, output[0], "\t\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}")
				require.Contains(t, output[1], "\t\x1b[34mINFO\x1b[0m\tfoobar\tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
			},
		},
		{"enable-stacktrace",
			func(config *Configurator) {
				config.EnableStacktrace()
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.False(t, config.DisableStacktrace)
				require.Contains(t, output[0], "\t\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}")
				require.Contains(t, output[1], "\t\x1b[34mINFO\x1b[0m\tfoobar\tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
			},
		},
		{"set-env-default",
			func(config *Configurator) {
				os.Setenv("ENVIRONMENT", "")
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.False(t, config.Development)
				require.Contains(t, output[0], "\t\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}")
				require.Contains(t, output[1], "\t\x1b[34mINFO\x1b[0m\tfoobar\tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
			},
		},
		{"set-env-production",
			func(config *Configurator) {
				os.Setenv("ENVIRONMENT", "production")
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.False(t, config.Development)
				require.Contains(t, output[0], "\t\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}")
				require.Contains(t, output[1], "\t\x1b[34mINFO\x1b[0m\tfoobar\tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
			},
		},
		{"set-env-development",
			func(config *Configurator) {
				os.Setenv("ENVIRONMENT", "development")
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.True(t, config.Development)
				require.Contains(t, output[0], "\t\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}")
				require.Contains(t, output[1], "\t\x1b[34mINFO\x1b[0m\tfoobar\tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
			},
		},
		{"set-invalid-preset",
			func(config *Configurator) {
				config.SetPreset("INVALID")
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.Error(t, err)
				require.Equal(t, `unknown preset: "INVALID"`, err.Error())
			},
		},
		{"set-preset-light-console",
			func(config *Configurator) {
				config.SetPreset("light-console")
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.Equal(t, "\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}", output[0])
				require.Equal(t, "\x1b[34mINFO \x1b[0m\tfoobar            \tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}", output[1])
			},
		},
		{"set-preset-light-json",
			func(config *Configurator) {
				config.SetPreset("light-json")
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.Equal(t, `{"level":"debug","caller":"zapconfig/zapconfig_test.go:16","msg":"hello world!","foo":"bar"}`, output[0])
				require.Equal(t, `{"level":"info","logger":"foobar","caller":"zapconfig/zapconfig_test.go:17","msg":"hello world!","baz":42}`, output[1])
			},
		},
		{"set-preset-console",
			func(config *Configurator) {
				config.SetPreset("console")
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.Contains(t, output[0], "\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}")
				require.Contains(t, output[1], "\x1b[34mINFO \x1b[0m\tfoobar            \tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
				require.NotEqual(t, output[0], "\x1b[35mDEBUG\x1b[0m\tzapconfig/zapconfig_test.go:16\thello world!\t{\"foo\": \"bar\"}")
				require.NotEqual(t, output[1], "\x1b[34mINFO \x1b[0m\tfoobar            \tzapconfig/zapconfig_test.go:17\thello world!\t{\"baz\": 42}")
			},
		},
		{"set-preset-json",
			func(config *Configurator) {
				config.SetPreset("json")
			},
			func(t *testing.T, config *zap.Config, output []string, err error) {
				require.NoError(t, err)
				require.Contains(t, output[0], `{"level":"debug","ts":`)
				require.Contains(t, output[0], `,"caller":"zapconfig/zapconfig_test.go:16","msg":"hello world!","foo":"bar"}`)
				require.Contains(t, output[1], `{"level":"info","ts":`)
				require.Contains(t, output[1], `,"logger":"foobar","caller":"zapconfig/zapconfig_test.go:17","msg":"hello world!","baz":42}`)
				require.NotEqual(t, output[0], `{"level":"debug","caller":"zapconfig/zapconfig_test.go:16","msg":"hello world!","foo":"bar"}`)
				require.NotEqual(t, output[1], `{"level":"info","logger":"foobar","caller":"zapconfig/zapconfig_test.go:17","msg":"hello world!","baz":42}`)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			configurator := Configurator{}
			if tc.setupFn != nil {
				tc.setupFn(&configurator)
			}
			config, err := configurator.Config()
			if err == nil {
			}

			var output []string
			if err == nil {
				closer, err := u.CaptureStdoutAndStderr()
				require.NoError(t, err)

				logger, err := config.Build()
				require.NoError(t, err)
				require.NotNil(t, logger)

				logFoobar(t, logger)

				logger.Sync()
				output = strings.Split(closer(), "\n")
				require.Equal(t, output[len(output)-1], "") // last line is always empty
				output = output[:len(output)-1]
			}
			tc.checkFn(t, &config, output, err)
		})
	}
}
