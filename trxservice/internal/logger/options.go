package logger

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Option custom setup config
type Option func(*option)

type option struct {
	level          zapcore.Level
	fields         map[string]string
	file           io.Writer
	timeLayout     string
	disableConsole bool
}


func defaultOption() *option {
	return &option{
		level:         DefaultLevel,
		fields:        make(map[string]string),
		file:           &lumberjack.Logger{ // concurrent-safed
			Filename:   "./logs/app.log", // 文件路径
			MaxSize:    50,  // 单个文件最大尺寸，默认单位 M
			MaxBackups: 30,  // 最多保留 30 个备份
			MaxAge:     30,   // 最大时间30天，默认单位 day
			LocalTime:  true, // 使用本地时间
			Compress:   true, // 是否压缩 disabled by default
		},
		timeLayout:     DefaultTimeLayout,
		disableConsole: false,
	}
}

func generateOption(opts ...Option) *option {
	config := defaultOption()
	for _, opt := range opts {
		opt(config)
	}
	return config
}



func WithLogLevel(level string) Option {
	return func(opt *option) {
		switch strings.ToLower(level) {
		case "debug":
			opt.level = zapcore.DebugLevel
		case "info":
			opt.level = zapcore.InfoLevel
		case "warn":
			opt.level = zapcore.WarnLevel
		case "error":
			opt.level = zapcore.ErrorLevel
		case "fatal":
			opt.level = zapcore.FatalLevel
		case "panic":
			opt.level = zapcore.PanicLevel
		default:
			opt.level = zapcore.InfoLevel
		}
	}
}


// WithDebugLevel only greater than 'level' will output
func WithDebugLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.DebugLevel
	}
}

// WithInfoLevel only greater than 'level' will output
func WithInfoLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.InfoLevel
	}
}

// WithWarnLevel only greater than 'level' will output
func WithWarnLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.WarnLevel
	}
}

// WithErrorLevel only greater than 'level' will output
func WithErrorLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.ErrorLevel
	}
}

// WithField add some field(s) to log
func WithField(key, value string) Option {
	return func(opt *option) {
		opt.fields[key] = value
	}
}

// WithFileP write log to some file
func WithFileP(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0766)
	if err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = zapcore.Lock(f)
	}
}

// WithFileRotationP write log to some file with rotation
func WithFileRotationP(file string, maxSize, maxBackups, maxAge int) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = &lumberjack.Logger{ // concurrent-safed
			Filename:   file, // 文件路径
			MaxSize:    maxSize,  // 单个文件最大尺寸，默认单位 M
			MaxBackups: maxBackups,  // 最多保留 300 个备份
			MaxAge:     maxAge,   // 最大时间，默认单位 day
			LocalTime:  true, // 使用本地时间
			Compress:   true, // 是否压缩 disabled by default
		}
	}
}

// WithTimeLayout custom time format
func WithTimeLayout(timeLayout string) Option {
	return func(opt *option) {
		opt.timeLayout = timeLayout
	}
}

// WithDisableConsole WithEnableConsole write log to os.Stdout or os.Stderr
func WithDisableConsole() Option {
	return func(opt *option) {
		opt.disableConsole = true
	}
}
