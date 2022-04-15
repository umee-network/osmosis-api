package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/umee-network/osmosis-api/config"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		name      string
		cfg       config.Config
		expectErr bool
	}{
		{
			"valid config",
			config.Config{
				Server: config.Server{
					ListenAddr: "0.0.0.0:7171",
				},
			},
			false,
		},
		{
			"empty server",
			config.Config{
				Server: config.Server{
					ListenAddr: "",
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.cfg.Validate() != nil, tc.expectErr)
		})
	}
}

func TestParseConfig_Valid(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "config.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte(`
[server]
listen_addr = "0.0.0.0:99999"
read_timeout = "20s"
`)
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	cfg, err := config.ParseConfig(tmpFile.Name())
	require.NoError(t, err)

	require.Equal(t, "0.0.0.0:99999", cfg.Server.ListenAddr)
	require.Equal(t, "20s", cfg.Server.WriteTimeout)
	require.Equal(t, "20s", cfg.Server.ReadTimeout)
}
