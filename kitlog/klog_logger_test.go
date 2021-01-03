package kitlog_test

import (
	"bytes"
	"errors"
	"flag"
	"strconv"
	"strings"
	"testing"

	"github.com/gesellix/couchdb-prometheus-exporter/v30/kitlog"
	"k8s.io/klog/v2"
)

func TestKlogLogger(t *testing.T) {
	t.Parallel()

	klogFlags := flag.NewFlagSet("klog", flag.ContinueOnError)
	klog.InitFlags(klogFlags)
	klogFlags.Set("skip_headers", strconv.FormatBool(true))
	klog.LogToStderr(false)
	buf := new(bytes.Buffer)
	klog.SetOutput(buf)
	logger := kitlog.NewKlogLogger()

	if err := logger.Log("hello", "world"); err != nil {
		t.Fatal(err)
	}
	if want, have := "hello=\"world\"\n", strings.Split(buf.String(), " ")[1]; want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	buf.Reset()
	if err := logger.Log("a", 1, "err", errors.New("error")); err != nil {
		t.Fatal(err)
	}
	if want, have := "a=1 err=\"error\"", strings.TrimSpace(strings.SplitAfterN(buf.String(), " ", 2)[1]); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	buf.Reset()
	if err := logger.Log("a", 1, "b"); err != nil {
		t.Fatal(err)
	}
	if want, have := "a=1 b=\"(MISSING)\"", strings.TrimSpace(strings.SplitAfterN(buf.String(), " ", 2)[1]); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	buf.Reset()
	if err := logger.Log("my_map", mymap{0: 0}); err != nil {
		t.Fatal(err)
	}
	if want, have := "my_map=\"special_behavior\"", strings.TrimSpace(strings.Split(buf.String(), " ")[1]); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}

type mymap map[int]int

func (m mymap) String() string { return "special_behavior" }
