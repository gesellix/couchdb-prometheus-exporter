// Heavily copied and inspired by urfave/cli and https://stackoverflow.com/a/46860900/372019

package fileutil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func readPropertiesFile(filename string) (map[interface{}]interface{}, error) {
	config := map[interface{}]interface{}{}

	if len(filename) == 0 {
		return config, nil
	}

	raw, err := loadDataFrom(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(raw))
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Other file formats would provide more type information, and unmarshalling would
	// create correct types. Properties files aren't that useful, so we have to enforce
	// some types by ourselves. We have to keep an eye on other types beside boolean,
	// but this one seems good enough for now.
	coerceBooleanValues(config)

	return config, nil
}

func coerceBooleanValues(m map[interface{}]interface{}) {
	for k, v := range m {
		switch v := v.(type) {
		case string:
			parsed, err := strconv.ParseBool(v)
			if err == nil {
				m[k] = parsed
			}
		}
	}
}

type propertiesSourceContext struct {
	FilePath string
}

// NewPropertiesSourceFromFile creates a new Properties InputSourceContext from a filepath.
func NewPropertiesSourceFromFile(file string) (altsrc.InputSourceContext, error) {
	ysc := &propertiesSourceContext{FilePath: file}
	var results map[interface{}]interface{}
	results, err := readPropertiesFile(ysc.FilePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to load Properties file '%s': inner error: \n'%v'", ysc.FilePath, err.Error())
	}

	return altsrc.NewMapInputSource(file, results), nil
}

// NewPropertiesSourceFromFlagFunc creates a new Properties InputSourceContext from a provided flag name and source context.
func NewPropertiesSourceFromFlagFunc(flagFileName string) func(context *cli.Context) (altsrc.InputSourceContext, error) {
	return func(context *cli.Context) (altsrc.InputSourceContext, error) {
		filePath := context.String(flagFileName)
		return NewPropertiesSourceFromFile(filePath)
	}
}

func loadDataFrom(filePath string) ([]byte, error) {
	u, err := url.Parse(filePath)
	if err != nil {
		return nil, err
	}

	if u.Host != "" { // i have a host, now do i support the scheme?
		switch u.Scheme {
		case "http", "https":
			res, err := http.Get(filePath)
			if err != nil {
				return nil, err
			}
			return io.ReadAll(res.Body)
		default:
			return nil, fmt.Errorf("scheme of %s is unsupported", filePath)
		}
	} else if u.Path != "" { // I don't have a host, but I have a path. I am a local file.
		if _, notFoundFileErr := os.Stat(filePath); notFoundFileErr != nil {
			return nil, fmt.Errorf("cannot read from file: '%s' because it does not exist", filePath)
		}
		return os.ReadFile(filePath)
	} else if runtime.GOOS == "windows" && strings.Contains(u.String(), "\\") {
		// on Windows systems u.Path is always empty, so we need to check the string directly.
		if _, notFoundFileErr := os.Stat(filePath); notFoundFileErr != nil {
			return nil, fmt.Errorf("cannot read from file: '%s' because it does not exist", filePath)
		}
		return os.ReadFile(filePath)
	} else {
		return nil, fmt.Errorf("unable to determine how to load from path %s", filePath)
	}
}
