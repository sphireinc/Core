package core

import (
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
	"os"
	"path/filepath"
)

type Handler []routing.Handler

var (
	MethodGet    = []string{"GET"}
	MethodPut    = []string{"PUT"}
	MethodPost   = []string{"POST"}
	MethodDelete = []string{"DELETE"}
)

func EmptyJson() map[string]interface{} {
	return map[string]interface{}{}
}

func (c *Config) Methods() []string {
	return []string{"GET,PUT,POST,DELETE"}
}

func (c *Config) Method(method string) []string {
	return []string{method}
}

func (c *Config) M(method string) []string {
	return c.Method(method)
}

func findFileInDirectory(directoryPath, filename string) (string, error) {
	// Check if the directory exists
	fileInfo, err := os.Stat(directoryPath)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return "", fmt.Errorf("%s is not a directory", directoryPath)
	}

	// Walk through the directory
	var targetFile string
	err = filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == filename {
			targetFile = path
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if targetFile == "" {
		return "", fmt.Errorf("file %s not found in directory %s", filename, directoryPath)
	}

	return targetFile, nil
}

// configLocations returns a slice of possible locations for the configuration file
func configLocations() []string {
	var locations []string

	// Add current directory
	locations = append(locations, ".")

	// Add home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		locations = append(locations, homeDir)
	}

	// Add common system-wide locations
	commonLocations := []string{
		"/etc",
		"/usr/local/etc",
	}
	for _, loc := range commonLocations {
		locations = append(locations, loc)
	}

	// Add executable directory
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		locations = append(locations, exeDir)
	}

	return locations
}

// findConfigFile searches for the configuration file in the specified locations
func findConfigFile(locations []string, configFilename string) (string, error) {
	for _, loc := range locations {
		configFile := filepath.Join(loc, configFilename)
		_, err := os.Stat(configFile)
		if err == nil {
			return configFile, nil
		}
	}
	return "", fmt.Errorf("configuration file '%s' not found in any of the locations", configFilename)
}
