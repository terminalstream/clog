package clog

import "go.uber.org/zap/zapcore"

type entryFieldCallbacks struct {
	zapcore.Core
	cbs     []func(zapcore.Entry, []zapcore.Field)
	context []zapcore.Field
}

func (c *entryFieldCallbacks) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}

	return checked
}

func (c *entryFieldCallbacks) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	allFields := append(c.context, fields...)

	for i := range c.cbs {
		c.cbs[i](entry, allFields)
	}

	return c.Core.Write(entry, fields)
}

func (c *entryFieldCallbacks) With(fields []zapcore.Field) zapcore.Core {
	return &entryFieldCallbacks{
		Core:    c.Core.With(fields),
		cbs:     c.cbs,
		context: append(c.context, fields...),
	}
}
