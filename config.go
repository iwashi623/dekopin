package dekopin

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func getConfig(cmd *cobra.Command) (*DekopinConfig, error) {
	dekopinCmd, err := GetDekopinCommand(cmd.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get dekopin command: %w", err)
	}

	fileName, err := dekopinCmd.GetFileByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get filename: %w", err)
	}

	dekopinYaml, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config *DekopinConfig
	if err := yaml.Unmarshal(dekopinYaml, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	return config, nil
}

type DekopinConfig struct {
	Project string `yaml:"project"`
	Region  string `yaml:"region"`
	Service string `yaml:"service"`
	Runner  string `yaml:"runner"`
}

const (
	RUNNER_GITHUB_ACTIONS = "github-actions"
	RUNNER_CLOUD_BUILD    = "cloud-build"
	RUNNER_LOCAL          = "local"

	ENV_GITHUB_REF      = "GITHUB_REF"
	ENV_CLOUD_BUILD_REF = "REF_NAME"

	ENV_GITHUB_SHA      = "GITHUB_SHA"
	ENV_CLOUD_BUILD_SHA = "COMMIT_SHA"
)

var ValidRunners = []string{
	RUNNER_GITHUB_ACTIONS,
	RUNNER_CLOUD_BUILD,
	RUNNER_LOCAL,
}
