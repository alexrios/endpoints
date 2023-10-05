package main

import (
	"path/filepath"
	"syscall"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFirstRun(t *testing.T) {
	t.Run("is the first time running", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		firstRun, err := isFirstRun(fs, DefaultConfigurationFileName)

		assert.Nil(t, err)
		assert.True(t, firstRun)
	})
	t.Run("is the first time running given a config path", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		firstRun, err := isFirstRun(fs, filepath.Join("newpath", DefaultConfigurationFileName))

		assert.Nil(t, err)
		assert.True(t, firstRun)
	})
	t.Run("is NOT the first time running", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, _ = fs.Create(DefaultConfigurationFileName)
		firstRun, err := isFirstRun(fs, DefaultConfigurationFileName)

		assert.Nil(t, err)
		assert.False(t, firstRun)
	})
	t.Run("is NOT the first time running given a config path", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		path := filepath.Join("newpath", DefaultConfigurationFileName)
		_, _ = fs.Create(path)
		firstRun, err := isFirstRun(fs, path)

		assert.Nil(t, err)
		assert.False(t, firstRun)
	})
}

func TestConfigureFirstRun(t *testing.T) {
	t.Run("configure first time use files", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		err := configureFirstRun(fs, DefaultConfigurationFileName)

		assert.Nil(t, err)
	})

	t.Run("configure first time use files given a config path", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		err := configureFirstRun(fs, filepath.Join("newpath", DefaultConfigurationFileName))

		assert.Nil(t, err)
	})

	t.Run("cannot configure when is a read only dir", func(t *testing.T) {
		fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
		err := configureFirstRun(fs, DefaultConfigurationFileName)
		if assert.Error(t, err) {
			assert.Equal(t, syscall.Errno(1), err)
		}
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("loading config file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_ = configureFirstRun(fs, DefaultConfigurationFileName)

		config, err := loadConfig(fs, DefaultConfigurationFileName)
		assert.Nil(t, err)
		assert.Equal(t, ":8080", config.Addr)
		assert.Equal(t, 1, len(config.Responses))
		assert.Equal(t, 200, config.Responses[0].Status)
		assert.Equal(t, "GET", config.Responses[0].Method)
		assert.Equal(t, "0ms", config.Responses[0].Latency)
	})
	t.Run("loading config file", func(t *testing.T) {
		path := filepath.Join("newpath", DefaultConfigurationFileName)
		fs := afero.NewMemMapFs()
		_ = configureFirstRun(fs, path)

		config, err := loadConfig(fs, path)
		assert.Nil(t, err)
		assert.Equal(t, ":8080", config.Addr)
		assert.Equal(t, 1, len(config.Responses))
		assert.Equal(t, 200, config.Responses[0].Status)
		assert.Equal(t, "GET", config.Responses[0].Method)
		assert.Equal(t, "0ms", config.Responses[0].Latency)
	})
}
