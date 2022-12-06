package adapter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/iKuiki/go-component/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

// zapAdapter zap的封装实现
type zapAdapter struct {
	zapLogger *zap.SugaredLogger
	config    zapConfig

	// ========== 内部结构 ==========

	// 黑名单context
	blockContextMap map[string]bool
}

// zapOutPutConfig zap log config
type zapOutPutConfig struct {
	LogLevel         string   `json:"logLevel" yaml:"logLevel"`
	FileName         string   `json:"fileName" yaml:"fileName"`
	TimeFormat       string   `json:"timeFormat" yaml:"timeFormat"`
	TimeKey          string   `json:"timeKey" yaml:"timeKey"`
	LevelKey         string   `json:"levelKey" yaml:"levelKey"`
	CallerKey        string   `json:"callerKey" yaml:"callerKey"`
	MessageKey       string   `json:"messageKey" yaml:"messageKey"`
	BlockContextKeys []string `json:"blockContextKeys" yaml:"blockContextKeys,flow"`
}

// rotateConfig rotate config for zap
type rotateConfig struct {
	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxSize" yaml:"maxSize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxAge" yaml:"maxAge"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxBackups" yaml:"maxBackups"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress"`
}

// zapConfig full zap config
type zapConfig struct {
	Zap    zapOutPutConfig `json:"zap" yaml:"zap"`
	Rotate rotateConfig    `json:"rotate" yaml:"rotate"`
}

// NewZaplogAdapterFromFile 通过配置文件创建zap logger
func NewZaplogAdapterFromFile(logConfigFile, buildVersion string) (logger.Logger, error) {
	config, err := loadConfig(logConfigFile)
	if err != nil {
		return nil, err
	}

	lo := createRotateZapLogger(config)

	log := lo.WithOptions(zap.AddCallerSkip(1)).With(createInitialLogFields(buildVersion)...).Sugar()
	z := zapAdapter{
		zapLogger: log,
		config:    config,
	}
	z.blockContextMap = make(map[string]bool)
	for _, k := range config.Zap.BlockContextKeys {
		z.blockContextMap[k] = true
	}
	return &z, nil
}

// createEncoderConfig 从配置初始化zap EncoderConfig
func createEncoderConfig(config zapConfig) zapcore.EncoderConfig {
	// 创建日志参数配置对象
	encoderConfig := zap.NewProductionEncoderConfig()

	// 日志格式设置
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(config.Zap.TimeFormat))
	}
	encoderConfig.MessageKey = config.Zap.MessageKey
	encoderConfig.CallerKey = config.Zap.CallerKey
	encoderConfig.LevelKey = config.Zap.LevelKey
	encoderConfig.TimeKey = config.Zap.TimeKey

	return encoderConfig
}

// createRotateZapLogger 创建rotate日志zap Logger
func createRotateZapLogger(config zapConfig) *zap.Logger {
	encoderConfig := createEncoderConfig(config)

	level := zap.InfoLevel
	level.UnmarshalText([]byte(config.Zap.LogLevel))

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.Zap.FileName,
		MaxSize:    config.Rotate.MaxSize, // megabytes
		MaxBackups: config.Rotate.MaxBackups,
		MaxAge:     config.Rotate.MaxAge, //days
		Compress:   false,                // disabled by default
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		level,
	)
	return zap.New(core, zap.AddCaller())
}

// createNormalZapLogger 创建单文件日志zap Logger
func createNormalZapLogger(config zapConfig) (*zap.Logger, error) {
	// 创建日志参数配置对象
	cfg := zap.NewProductionConfig()

	// 日志输出文件
	cfg.OutputPaths = []string{config.Zap.FileName}

	// 日志格式设置
	encoderConfig := createEncoderConfig(config)
	cfg.EncoderConfig = encoderConfig
	cfg.DisableStacktrace = true

	cfg.Level.UnmarshalText([]byte(config.Zap.LogLevel))
	lo, err := cfg.Build()
	if err != nil {
		err = errors.Wrapf(err, "zap cfg.Build() error")
		return nil, err
	}
	return lo, err
}

var defaultZapConfig zapConfig

