package rcon

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// TODO: add logging to Config, i have not set it up yet!

type Config struct {
	Servers       map[string]Server `yaml:"servers"`
	DefaultServer string            `yaml:"default_server"`
	yamlDirectory string
	yamlFilePath  string
}

type Server struct {
	Host     string `yaml:"host"`
	Password string `yaml:"rcon_password"`
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
		DefaultServer: "",
		yamlDirectory: yamlDir,
		yamlFilePath:  yamlPath,
	}

	err := yamlConfig.checkMakeYaml()
	if err != nil {
		err := yamlConfig.createNewYaml()
		if err != nil {
			return nil, err
		}
	}

	//  when unmarshaling an empty config, it overwrites the paths- these holds the paths beforehand.
	yamlDirState := yamlConfig.yamlDirectory
	yamlFileState := yamlConfig.yamlFilePath

	err = os.MkdirAll(yamlDir, 0o755)
	if err != nil {
		return nil, err
	}

	yamlContent, err := os.ReadFile(yamlConfig.yamlFilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &yamlConfig)
	if err != nil {
		return nil, err
	}

	if yamlConfig.yamlDirectory == "" {
		yamlConfig.yamlDirectory = yamlDirState
	}
	if yamlConfig.yamlFilePath == "" {
		yamlConfig.yamlFilePath = yamlFileState
	}

	return &yamlConfig, nil
}

// AddServerEntry adds a server entry to the YAML file with a given name.
// A tag is required to distinguish the entry from other entries. If an
// entry exists with the given tag, then the previous tag will be replaced by the new tag.
//
// An error is returned if there are issues interacting with the YAML file.
func (c *Config) AddServerEntry(serverTag string, host string, password string) error {
	newEntry := Server{
		Host:     host,
		Password: password,
	}

	if c.Servers == nil {
		c.Servers = make(map[string]Server)
	}

	c.Servers[serverTag] = newEntry

	err := c.checkMakeYaml()
	if err != nil {
		err := c.createNewYaml()
		if err != nil {
			return err
		}
	}

	err = c.updateYaml()
	if err != nil {
		return err
	}

	return nil
}

// RemoveServerEntry takes a server tag and removes it from the YAML file.
//
// If no entries exist then nothing occurs.
// If the file fails to update then an error is returned.
func (c *Config) RemoveServerEntry(serverTag string) error {
	if _, ok := c.Servers[serverTag]; !ok {
		fmt.Printf("no entries found with tag %s\n", serverTag)
		return nil
	} else {
		delete(c.Servers, serverTag)
	}

	err := c.updateYaml()
	if err != nil {
		return err
	}

	return nil
}

// GetServer returns the server of the matching server tag in the entries.
//
// If the name entry does not exist, then return an error.
func (c *Config) GetServerEntry(serverTag string) (*Server, error) {
	if s, ok := c.Servers[serverTag]; ok {
		return &s, nil
	}

	errMsg := fmt.Sprintf("no entries found with tag %s", serverTag)
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

// checkMakeYaml checks if the YAML exists, if not then it creates an empty YAML file.
// c.yamlFilePath will be updated to the matching YAML file if the extension
// does not end in .yaml.
//
// It returns nil if the file exists or if the file is successfully created.
// If the file fails to write then it will return an error.
func (c *Config) checkMakeYaml() error {
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

	return fmt.Errorf("no such file, path: %s", c.yamlFilePath)
}

// createNewYaml creates a new empty YAML file.
func (c *Config) createNewYaml() error {
	err := os.WriteFile(c.yamlFilePath, []byte{}, 0o644)
	if err != nil {
		return err
	}

	return nil
}
