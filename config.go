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
		return fmt.Errorf("ファイル名の取得に失敗しました: %w", err)
	}

	dekopinYaml, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}

	err = yaml.Unmarshal(dekopinYaml, &config)
	if err != nil {
		return fmt.Errorf("ファイルのパースに失敗しました: %w", err)
	}

	err = config.validate()
	if err != nil {
		return fmt.Errorf("設定の検証に失敗しました: %w", err)
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

	ENV_GITHUB_TAG      = "GITHUB_REF_NAME"
	ENV_CLOUD_BUILD_TAG = "TAG_NAME"

	ENV_GITHUB_SHA      = "GITHUB_SHA"
	ENV_CLOUD_BUILD_SHA = "COMMIT_SHA"
)

var validRunners = map[string]bool{
	RUNNER_GITHUB_ACTIONS: true,
	RUNNER_CLOUD_BUILD:    true,
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
