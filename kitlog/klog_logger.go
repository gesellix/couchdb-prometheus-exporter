// Package kitlog provides a klog adapter to the
// go-kit log.Logger interface.
package kitlog

import (
	kitlog "github.com/go-kit/kit/log"
	"k8s.io/klog/v2"
)

type klogLogger struct {
}

// NewKlogLogger returns a go-kit log.Logger that sends log events to a k8s.io/klog logger.
func NewKlogLogger() kitlog.Logger {
	return &klogLogger{}
}

// Log the keyvals to the INFO log.
func (l klogLogger) Log(keyvals ...interface{}) error {
	klog.InfoS("-", keyvals...)
	return nil
}
