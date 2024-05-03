// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package zerozap

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap/zapcore"
)

var levelMap = map[zapcore.Level]zerolog.Level{
	zapcore.DebugLevel:  zerolog.DebugLevel,
	zapcore.InfoLevel:   zerolog.InfoLevel,
	zapcore.WarnLevel:   zerolog.WarnLevel,
	zapcore.ErrorLevel:  zerolog.ErrorLevel,
	zapcore.DPanicLevel: zerolog.PanicLevel,
	zapcore.PanicLevel:  zerolog.PanicLevel,
	zapcore.FatalLevel:  zerolog.FatalLevel,
}

type ZeroZap struct {
	zerolog.Logger
}

var _ zapcore.Core = (*ZeroZap)(nil)

func (z *ZeroZap) Enabled(level zapcore.Level) bool {
	return z.GetLevel() <= levelMap[level]
}

func (z *ZeroZap) With(fields []zapcore.Field) zapcore.Core {
	logWith := z.Logger.With()
	for _, f := range fields {
		switch f.Type {
		case zapcore.ArrayMarshalerType:
			ap := &arrayProxy{arr: f.Interface.(zapcore.ArrayMarshaler)}
			logWith = logWith.Array(f.Key, ap)
			// TODO why doesn't this function return errors when AddObject and others do?
			if ap.err != nil {
				panic(ap.err)
			}
		case zapcore.ObjectMarshalerType:
			op := &objectProxy{obj: f.Interface.(zapcore.ObjectMarshaler)}
			logWith = logWith.Object(f.Key, op)
			if op.err != nil {
				panic(op.err)
			}
		case zapcore.InlineMarshalerType:
			op := &objectProxy{obj: f.Interface.(zapcore.ObjectMarshaler)}
			logWith = logWith.EmbedObject(op)
			if op.err != nil {
				panic(op.err)
			}
		case zapcore.BinaryType:
			logWith = logWith.Str(f.Key, base64.StdEncoding.EncodeToString(f.Interface.([]byte)))
		case zapcore.BoolType:
			logWith = logWith.Bool(f.Key, f.Integer == 1)
		case zapcore.ByteStringType:
			logWith = logWith.Bytes(f.Key, f.Interface.([]byte))
		case zapcore.Complex128Type:
			logWith = logWith.Str(f.Key, strconv.FormatComplex(f.Interface.(complex128), 'f', -1, 128))
		case zapcore.Complex64Type:
			logWith = logWith.Str(f.Key, strconv.FormatComplex(complex128(f.Interface.(complex64)), 'f', -1, 64))
		case zapcore.DurationType:
			logWith = logWith.Dur(f.Key, time.Duration(f.Integer))
		case zapcore.Float64Type:
			logWith = logWith.Float64(f.Key, math.Float64frombits(uint64(f.Integer)))
		case zapcore.Float32Type:
			logWith = logWith.Float32(f.Key, math.Float32frombits(uint32(f.Integer)))
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
			logWith = logWith.Int64(f.Key, f.Integer)
		case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type, zapcore.UintptrType:
			logWith = logWith.Uint64(f.Key, uint64(f.Integer))
		case zapcore.StringType:
			logWith = logWith.Str(f.Key, f.String)
		case zapcore.TimeType:
			if f.Interface != nil {
				logWith = logWith.Time(f.Key, time.Unix(0, f.Integer).In(f.Interface.(*time.Location)))
			} else {
				// Fall back to UTC if location is nil.
				logWith = logWith.Time(f.Key, time.Unix(0, f.Integer))
			}
		case zapcore.TimeFullType:
			logWith = logWith.Time(f.Key, f.Interface.(time.Time))
		case zapcore.ReflectType:
			logWith = logWith.Any(f.Key, f.Interface)
		case zapcore.NamespaceType:
			// TODO implement
			panic("unsupported field type namespace")
		case zapcore.StringerType:
			// TODO catch panics like zap does in encodeStringer?
			logWith = logWith.Stringer(f.Key, f.Interface.(fmt.Stringer))
		case zapcore.ErrorType:
			logWith = logWith.AnErr(f.Key, f.Interface.(error))
		case zapcore.SkipType:
			// noop
		default:
			panic(fmt.Sprintf("unknown field type: %v", f))
		}
	}
	return &ZeroZap{Logger: logWith.Logger()}
}

