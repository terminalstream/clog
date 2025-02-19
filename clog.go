// Copyright 2024 Terminal Stream Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clog

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level represents a log level.
type Level int8

func (l Level) String() string {
	return zapcore.Level(l).String()
}

const (
	// DefaultLevel is the default logging level.
	DefaultLevel = InfoLevel
	// DebugLevel represents the DEBUG level.
	DebugLevel = Level(zapcore.DebugLevel)
	// InfoLevel represents the INFO level.
	InfoLevel = Level(zapcore.InfoLevel)
	// WarnLevel represents the WARN level.
	WarnLevel = Level(zapcore.WarnLevel)
	// ErrorLevel represents the ERROR level.
	ErrorLevel = Level(zapcore.ErrorLevel)
	// PanicLevel represents the PANIC level.
	PanicLevel = Level(zapcore.PanicLevel)
)

const (
	// DefaultMessageKey is the default key that has the log message as value.
	DefaultMessageKey = "msg"
	// DefaultEncoding is the default logging encoding. Possible values are "console" or "json".
	DefaultEncoding = "console"
	// DefaultLevelKey is the default key that has the LogLevel value.
	DefaultLevelKey = "severity"
	// DefaultErrorKey is the default key that has as value any errors logged with WithError.
	DefaultErrorKey = "error"
	// DefaultTimeKey is the default key that holds the timestamp.
	DefaultTimeKey = "time"
)

// Fields is a collection of key-values that you can provide to a logging context or to a log
// record.
type Fields map[string]any

type logKeyType string

var (
	loggerKey logKeyType = "logger"
	levelKey  logKeyType = "level_key"
	errorKey  logKeyType = "error_key"
)

// Option allows extending individual log records with additional structured data.
type Option func(*options)

type options struct {
	err    error
	fields map[string]any
}

// WithError adds an error field to the log record.
func WithError(err error) Option {
	return func(o *options) {
		o.err = err
	}
}

// WithField adds a field to the log record.
func WithField(key string, value any) Option {
	return func(o *options) {
		if o.fields == nil {
			o.fields = make(Fields)
		}

		o.fields[key] = value
	}
}

// WithFields adds multiple fields to the log record.
func WithFields(fields Fields) Option {
	return func(o *options) {
		if o.fields == nil {
			o.fields = make(Fields)
		}

		for k, v := range fields {
			o.fields[k] = v
		}
	}
}

// ContextOption allows customization of a few aspects of a logging context.
type ContextOption func(*contextOptions)

type contextOptions struct {
	encoding            string
	level               Level
	outputPath          string
	levelKey            string
	msgKey              string
	timeKey             string
	errorKey            string
	entryFieldCallbacks []func(zapcore.Entry, []zapcore.Field)
}

// WithLevel lets the logging context's Level to level. InfoLevel is the default Level.
func WithLevel(level Level) ContextOption {
	return func(o *contextOptions) {
		o.level = level
	}
}

// WithJSONEncoding sets the logging format to JSON. 'Console' format is the default format.
func WithJSONEncoding() ContextOption {
	return func(o *contextOptions) {
		o.encoding = "json"
	}
}

// WithConsoleEncoding sets the logging format to 'console' (this is the default anyway).
func WithConsoleEncoding() ContextOption {
	return func(o *contextOptions) {
		o.encoding = "console"
	}
}

// OutputToStdout redirects logging output to os.Stdout (default is os.Stderr).
func OutputToStdout() ContextOption {
	return func(o *contextOptions) {
		o.outputPath = "stdout"
	}
}

// WithLevelKey allows switching away from the DefaultLevelKey.
func WithLevelKey(key string) ContextOption {
	return func(o *contextOptions) {
		o.levelKey = key
	}
}

// WithMessageKey allows switching away from the DefaultMessageKey.
func WithMessageKey(key string) ContextOption {
	return func(o *contextOptions) {
		o.msgKey = key
	}
}

// WithTimeKey allows switching away from the DefaultTimeKey.
func WithTimeKey(key string) ContextOption {
	return func(o *contextOptions) {
		o.timeKey = key
	}
}

// WithNoTimeKey disables timestamps in log messages.
func WithNoTimeKey() ContextOption {
	return func(o *contextOptions) {
		o.timeKey = ""
	}
}

// WithErrorKey allows switching away from the DefaultErrorKey.
func WithErrorKey(key string) ContextOption {
	return func(o *contextOptions) {
		o.errorKey = key
	}
}

func WithEntryFieldCallbacks(cbs ...func(zapcore.Entry, []zapcore.Field)) ContextOption {
	return func(o *contextOptions) {
		o.entryFieldCallbacks = append(o.entryFieldCallbacks, cbs...)
	}
}

// ParseLevel parses the given level.
func ParseLevel(level string) (Level, error) {
	l, err := zapcore.ParseLevel(level)
	if err != nil {
		return InfoLevel, fmt.Errorf("invalid level: %w", err)
	}

	return Level(l), nil
}

