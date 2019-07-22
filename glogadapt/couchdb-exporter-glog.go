// Adapts parts of the golang/glog package
// to allow us passing glog flags via our own flags package.
// Idea taken from https://github.com/kubernetes/kubernetes/pull/3342

package glogadapt

import (
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/golang/glog"
)

// severity identifies the sort of log: info, warning etc. It also implements
// the flag.Value interface. The -stderrthreshold flag is of type severity and
// should be modified only through the flag.Value interface. The values match
// the corresponding constants in C++.
type severity int32 // sync/atomic int32

// These constants identify the log levels in order of increasing severity.
// A message written to a high-severity log file is also written to each
// lower-severity log file.
const (
	infoLog severity = iota
	warningLog
	errorLog
	fatalLog
)

var severityName = []string{
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

// get returns the value of the severity.
func (s *severity) get() severity {
	return severity(atomic.LoadInt32((*int32)(s)))
}

// set sets the value of the severity.
func (s *severity) set(val severity) {
	atomic.StoreInt32((*int32)(s), int32(val))
}

// String is part of the flag.Value interface.
func (s *severity) String() string {
	return strconv.FormatInt(int64(*s), 10)
}

// Set is part of the flag.Value interface.
func (s *severity) Set(value string) error {
	var threshold severity
	// Is it a known name?
	if v, ok := severityByName(value); ok {
		threshold = v
	} else {
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		threshold = severity(v)
	}
	Logging.StderrThreshold.set(threshold)
	return nil
}

func severityByName(s string) (severity, bool) {
	s = strings.ToUpper(s)
	for i, name := range severityName {
		if name == s {
			return severity(i), true
		}
	}
	return 0, false
}

// loggingT collects all the global state of the logging setup.
type LoggingT struct {
	// Boolean flags. Not handled atomically because the flag.Value interface
	// does not let us avoid the =true, and that shorthand is necessary for
	// compatibility. TODO: does this matter enough to fix? Seems unlikely.
	ToStderr     bool // The -logtostderr flag.
	AlsoToStderr bool // The -alsologtostderr flag.

	// Level flag. Handled atomically.
	StderrThreshold severity // The -stderrthreshold flag.

	Verbosity glog.Level // V logging level, the value of the -v flag/

	// If non-empty, overrides the choice of directory in which to write logs.
	// See glog.createLogDirs for the full list of possible destinations.
	LogDir string
}

var Logging LoggingT

func init() {
	// Default stderrThreshold is ERROR.
	Logging.StderrThreshold = errorLog
}
