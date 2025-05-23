package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())

}

type DiscardHandler struct {}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (dh *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (dh *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (dh *DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return dh
}

func (dh *DiscardHandler) WithGroup(name string) slog.Handler {
	return dh
}
