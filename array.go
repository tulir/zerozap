// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package zerozap

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap/zapcore"
)

type zeroArray struct {
	evt *zerolog.Array
}

var _ zapcore.ArrayEncoder = (*zeroArray)(nil)

type arrayProxy struct {
	arr zapcore.ArrayMarshaler
	err error
}

var _ zerolog.LogArrayMarshaler = (*arrayProxy)(nil)

func (ap *arrayProxy) MarshalZerologArray(arr *zerolog.Array) {
	ap.err = ap.arr.MarshalLogArray(&zeroArray{evt: arr})
}

func (z *zeroArray) AppendArray(marshaler zapcore.ArrayMarshaler) error {
	// TODO why does zerolog not support nested arrays?
	//ap := &arrayProxy{arr: marshaler}
	//z.evt.Array(key, ap)
	//return ap.err
	return fmt.Errorf("zerolog doesn't support nested arrays")
}

func (z *zeroArray) AppendObject(marshaler zapcore.ObjectMarshaler) error {
	op := &objectProxy{obj: marshaler}
	z.evt.Object(op)
	return op.err
}

func (z *zeroArray) AppendBinary(value []byte) {
	z.evt.Str(base64.StdEncoding.EncodeToString(value))
}

func (z *zeroArray) AppendByteString(value []byte) {
	z.evt.Bytes(value)
}

func (z *zeroArray) AppendBool(value bool) {
	z.evt.Bool(value)
}

func (z *zeroArray) AppendComplex128(value complex128) {
	z.evt.Str(strconv.FormatComplex(value, 'f', -1, 128))
}

func (z *zeroArray) AppendComplex64(value complex64) {
	z.evt.Str(strconv.FormatComplex(complex128(value), 'f', -1, 64))
}

func (z *zeroArray) AppendDuration(value time.Duration) {
	z.evt.Dur(value)
}

func (z *zeroArray) AppendFloat64(value float64) {
	z.evt.Float64(value)
}

func (z *zeroArray) AppendFloat32(value float32) {
	z.evt.Float32(value)
}

func (z *zeroArray) AppendInt(value int) {
	z.evt.Int(value)
}

func (z *zeroArray) AppendInt64(value int64) {
	z.evt.Int64(value)
}

func (z *zeroArray) AppendInt32(value int32) {
	z.evt.Int32(value)
}

func (z *zeroArray) AppendInt16(value int16) {
	z.evt.Int16(value)
}

func (z *zeroArray) AppendInt8(value int8) {
	z.evt.Int8(value)
}

func (z *zeroArray) AppendString(value string) {
	z.evt.Str(value)
}

func (z *zeroArray) AppendTime(value time.Time) {
	z.evt.Time(value)
}

func (z *zeroArray) AppendUint(value uint) {
	z.evt.Uint(value)
}

func (z *zeroArray) AppendUint64(value uint64) {
	z.evt.Uint64(value)
}

func (z *zeroArray) AppendUint32(value uint32) {
	z.evt.Uint32(value)
}

func (z *zeroArray) AppendUint16(value uint16) {
	z.evt.Uint16(value)
}

func (z *zeroArray) AppendUint8(value uint8) {
	z.evt.Uint8(value)
}

func (z *zeroArray) AppendUintptr(value uintptr) {
	z.evt.Uint(uint(value))
}

func (z *zeroArray) AppendReflected(value interface{}) error {
	z.evt.Interface(value)
	return nil
}
