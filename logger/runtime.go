package logger

import (
	"context"
	"fmt"
	"time"
)

type Metadata string

type Timer interface {
	Debug()
	Info()
	Warn()
	Error()
}

type timer struct {
	metadata Metadata
	t        time.Time
	ms       int64
	err      error
}

func (t timer) Debug() {
	if t.err == nil {
		Debug("%s use time %d ms", t.metadata, t.ms)
	}
}

func (t timer) Info() {
	if t.err == nil {
		Info("%s use time %d ms", t.metadata, t.ms)
	}
}

func (t timer) Warn() {
	if t.err == nil {
		Warn("%s use time %d ms", t.metadata, t.ms)
	}
}

func (t timer) Error() {
	if t.err == nil {
		Error("%s use time %d ms", t.metadata, t.ms)
	}
}

func Time(ctx context.Context, metadata Metadata) context.Context {
	return context.WithValue(context.TODO(), metadata, timer{
		metadata: metadata,
		t:        time.Now(),
	})
}

func TimeEnd(ctx context.Context, metadata Metadata) (t timer) {
	if meta, ok := ctx.Value(metadata).(timer); ok {
		t.metadata = meta.metadata
		t.t = meta.t
		t.ms = time.Since(meta.t).Milliseconds()
		return
	}
	t.err = fmt.Errorf("fail to get metadata %s use time", metadata)
	Error("fail to get metadata %s use time", metadata)
	return
}
