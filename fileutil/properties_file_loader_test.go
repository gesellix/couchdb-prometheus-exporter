// Heavily copied and inspired by urfave/cli and https://stackoverflow.com/a/46860900/372019

package fileutil

import (
	"testing"
)

func TestReadPropertiesFile(t *testing.T) {
	props, err := readPropertiesFile("test.properties")
	if err != nil {
		t.Error("Error while reading properties file")
	}

	if props["host"] != "localhost" || props["proxyHost"] != "test" || props["protocol"] != "https://" || props["chunk"] != "" {
		t.Error("Error properties not loaded correctly")
	}
}
