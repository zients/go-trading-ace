package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerConfigRequestTimeoutDefaultsToTenSeconds(t *testing.T) {
	cfg := ServerConfig{}

	assert.Equal(t, 10*time.Second, cfg.RequestTimeout())
}

func TestServerConfigRequestTimeoutUsesConfiguredSeconds(t *testing.T) {
	cfg := ServerConfig{RequestTimeoutSeconds: 3}

	assert.Equal(t, 3*time.Second, cfg.RequestTimeout())
}

func TestLoadConfigDefaultsInvalidRequestTimeoutToTenSeconds(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(`
server:
  port: 8080
  request_timeout_seconds: "invalid"
`), 0600)
	require.NoError(t, err)

	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(dir)
	v.SetConfigType("yaml")

	cfg, err := loadConfig(v)

	require.NoError(t, err)
	assert.Equal(t, 10*time.Second, cfg.Server.RequestTimeout())
}
