package clog

import "go.uber.org/zap/zapcore"

type entryFieldCallbacks struct {
	zapcore.Core
	cbs []func(zapcore.Entry, []zapcore.Field)
}

func (c *entryFieldCallbacks) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}

	return checked
}

func (c *entryFieldCallbacks) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	err := c.Core.Write(entry, fields)

	for _, cb := range c.cbs {
		cb(entry, fields)
	}

	return err
}

func (c *entryFieldCallbacks) With(fields []zapcore.Field) zapcore.Core {
	c.Core = c.Core.With(fields)

	return c
}
