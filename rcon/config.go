package rcon

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Servers        map[string]server
	Default_Server string
	yamlDirectory  string
	yamlFilePath   string
}

type server struct {
	Host          string
	RCON_Password string
}

const yamlName = "mc-rcon"

// NewConfig creates a new config object for YAML config by reading the
// YAML file.
//
// config is returned if the YAML can be successfully read.
// If the YAML file cannot be read, an error is returned and config will be nil.
func NewConfig(yamlDir string) (*Config, error) {
	// checkYamlPath validates this and changes it to a matching YAML extension.
	var yamlPath string = fmt.Sprintf("%s/%s.yaml", yamlDir, yamlName)
	yamlConfig := Config{
		yamlDirectory: yamlDir,
		yamlFilePath:  yamlPath,
	}

	err := yamlConfig.checkYamlPath()
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(yamlDir, 0o644)
	if err != nil {
		return nil, err
	}

	yamlContent, err := os.ReadFile(yamlConfig.yamlFilePath)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(yamlContent, &yamlConfig)
	return &yamlConfig, nil
}

// AddServerEntry adds a server entry to the YAML file with a given name.
// A name is required to distinguish the entry from other entries. If an
// entry exists with the same name, then it will get replaced by the new entry.
//
// An error is returned if there are issues interacting with the YAML file.
func (c *Config) AddServerEntry(name string, host string, password string) error {
	newEntry := server{
		Host:          host,
		RCON_Password: password,
	}

	c.Servers[name] = newEntry

	err := c.checkYamlPath()
	if err != nil {
		return err
	}

	err = c.updateYaml()
	if err != nil {
		return err
	}

	return nil
}

// RemoveServerEntry takes a server name and removes it from the YAML file.
//
// If no entries exist then nothing occurs.
// If the file fails to update then an error is returned.
func (c *Config) RemoveServerEntry(name string) error {
	if _, ok := c.Servers[name]; !ok {
		fmt.Printf("no entries found with name %s\n", name)
		return nil
	} else {
		delete(c.Servers, name)
	}

	err := c.updateYaml()
	if err != nil {
		return err
	}

	return nil
}

// GetServer returns the server of the given name.
//
// If the name entry does not exist, then return an error.
func (c *Config) GetServerEntry(name string) (*server, error) {
	if s, ok := c.Servers[name]; ok {
		return &s, nil
	}

	errMsg := fmt.Sprintf("no entries found with the name %s", name)
	return nil, errors.New(errMsg)
}

// GetYamlPath returns the full path to the YAML file.
func (c *Config) GetYamlPath() string {
	return c.yamlFilePath
}

// updateYaml updates the YAML file when a modification happens.
func (c *Config) updateYaml() error {
	yamlContent, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(c.yamlFilePath, yamlContent, 0o644)
	if err != nil {
		return err
	}

	return nil
}

// checkYamlPath checks if the YAML path is valid.
// It checks for the extension of the YAML file.
// It will update c.yamlFilePath to the matching file name, if found.
//
// It returns nil if the file exists.
// An error is return if the file does not exist.
func (c *Config) checkYamlPath() error {
	yamlExtensions := []string{"yaml", "yml", "YAML", "YML"}

	for _, ext := range yamlExtensions {
		// either no changes to the original path, or gets replaced- causing an update.
		newFileName := fmt.Sprintf("/%s.%s", yamlName, ext)
		newYamlPath := c.yamlDirectory + newFileName

		_, err := os.Stat(newYamlPath)
		if err == nil {
			c.yamlFilePath = newYamlPath
			return nil
		}
	}

	errMsg := fmt.Sprintf("no YAML config found in %s", c.yamlDirectory)
	return errors.New(errMsg)
}
