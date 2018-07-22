package read

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// Type RawScript representing .sql script. It contains script's name, it's path
// and script's content as a single string.
// RawScript is non-transformed raw .sql script type.
type RawScript struct {
	Name     string
	FullPath string
	Content  string
}

// Function SQLScript reads .sql script as a single string.
// In case of error while reading the file function return also
// non nil error value.
func SQLScript(path string) (RawScript, error) {
	file, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		log.Printf("Could not read .sql script from path = [%s] \n", path)
		return RawScript{}, readErr
	}

	scriptName, nameErr := parseScriptName(path)
	if nameErr != nil {
		log.Printf("Could not parse script name from path = [%s] \n", path)
		return RawScript{}, nameErr
	}

	fullPath, pathErr := filepath.Abs(path)
	if pathErr != nil {
		log.Printf("Could not parse full path script from path = [%s] \n", path)
		return RawScript{}, pathErr
	}

	return RawScript{scriptName, fullPath, string(file)}, nil
}

// This function parse script name from it's path
func parseScriptName(path string) (string, error) {
	pathParts := strings.Split(path, ".")
	if len(pathParts) == 0 {
		return "", errors.New("Path does not contain a dot.")
	}

	return filepath.Base(path), nil
}
