package cache

import (
	"github.com/linmadan/egglib-go/log"
	"time"
)

type Options struct {
	CleanInterval time.Duration
	DebugMode     bool
	Log           func() log.Logger
}

type Option func(o *Options)

func NewOptions(option ...Option) *Options {
	o := &Options{}
	for i := range option {
		option[i](o)
	}
	return o
}

func WithCleanInterval(i time.Duration) Option {
	return func(o *Options) {
		o.CleanInterval = i
	}
}

func WithDebugLog(DebugModule bool, log func() log.Logger) Option {
	return func(o *Options) {
		o.DebugMode = DebugModule
		o.Log = log
	}
}
