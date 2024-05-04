# zerozap
An implementation of [zapcore.Core](https://pkg.go.dev/go.uber.org/zap/zapcore#Core)
to feed all data to a [zerolog](https://github.com/rs/zerolog) logger.

## Basic example

```go
package main

import (
	"os"

	"github.com/rs/zerolog"
	"go.mau.fi/zerozap"
	"go.uber.org/zap"
)

func main() {
	mainLog := zerolog.New(os.Stdout)
	zapLogger := zap.New(zerozap.New(mainLog))
	zapLogger.Info("Hello, world!")
}
```
