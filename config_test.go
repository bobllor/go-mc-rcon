package main

import (
	"fmt"
	"os"
	"rcon/rcon"
	"strings"
	"testing"
)

var testYaml string = `
servers:
  default:
    host: "some_ip:at_some_port"
    rcon_password: "go_test_password"
default_server: "default"
`

func TestCreateConfig(t *testing.T) {
	_, err := getConfig(t, testYaml)
	if err != nil {
		t.Error(err)
	}
}

func TestEmptyConfig(t *testing.T) {
	config, err := getConfig(t, "")
	if err != nil {
		t.Error(err)
	}

	if len(config.Servers) != 0 {
		t.Errorf("servers is not empty")
	}
	if config.DefaultServer != "" {
		t.Errorf("default server is not empty")
	}
}

func TestAddEntry(t *testing.T) {
	config, err := getConfig(t, testYaml)
	if err != nil {
		t.Error(err)
	}

	serverName := "test_name"
	serverHost := "another_port"
	serverPassword := "an0therR4nd0m:p@ssw0rd!"

	wantServerAmount := len(config.Servers)
	err = config.AddServerEntry(serverName, serverHost, serverPassword)
	if err != nil {
		t.Error(err)
	}

	gotServerAmount := len(config.Servers)

	if wantServerAmount == gotServerAmount {
		msg := fmt.Sprintf(
			"base server amount: %d got server amount: %d",
			wantServerAmount,
			gotServerAmount)

		t.Error(msg)
	}

	contentBytes, err := os.ReadFile(config.GetYamlPath())
	if err != nil {
		t.Error(err)
	}
	yamlStr := string(contentBytes)

	addServerData := []string{serverName, serverHost, serverPassword}
	for _, data := range addServerData {
		if !strings.Contains(yamlStr, data) {
			msg := fmt.Sprintf("yaml file is missing %s, failed to write", data)
			t.Error(msg)
		}
	}
}

func TestRemoveEntry(t *testing.T) {
	config, err := getConfig(t, testYaml)
	if err != nil {
		t.Error(err)
	}

	serverName := "default"

	err = config.RemoveServerEntry(serverName)
	if err != nil {
		t.Error(err)
	}

	if _, ok := config.Servers[serverName]; ok {
		t.Error("failed to remove file, existing entry found")
	}
}

func getConfig(t *testing.T, dataToWrite string) (*rcon.Config, error) {
	tempPath := t.TempDir()
	tempYamlPath := tempPath + "/mc-rcon.yml"

	err := os.WriteFile(tempYamlPath, []byte(dataToWrite), 0o644)
	if err != nil {
		return nil, err
	}

	config, err := rcon.NewConfig(tempPath)
	if err != nil {
		return nil, err
	}

	return config, err
}
