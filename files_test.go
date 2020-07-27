package main

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"syscall"
	"testing"
)

func TestFirstRun(t *testing.T) {
	t.Run("is the first time running", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		firstRun, err := isFirstRun(fs)

		assert.Nil(t, err)
		assert.True(t, firstRun)
	})
	t.Run("is NOT the first time running", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, _ = fs.Create(DefaultConfigurationFileName)
		firstRun, err := isFirstRun(fs)

		assert.Nil(t, err)
		assert.False(t, firstRun)
	})
}

func TestConfigureFirstRun(t *testing.T) {
	t.Run("configure first time use files", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		err := configureFirstRun(fs)

		assert.Nil(t, err)
	})

	t.Run("cannot configure when is a read only dir", func(t *testing.T) {
		fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
		err := configureFirstRun(fs)
		if assert.Error(t, err) {
			assert.Equal(t, syscall.Errno(1), err)
		}
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("loading config file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_ = configureFirstRun(fs)

		config, err := loadConfig(fs)
		assert.Nil(t, err)
		assert.Equal(t, ":8080", config.Addr)
		assert.Equal(t, 1, len(config.Responses))
		assert.Equal(t, 200, config.Responses[0].Status)
		assert.Equal(t, "GET", config.Responses[0].Method)
		assert.Equal(t, "0ms", config.Responses[0].Latency)
	})
}
