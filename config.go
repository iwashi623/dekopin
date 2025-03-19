package dekopin

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var config DekopinConfig

func setConfig(cmd *cobra.Command) error {
	fileName, err := cmd.Flags().GetString("file")
	if err != nil {
		return fmt.Errorf("failed to get filename: %w", err)
	}

	dekopinYaml, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := yaml.Unmarshal(dekopinYaml, &config); err != nil {
		return fmt.Errorf("failed to parse configuration file: %w", err)
	}

	if err := config.validate(); err != nil {
		return fmt.Errorf("failed to validate configuration: %w", err)
	}

	return nil
}

func (c *DekopinConfig) validate() error {
	if err := validateRunner(c.Runner); err != nil {
		return err
	}

	return nil
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

var validRunners = map[string]bool{
	RUNNER_GITHUB_ACTIONS: true,
	RUNNER_CLOUD_BUILD:    true,
	RUNNER_LOCAL:          true,
}

func validateRunnerFunc(cmd *cobra.Command, args []string) error {
	runner, _ := cmd.Flags().GetString("runner")
	if runner == "" {
		return nil
	}

	if err := validateRunner(runner); err != nil {
		return err
	}
	return nil
}