// init 初始化缺省defaultZapConfig. 特意设置缺省为测试状态的配置，以方便区分是否从配置文件读取的
func init() {
	defaultZapConfig = zapConfig{
		Zap: zapOutPutConfig{
			LogLevel:         "info",
			FileName:         fmt.Sprintf("./zap-%d.log", time.Now().Unix()),
			TimeFormat:       "2006-01-02 15:04:05.000000",
			TimeKey:          "time",
			LevelKey:         "level",
			CallerKey:        "caller",
			MessageKey:       "msg",
			BlockContextKeys: []string{},
		},
		Rotate: rotateConfig{
			MaxSize:    20,
			MaxAge:     30,
			MaxBackups: 50,
			Compress:   false,
		},
	}
}

// loadConfig 从配置文件读取日志配置
func loadConfig(configFile string) (retConfig zapConfig, err error) {
	retConfig = defaultZapConfig
	yamlFile, err := ioutil.ReadFile(configFile)

	if err != nil {
		err = errors.Wrapf(err, "ioutil.ReadFile(%s)", configFile)
		return
	}
	err = yaml.Unmarshal(yamlFile, &retConfig)
	if err != nil {
		err = errors.Wrapf(err, "Zap ConfigFile [%s] yaml.Unmarshal error", configFile)
		return
	}
	return retConfig, nil
}

// createInitialLogFields 创建每条日志都需要的host/project/version字段
func createInitialLogFields(buildVersion string) (fields []zap.Field) {
	host, _ := os.Hostname()

	version := buildVersion
	project := filepath.Base(os.Args[0])

	return []zap.Field{
		zap.String("host", host),
		zap.String("project", project),
		zap.String("version", version),
	}
}

func (adapter *zapAdapter) Print(v ...interface{}) {
	adapter.zapLogger.Info(v...)
}

func (adapter *zapAdapter) Printf(format string, args ...interface{}) {
	adapter.zapLogger.Infof(format, args...)
}

func (adapter *zapAdapter) Println(v ...interface{}) {
	adapter.zapLogger.Info(v...)
}

func (adapter *zapAdapter) Error(v ...interface{}) {
	adapter.zapLogger.Error(v...)
}

func (adapter *zapAdapter) Errorf(format string, args ...interface{}) {
	adapter.zapLogger.Errorf(format, args...)
}

func (adapter *zapAdapter) Warn(v ...interface{}) {
	adapter.zapLogger.Warn(v...)
}

func (adapter *zapAdapter) Warnf(format string, args ...interface{}) {
	adapter.zapLogger.Warnf(format, args...)
}

func (adapter *zapAdapter) Info(v ...interface{}) {
	adapter.zapLogger.Info(v...)
}

func (adapter *zapAdapter) Infof(format string, args ...interface{}) {
	adapter.zapLogger.Infof(format, args...)
}

func (adapter *zapAdapter) Debug(v ...interface{}) {
	adapter.zapLogger.Debug(v...)
}

func (adapter *zapAdapter) Debugf(format string, args ...interface{}) {
	adapter.zapLogger.Debugf(format, args...)
}

func (adapter *zapAdapter) Fatal(v ...interface{}) {
	adapter.zapLogger.Fatal(v...)
	panic(errors.Errorf("Fatal: %v", v))
}

func (adapter *zapAdapter) Fatalf(format string, args ...interface{}) {
	adapter.zapLogger.Fatalf(format, args...)
	panic(errors.Errorf(format, args...))
}

func (adapter *zapAdapter) Sync() {
	adapter.zapLogger.Sync()
}

func (adapter *zapAdapter) With(values logger.Values) logger.Logger {
	fields := []interface{}{}
	var overwriteCaller bool
	if values != nil {
		values.RangeValues(func(key, value string) bool {
			if _, ok := adapter.blockContextMap[key]; !ok && value != "" {
				// 不在黑名单，则输出
				fields = append(fields, zap.String(key, value))
			}
			return true
		})
		// 固定字段
		if str := values.GetString("caller"); str != "" {
			// 重写了caller
			overwriteCaller = true
			// fields = append(fields, zap.String("caller", str))
		}
	}
	if len(fields) > 0 {
		cfg := adapter.config
		zLog := adapter.zapLogger
		if overwriteCaller {
			// 增加多一层offset
			zLog = zLog.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar()
		}
		ret := &zapAdapter{
			zapLogger: zLog.With(fields...),
			config:    cfg,
		}
		return ret
	}
	return adapter
}
