// Copyright 2022 NetEase Media Technology（Beijing）Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xzap

import (
	"strings"
	"sync"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _textPool = sync.Pool{New: func() interface{} {
	return &textEncoder{}
}}

func getTextEncoder() *textEncoder {
	return _textPool.Get().(*textEncoder)
}

func putTextEncoder(enc *textEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.openNamespaces = 0
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	_textPool.Put(enc)
}

type textEncoder struct {
	*jsonEncoder
	opt *Option
}

// NewTextEncoder creates a fast, low-allocation Text encoder. The encoder
// appropriately escapes all field keys and values.
//
// Note that the encoder doesn't deduplicate keys, so it's possible to produce
// a message like
//
//	{"foo":"bar","foo":"baz"}
//
// This is permitted by the Text specification, but not encouraged. Many
// libraries will ignore duplicate key-value pairs (typically keeping the last
// pair) when unmarshaling, but users should attempt to avoid adding duplicate
// keys.
func NewTextEncoder(opt *Option, cfg zapcore.EncoderConfig) zapcore.Encoder {
	return newTextEncoder(opt, cfg)
}

func newTextEncoder(opt *Option, cfg zapcore.EncoderConfig) *textEncoder {
	return &textEncoder{
		opt:         opt,
		jsonEncoder: newJSONEncoder(opt, cfg, false),
	}
}

func (enc *textEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *textEncoder) clone() *textEncoder {
	clone := getTextEncoder()
	clone.jsonEncoder = enc.jsonEncoder.clone()
	clone.opt = enc.opt
	return clone
}

func (enc *textEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	var pkg, f string
	if ent.Caller.Defined {
		i := strings.LastIndex(ent.Caller.Function, ".")
		pkg = ent.Caller.Function[:i]
		if level, ok := enc.opt.packageLogLevel[pkg]; ok && level > ent.Level {
			return emptyBuffer, nil
		}
		f = ent.Caller.Function[i+1:]
	}

	final := enc.clone()
	if final.TimeKey != "" {
		cur := final.buf.Len()
		if e := final.EncodeTime; e != nil {
			e(ent.Time, final)
		}
		if cur == final.buf.Len() {
			// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
			// output JSON valid.
			final.AppendInt64(ent.Time.UnixNano())
		}
		final.AppendString(" ")
	}

	if final.LevelKey != "" && final.EncodeLevel != nil {
		final.AppendString("[")
		cur := final.buf.Len()
		final.EncodeLevel(ent.Level, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeLevel was a no-op. Fall back to strings to keep
			// output Text valid.

			final.AppendString(ent.Level.String())
		}
		final.AppendString("] ")
	}

	if ent.LoggerName != "" && final.NameKey != "" {
		cur := final.buf.Len()
		nameEncoder := final.EncodeName

		// if no name encoder provided, fall back to FullNameEncoder for backwards
		// compatibility
		if nameEncoder == nil {
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(ent.LoggerName, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeName was a no-op. Fall back to strings to
			// keep output Text valid.
			final.AppendString(ent.LoggerName)
		}
		final.AppendString(" ")
	}
	if ent.Caller.Defined {
		if final.CallerKey != "" {
			final.AppendString("[")
			cur := final.buf.Len()
			final.EncodeCaller(ent.Caller, final)
			if cur == final.buf.Len() {
				// User-supplied EncodeCaller was a no-op. Fall back to strings to
				// keep output Text valid.
				final.AppendString(ent.Caller.String())
			}
			final.AppendString("] ")
		}
		if final.FunctionKey != "" {
			final.AppendString("[")
			final.AppendString(f)
			final.AppendString("] ")
		}
	}
	if final.MessageKey != "" {
		final.AppendString(ent.Message)
	}
	fs := len(fields) > 0
	buf := enc.buf.Len() > 0
	if fs || buf {
		final.AppendString(" {")
		if buf {
			final.buf.Write(enc.buf.Bytes())
		}
		if fs {
			if buf {
				final.addElementSeparator()
			}
			addFields(final.jsonEncoder, fields)
		}
		final.AppendString("}")
	}

	if ent.Stack != "" && final.StacktraceKey != "" {
		final.buf.AppendByte('\n')
		final.buf.AppendString(ent.Stack)
	}
	final.buf.AppendString(final.LineEnding)

	ret := final.buf
	putTextEncoder(final)
	return ret, nil
}

func (enc *textEncoder) AppendString(val string) {
	enc.safeAddString(val)
}

func (enc *textEncoder) AppendInt64(val int64) {
	enc.buf.AppendInt(val)
}

func (enc *textEncoder) AppendTimeLayout(time time.Time, layout string) {
	enc.buf.AppendTime(time, layout)
}
