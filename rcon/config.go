package rcon

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type config struct {
	Servers        map[string]server
	Default_Server string
}

type server struct {
	Host          string
	RCON_Password string
}

var YamlConfig config
var yamlLocation string = os.Getenv("HOME") + "/.config/.mcrcon"

// will be concatenated with .yml or .yaml depending on the config extension.
const yamlName = "mc-rcon"

// file path to the yaml config
var yamlFilePath string = fmt.Sprintf("%s/%s.yaml", yamlLocation, yamlName)

// Initializes the YAML struct data and validates the configuration directory and file.
func init() {
	// TODO: remove the panics in the future.
	err := checkYamlPath()
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(yamlLocation, 0644)
	if err != nil {
		panic(err)
	}

	yamlContent, err := os.ReadFile(yamlFilePath)
	if err != nil {
		// TODO: fix later, i am not sure what error can go here.
		panic(err)
	}

	yaml.Unmarshal(yamlContent, &YamlConfig)
}

// AddServerEntry adds a server entry to the YAML file with a given name.
// A name is required to distinguish the entry from other entries. If an
// entry exists with the same name, then it will get replaced by the new entry.
//
// An error is returned if there are issues interacting with the YAML file.
func AddServerEntry(name string, host string, password string) error {
	newEntry := server{
		Host:          host,
		RCON_Password: password,
	}

	YamlConfig.Servers[name] = newEntry

	err := checkYamlPath()
	if err != nil {
		return err
	}

	yamlContent, err := yaml.Marshal(YamlConfig)
	if err != nil {
		return err
	}

	err = os.WriteFile(yamlFilePath, yamlContent, 0644)
	if err != nil {
		// TODO: fix later, not sure what error can occur here due to the check above.
		return err
	}

	return nil
}

// checkYamlPath checks if the YAML path is valid.
// It checks for the file ending in .yaml or .yml.
// It will update the yamlFilePath to end in .yml if a .yml file exists.
//
// It returns nil if there is no error with the path.
// Otherwise an error will return if it fails.
func checkYamlPath() error {
	yamlExtensions := []string{"yaml", "yml"}
	// will contain an err if the loop does not return early
	var err error = nil

	for _, ext := range yamlExtensions {
		newYamlPath := fmt.Sprintf("%s/%s.%s", yamlLocation, yamlName, ext)

		_, err = os.Stat(newYamlPath)
		if err == nil {
			yamlFilePath = newYamlPath
			return nil
		}
	}

	return err
}
