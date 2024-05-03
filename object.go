// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package zerozap

import (
	"encoding/base64"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap/zapcore"
)

type zeroObject struct {
	evt    *zerolog.Event
	finish func()
}

var _ zapcore.ObjectEncoder = (*zeroObject)(nil)

type objectProxy struct {
	obj zapcore.ObjectMarshaler
	err error
}

var _ zerolog.LogObjectMarshaler = (*objectProxy)(nil)

func (op *objectProxy) MarshalZerologObject(evt *zerolog.Event) {
	zo := &zeroObject{evt: evt}
	op.err = op.obj.MarshalLogObject(zo)
	if zo.finish != nil {
		zo.finish()
	}
}

func (z *zeroObject) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	ap := &arrayProxy{arr: marshaler}
	z.evt.Array(key, ap)
	return ap.err
}

func (z *zeroObject) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	op := &objectProxy{obj: marshaler}
	z.evt.Object(key, op)
	return op.err
}

func (z *zeroObject) AddBinary(key string, value []byte) {
	z.evt.Str(key, base64.StdEncoding.EncodeToString(value))
}

func (z *zeroObject) AddByteString(key string, value []byte) {
	z.evt.Bytes(key, value)
}

func (z *zeroObject) AddBool(key string, value bool) {
	z.evt.Bool(key, value)
}

func (z *zeroObject) AddComplex128(key string, value complex128) {
	z.evt.Str(key, strconv.FormatComplex(value, 'f', -1, 128))
}

func (z *zeroObject) AddComplex64(key string, value complex64) {
	z.evt.Str(key, strconv.FormatComplex(complex128(value), 'f', -1, 64))
}

func (z *zeroObject) AddDuration(key string, value time.Duration) {
	z.evt.Dur(key, value)
}

func (z *zeroObject) AddFloat64(key string, value float64) {
	z.evt.Float64(key, value)
}

func (z *zeroObject) AddFloat32(key string, value float32) {
	z.evt.Float32(key, value)
}

func (z *zeroObject) AddInt(key string, value int) {
	z.evt.Int(key, value)
}

func (z *zeroObject) AddInt64(key string, value int64) {
	z.evt.Int64(key, value)
}

func (z *zeroObject) AddInt32(key string, value int32) {
	z.evt.Int32(key, value)
}

func (z *zeroObject) AddInt16(key string, value int16) {
	z.evt.Int16(key, value)
}

func (z *zeroObject) AddInt8(key string, value int8) {
	z.evt.Int8(key, value)
}

func (z *zeroObject) AddString(key, value string) {
	z.evt.Str(key, value)
}

func (z *zeroObject) AddTime(key string, value time.Time) {
	z.evt.Time(key, value)
}

func (z *zeroObject) AddUint(key string, value uint) {
	z.evt.Uint(key, value)
}

func (z *zeroObject) AddUint64(key string, value uint64) {
	z.evt.Uint64(key, value)
}

func (z *zeroObject) AddUint32(key string, value uint32) {
	z.evt.Uint32(key, value)
}

func (z *zeroObject) AddUint16(key string, value uint16) {
	z.evt.Uint16(key, value)
}

func (z *zeroObject) AddUint8(key string, value uint8) {
	z.evt.Uint8(key, value)
}

func (z *zeroObject) AddUintptr(key string, value uintptr) {
	z.evt.Uint(key, uint(value))
}

func (z *zeroObject) AddReflected(key string, value interface{}) error {
	z.evt.Any(key, value)
	return nil
}

func (z *zeroObject) OpenNamespace(key string) {
	parentEvt := z.evt
	prevFinish := z.finish
	subEvt := zerolog.Dict()

	z.evt = subEvt
	z.finish = func() {
		parentEvt.Dict(key, subEvt)
		if prevFinish != nil {
			prevFinish()
		}
	}
}
