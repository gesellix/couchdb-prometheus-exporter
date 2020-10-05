// Heavily copied and inspired by urfave/cli and https://stackoverflow.com/a/46860900/372019

package fileutil

import (
	"fmt"
	"testing"
)

func TestReadPropertiesFile(t *testing.T) {
	props, err := readPropertiesFile("test.properties")
	if err != nil {
		t.Error("Error while reading properties file")
		t.FailNow()
	}

	if props["host"] != "localhost" || props["proxyHost"] != "test" || props["protocol"] != "https://" || props["chunk"] != "" || props["boolean"] != "true" {
		fmt.Printf("props: %q", props)
		t.Error("Error properties not loaded correctly")
		t.Fail()
	}
}

func TestReadBoolPropertyFromFile(t *testing.T) {
	props, err := NewPropertiesSourceFromFile("test.properties")
	if err != nil {
		t.Error("Error while reading properties file")
		t.FailNow()
	}

	if b, _ := props.Bool("boolean"); !b {
		fmt.Printf("props: %q", props)
		t.Error("Error properties not loaded correctly")
	}
}
