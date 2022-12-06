package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func TestDefaultZapConfig(t *testing.T) {
	// 允许并行测试
	t.Parallel()
	assert.Equal(t, defaultZapConfig.Zap.LogLevel, "info")
	assert.Equal(t, defaultZapConfig.Zap.MessageKey, "msg")
	bytes, err := yaml.Marshal(defaultZapConfig)
	t.Logf("yaml.Marshal return: [%s]", string(bytes))
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
}

func TestLoadConfig(t *testing.T) {
	// 允许并行测试
	t.Parallel()
	config, err := loadConfig("./zaplog_test_config.yaml")
	assert.NoError(t, err)
	t.Logf("loadConfig return: [%v]", config)
	assert.Equal(t, config.Zap.LogLevel, "debug")
	assert.Equal(t, config.Zap.MessageKey, "message")
	assert.Equal(t, config.Zap.FileName, "./main.json.log")

	assert.Equal(t, 30, config.Rotate.MaxSize)
	assert.Equal(t, 30, config.Rotate.MaxAge)
	assert.Equal(t, 30, config.Rotate.MaxBackups)
}

func TestCreateInitialLogFields(t *testing.T) {
	// 允许并行测试
	t.Parallel()
	tag := "Test2020123"
	fields := createInitialLogFields(tag)
	t.Logf("createInitialLogFields(%s) return:[%v]", tag, fields)
	assert.True(t, len(fields) >= 3)
	for _, field := range fields {
		if field.Key == "version" {
			assert.Equal(t, tag, field.String)
		} else {
			assert.True(t, len(field.String) > 0)
		}
	}
}

func TestCreateNormalZapLogger(t *testing.T) {
	// 允许并行测试
	t.Parallel()
	config := defaultZapConfig
	config.Zap.LogLevel = "debug"
	logger, err := createNormalZapLogger(config)
	t.Logf("createNormalZapLogger(%v) return: [%+v]", config, logger)
	assert.Nil(t, err)
	logger.Debug("Without dynamic context")
	logger.With(zap.String("request_id", "test_request")).Debug("With dynamic context")
}

func TestCreateRotateZapLogger(t *testing.T) {
	// 允许并行测试
	t.Parallel()
	config := defaultZapConfig
	config.Zap.LogLevel = "debug"
	logger := createRotateZapLogger(config)
	t.Logf("createRotateZapLogger(%v) return: [%+v]", config, logger)

	logger.Debug("Without dynamic context")
	logger.With(zap.String("request_id", "test_request_for_rotate")).Debug("With dynamic context")
}

func TestNewZaplogAdapterFromFile(t *testing.T) {
	// 允许并行测试
	t.Parallel()
	version := "TestNewZaplogAdapterFromFile"
	configFile := "./zaplog_test_config.yaml"

	logger, err := NewZaplogAdapterFromFile(configFile, version)
	assert.NoError(t, err)
	t.Logf("NewZaplogAdapterFromFile(%s, %s) return: [%+v]", configFile, version, logger)

	logger.Debug("Debug:Without dynamic context")
	logger.Debugf("Debugf:Without dynamic context, version: %s", version)
	logger.Info("Info:Without dynamic context")
	logger.Infof("Infof:Without dynamic context, version: %s", version)
	logger.Warn("Warn:Without dynamic context")
	logger.Warnf("Warnf:Without dynamic context, version: %s", version)
	logger.Error("Error:Without dynamic context")
	logger.Errorf("Errorf:Without dynamic context, version: %s", version)
}
