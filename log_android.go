package main

/*
#cgo LDFLAGS: -landroid -llog

#include <android/log.h>
#include <string.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/sirupsen/logrus"
)

var levels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
}

type androidHook struct {
	tag *C.char
	fmt logrus.Formatter
}

type androidFormatter struct{}

func (f *androidFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

func (hook *androidHook) Levels() []logrus.Level {
	return levels
}

func (hook *androidHook) Fire(e *logrus.Entry) error {
	var priority C.int

	formatted, err := hook.fmt.Format(e)
	if err != nil {
		return err
	}
	str := C.CString(string(formatted))

	switch e.Level {
	case logrus.PanicLevel:
		priority = C.ANDROID_LOG_FATAL
	case logrus.FatalLevel:
		priority = C.ANDROID_LOG_FATAL
	case logrus.ErrorLevel:
		priority = C.ANDROID_LOG_ERROR
	case logrus.WarnLevel:
		priority = C.ANDROID_LOG_WARN
	case logrus.InfoLevel:
		priority = C.ANDROID_LOG_INFO
	case logrus.DebugLevel:
		priority = C.ANDROID_LOG_DEBUG
	}
	C.__android_log_write(priority, hook.tag, str)
	C.free(unsafe.Pointer(str))
	return nil
}

// create a logrus Hook that forward entries to logcat
func AndroidLogHook(tag string) logrus.Hook {
	return &androidHook{
		tag: C.CString(tag),
		fmt: &androidFormatter{},
	}
}

func addAndroidLogHook() {
	logrus.AddHook(AndroidLogHook("Torrent-Go"))
}
