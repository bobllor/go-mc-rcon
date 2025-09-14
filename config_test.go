package main

import (
	"fmt"
	"os"
	"rcon/rcon"
	"testing"
)

var testYaml string = `
servers:
  default:
    host: "some_ip:at_some_port"
    rcon_password: "go_test_password"
default_server: "default"
`

func getConfig(tempPath string) (*rcon.Config, error) {
	tempYamlPath := tempPath + "/mc-rcon.yml"

	err := os.WriteFile(tempYamlPath, []byte(testYaml), 0o644)
	if err != nil {
		return nil, err
	}

	config, err := rcon.NewConfig(tempPath)
	if err != nil {
		return nil, err
	}

	return config, err
}

func TestCreateConfig(t *testing.T) {
	dir := t.TempDir()

	_, err := getConfig(dir)
	if err != nil {
		t.Error(err)
	}
}

func TestAddEntry(t *testing.T) {
	dir := t.TempDir()
	config, err := getConfig(dir)
	if err != nil {
		t.Error(err)
	}

	wantServerAmount := len(config.Servers)
	err = config.AddServerEntry("test_name", "another_port", "an0therR4nd0m:p@ssw0rd!")
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
}

func TestRemoveEntry(t *testing.T) {
	dir := t.TempDir()

	config, err := getConfig(dir)
	if err != nil {
		t.Error(err)
	}

	wantServerAmount := len(config.Servers)

	err = config.RemoveServerEntry("default")
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
}
