// Copyright 2025 Terminal Stream Inc.
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

import "go.uber.org/zap/zapcore"

type hooksLogger struct {
	zapcore.Core
	hooks   []func(zapcore.Entry, []zapcore.Field)
	context []zapcore.Field // https://github.com/terminalstream/clog/issues/3
}

func (c *hooksLogger) Check(
	entry zapcore.Entry, checked *zapcore.CheckedEntry,
) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}

	return checked
}

func (c *hooksLogger) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	allFields := append(c.context, fields...)

	for i := range c.hooks {
		c.hooks[i](entry, allFields)
	}

	return c.Core.Write(entry, fields)
}

func (c *hooksLogger) With(fields []zapcore.Field) zapcore.Core {
	return &hooksLogger{
		Core:    c.Core.With(fields),
		hooks:   c.hooks,
		context: append(c.context, fields...),
	}
}
