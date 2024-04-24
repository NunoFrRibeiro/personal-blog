package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var logger *slog.Logger

func init() {

  replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey  {
      a.Key = "date"
      a.Value = slog.Int64Value(time.Now().Unix())
		}
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace}))
}

func Infof(format string, args ...any) {
  var pcs [1]uintptr
  runtime.Callers(2, pcs[:])
  r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(format, args...), pcs[0])
  _ = logger.Handler().Handle(context.Background(), r)
}

func Warnf(format string, args ...any) {
  var pcs [1]uintptr
  runtime.Callers(2, pcs[:])
  r := slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf(format, args...), pcs[0])
  _ = logger.Handler().Handle(context.Background(), r)
}

func Errorf(format string, args ...any) {
  var pcs [1]uintptr
  runtime.Callers(2, pcs[:])
  r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(format, args...), pcs[0])
  _ = logger.Handler().Handle(context.Background(), r)
}