func (z *ZeroZap) Check(entry zapcore.Entry, entry2 *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	//TODO wtf is this supposed to do?
	panic("implement me")
}

func (z *ZeroZap) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	evt := z.Logger.WithLevel(levelMap[entry.Level])
	// TODO is this a good idea? it'll probably lead to the field being duplicated
	evt.Time(zerolog.TimestampFieldName, entry.Time)
	if entry.Stack != "" {
		evt.Str(zerolog.ErrorStackFieldName, entry.Stack)
	}
	if entry.Caller.Defined {
		evt.Str(zerolog.CallerFieldName, zerolog.CallerMarshalFunc(entry.Caller.PC, entry.Caller.File, entry.Caller.Line))
	}
	for _, f := range fields {
		switch f.Type {
		case zapcore.ArrayMarshalerType:
			ap := &arrayProxy{arr: f.Interface.(zapcore.ArrayMarshaler)}
			evt.Array(f.Key, ap)
			if ap.err != nil {
				return ap.err
			}
		case zapcore.ObjectMarshalerType:
			op := &objectProxy{obj: f.Interface.(zapcore.ObjectMarshaler)}
			evt.Object(f.Key, op)
			if op.err != nil {
				return op.err
			}
		case zapcore.InlineMarshalerType:
			op := &objectProxy{obj: f.Interface.(zapcore.ObjectMarshaler)}
			evt.EmbedObject(op)
			if op.err != nil {
				return op.err
			}
		case zapcore.BinaryType:
			evt.Str(f.Key, base64.StdEncoding.EncodeToString(f.Interface.([]byte)))
		case zapcore.BoolType:
			evt.Bool(f.Key, f.Integer == 1)
		case zapcore.ByteStringType:
			evt.Bytes(f.Key, f.Interface.([]byte))
		case zapcore.Complex128Type:
			evt.Str(f.Key, strconv.FormatComplex(f.Interface.(complex128), 'f', -1, 128))
		case zapcore.Complex64Type:
			evt.Str(f.Key, strconv.FormatComplex(complex128(f.Interface.(complex64)), 'f', -1, 64))
		case zapcore.DurationType:
			evt.Dur(f.Key, time.Duration(f.Integer))
		case zapcore.Float64Type:
			evt.Float64(f.Key, math.Float64frombits(uint64(f.Integer)))
		case zapcore.Float32Type:
			evt.Float32(f.Key, math.Float32frombits(uint32(f.Integer)))
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
			evt.Int64(f.Key, f.Integer)
		case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type, zapcore.UintptrType:
			evt.Uint64(f.Key, uint64(f.Integer))
		case zapcore.StringType:
			evt.Str(f.Key, f.String)
		case zapcore.TimeType:
			if f.Interface != nil {
				evt.Time(f.Key, time.Unix(0, f.Integer).In(f.Interface.(*time.Location)))
			} else {
				// Fall back to UTC if location is nil.
				evt.Time(f.Key, time.Unix(0, f.Integer))
			}
		case zapcore.TimeFullType:
			evt.Time(f.Key, f.Interface.(time.Time))
		case zapcore.ReflectType:
			evt.Any(f.Key, f.Interface)
		case zapcore.NamespaceType:
			// TODO implement
			return fmt.Errorf("unsupported field type namespace")
		case zapcore.StringerType:
			// TODO catch panics like zap does in encodeStringer?
			evt.Stringer(f.Key, f.Interface.(fmt.Stringer))
		case zapcore.ErrorType:
			evt.AnErr(f.Key, f.Interface.(error))
		case zapcore.SkipType:
			// noop
		default:
			return fmt.Errorf("unknown field type: %v", f)
		}
	}
	evt.Msg(entry.Message)
	return nil
}

func (z *ZeroZap) Sync() error {
	return nil
}
