// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package zerozap_test

import (
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"go.uber.org/zap"

	"go.mau.fi/zerozap"
)

func init() {
	zerozap.CopyTime = false
}

func ExampleNew() {
	zlog := zerolog.New(os.Stdout)

	zaplog := zap.New(zerozap.New(zlog))
	zaplog.Info("Hello, world!")
	// Output: {"level":"info","message":"Hello, world!"}
}

func ExampleOption() {
	zlog := zerolog.New(os.Stdout)

	zaplog := zap.NewNop()
	zaplog.Info("This is not logged")

	zaplog = zaplog.WithOptions(zerozap.Option(zlog))
	zaplog.Info("This is logged")
	// Output: {"level":"info","message":"This is logged"}
}

func TestZeroZap(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		fn       func(*zap.Logger)
	}{
		{
			name: "Generic",
			expected: `{"level":"info","message":"Hello, world!"}
{"level":"info","int":42,"str":"meow","true":false,"message":"Normal fields"}
{"level":"info","meow??":{"subfield":1,"meow!!":{"subsubfield":2}},"message":"Namespaced fields"}
{"level":"info","meow":["me","o","w"],"message":"Array"}
`,
			fn: func(logger *zap.Logger) {
				logger.Info("Hello, world!")
				logger.Info("Normal fields", zap.Int("int", 42), zap.String("str", "meow"), zap.Bool("true", false))
				logger.Info("Namespaced fields", zap.Namespace("meow??"), zap.Int("subfield", 1), zap.Namespace("meow!!"), zap.Int("subsubfield", 2))
				logger.Info("Array", zap.Strings("meow", []string{"me", "o", "w"}))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf strings.Builder
			test.fn(zap.New(&zerozap.ZeroZap{Logger: zerolog.New(&buf)}))
			if out := buf.String(); out != test.expected {
				t.Errorf("expected:\n%s\ngot:\n%s", test.expected, out)
			}
		})
	}
}