// Context returns a new contextual logging context.
//
// The returned context is a child of parent unless parent is nil; in that case the returned
// context is derived from the background context.
//
// It is important to obtain a logging context with this function first before invoking any
// of the rest. Not doing so renders all other functions as no-ops.
func Context(parent context.Context, opts ...ContextOption) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	o := &contextOptions{
		encoding:   DefaultEncoding,
		level:      DefaultLevel,
		levelKey:   DefaultLevelKey,
		msgKey:     DefaultMessageKey,
		timeKey:    DefaultTimeKey,
		errorKey:   DefaultErrorKey,
		outputPath: "stderr",
	}

	for i := range opts {
		opts[i](o)
	}

	level := zap.NewAtomicLevelAt(zapcore.Level(o.level))

	zapConfig := zap.Config{
		Level:             level,
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoding:          o.encoding,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  o.msgKey,
			LevelKey:    o.levelKey,
			TimeKey:     o.timeKey,
			EncodeTime:  zapcore.RFC3339TimeEncoder,
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
		OutputPaths: []string{o.outputPath},
	}

	logger := zap.Must(zapConfig.Build())

	if len(o.entryFieldCallbacks) > 0 {
		logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return &entryFieldCallbacks{
				Core: core,
				cbs:  o.entryFieldCallbacks,
			}
		}))
	}

	return context.WithValue(
		context.WithValue(
			context.WithValue(parent, loggerKey, logger),
			levelKey,
			&level,
		),
		errorKey,
		o.errorKey,
	)
}

// CopyContext copies the logging context from 'from' into a new context derived from 'to'.
//
// This is a no-op if 'from' is not a logging context ('to' is returned as-is).
func CopyContext(to, from context.Context) context.Context {
	logger, ok := from.Value(loggerKey).(*zap.Logger)
	if !ok {
		return to
	}

	return context.WithValue(to, loggerKey, logger)
}

// ContextWithField returns a new logging context derived from parent and including
// the given key and value.
//
// If parent is not a logging context then parent is returned as-is.
func ContextWithField(parent context.Context, k string, v any) context.Context {
	logger, ok := parent.Value(loggerKey).(*zap.Logger)
	if !ok {
		return parent
	}

	logger = logger.With(zap.Any(k, v))

	return context.WithValue(parent, loggerKey, logger)
}

// ContextWithFields returns a new logging context derived from parent and including
// the given keys and values.
//
// If parent is not a logging context then parent is returned as-is.
func ContextWithFields(parent context.Context, fields Fields) context.Context {
	logger, ok := parent.Value(loggerKey).(*zap.Logger)
	if !ok {
		return parent
	}

	zf := make([]zap.Field, 0, len(fields))

	for k, v := range fields {
		zf = append(zf, zap.Any(k, v))
	}

	logger = logger.With(zf...)

	return context.WithValue(parent, loggerKey, logger)
}

// SetLevel adjusts the logging level on the given logging context.
//
// If 'ctx' is not a logging context then this is a no-op.
func SetLevel(ctx context.Context, level Level) {
	l, ok := ctx.Value(levelKey).(*zap.AtomicLevel)
	if !ok {
		return
	}

	l.SetLevel(zapcore.Level(level))
}

// DebugEnabled indicates whether DebugLevel is enabled on the given context.
//
// If ctx is not a logging context then false is returned.
func DebugEnabled(ctx context.Context) bool {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return false
	}

	return logger.Level().Enabled(zapcore.DebugLevel)
}

// Debug will log at the DebugLevel.
func Debug(ctx context.Context, msg string, opts ...Option) {
	if !DebugEnabled(ctx) {
		return
	}

	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return
	}

	logger.Debug(msg, getFields(ctx, opts)...)
}

// InfoEnabled indicates whether InfoLevel is enabled on the given context.
//
// If ctx is not a logging context then false is returned.
func InfoEnabled(ctx context.Context) bool {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return false
	}

	return logger.Level().Enabled(zapcore.InfoLevel)
}

// Info logs at the InfoLevel.
func Info(ctx context.Context, msg string, opts ...Option) {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return
	}

	if !logger.Level().Enabled(zapcore.InfoLevel) {
		return
	}

	logger.Info(msg, getFields(ctx, opts)...)
}

// WarnEnabled indicates whether WarnLevel is enabled on the given context.
//
// If ctx is not a logging context then false is returned.
func WarnEnabled(ctx context.Context) bool {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return false
	}

	return logger.Level().Enabled(zapcore.WarnLevel)
}

// Warn logs at the WarnLevel.
func Warn(ctx context.Context, msg string, opts ...Option) {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return
	}

	if !logger.Level().Enabled(zapcore.WarnLevel) {
		return
	}

	logger.Warn(msg, getFields(ctx, opts)...)
}

// ErrorEnabled indicates whether ErrorLevel is enabled on the given context.
//
// If ctx is not a logging context then false is returned.
func ErrorEnabled(ctx context.Context) bool {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return false
	}

	return logger.Level().Enabled(zapcore.ErrorLevel)
}

// Error logs at the ErrorLevel.
func Error(ctx context.Context, msg string, opts ...Option) {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return
	}

	if !logger.Level().Enabled(zapcore.ErrorLevel) {
		return
	}

	logger.Error(msg, getFields(ctx, opts)...)
}

// Panic logs at the PanicLevel.
func Panic(ctx context.Context, msg string, opts ...Option) {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return
	}

	if !logger.Level().Enabled(zapcore.PanicLevel) {
		return
	}

	logger.Panic(msg, getFields(ctx, opts)...)
}

func getFields(ctx context.Context, opts []Option) []zap.Field {
	o := &options{}

	for i := range opts {
		opts[i](o)
	}

	zf := make([]zap.Field, 0, len(o.fields))

	for k, v := range o.fields {
		zf = append(zf, zap.Any(k, v))
	}

	if o.err != nil {
		errKey, ok := ctx.Value(errorKey).(string)
		if ok {
			zf = append(zf, zap.NamedError(errKey, o.err))
		} else {
			zf = append(zf, zap.NamedError(DefaultErrorKey, o.err))
		}
	}

	return zf
}
