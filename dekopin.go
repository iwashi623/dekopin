package dekopin

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type gcloudCmdKey struct{}

func Run(ctx context.Context) int {
	gcloudCmd := NewGcloudCommand(os.Stdout, os.Stderr)
	ctx = context.WithValue(ctx, gcloudCmdKey{}, gcloudCmd)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Printf("ERROR: %s", err)
		return 1
	}
	return 0
}

var rootCmd = &cobra.Command{
	Use:   "dekopin",
	Short: "Dekopin is a Cloud Run deployment tool",
	Long:  "Dekopin is a tool to deploy Cloud Run services with tags and traffic management.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateRunnerFunc(cmd, args); err != nil {
			return err
		}
		return setConfig(cmd)
	},
}

var createRevisionCmd = &cobra.Command{
	Use:   "create-revision",
	Short: "Create a new Cloud Run revision",
	RunE:  CreateRevisionCommand,
}

var createTagCmd = &cobra.Command{
	Use:   "create-tag",
	Short: "Assign a Revision tag to a Cloud Run revision",
	RunE:  CreateTagCommand,
}

var removeTagCmd = &cobra.Command{
	Use:   "remove-tag",
	Short: "Remove a Revision tag from a Cloud Run revision",
	RunE:  RemoveTagCommand,
}

var deployCmd = &cobra.Command{
	Use: "deploy",
	// イメージを指定して新しいリビジョンを作成
	Short: "Deploy new revision with image",
	RunE:  DeployCommand,
}

var srDeployCmd = &cobra.Command{
	Use:   "sr-deploy",
	Short: "Switch Revision Deploy(Deploy new revision with revision name)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Switching revision")
	},
}

var stDeployCmd = &cobra.Command{
	Use:   "st-deploy",
	Short: "Switch Tag Deploy(Assign a Revision tag to a Cloud Run revision)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Switching tag")
	},
}

func init() {
	setRootFlags(rootCmd)

	rootCmd.AddCommand(createRevisionCmd)
	createRevisionCmd.Flags().StringP("image", "i", "", "container image")
	createRevisionCmd.MarkFlagRequired("image")

	rootCmd.AddCommand(createTagCmd)
	createTagCmd.Flags().StringP("tag", "t", "", "tag name")
	createTagCmd.Flags().String("revision", "", "revision name")

	rootCmd.AddCommand(removeTagCmd)
	removeTagCmd.Flags().StringP("tag", "t", "", "tag name")

	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(srDeployCmd)
	rootCmd.AddCommand(stDeployCmd)
}

func setRootFlags(cmd *cobra.Command) {
	rootCmd.PersistentFlags().String("project", "", "GCP project id")
	rootCmd.PersistentFlags().String("region", "", "region")
	rootCmd.PersistentFlags().String("service", "", "service name")
	rootCmd.PersistentFlags().String("runner", "", "runner type")
	rootCmd.PersistentFlags().StringP("file", "f", "dekopin.yml", "config file name")
}

const (
	COMMIT_HASH_LENGTH = 7
)

func getCommitHash(cmd *cobra.Command) string {
	if config.Runner == RUNNER_GITHUB_ACTIONS {
		sha := os.Getenv(ENV_GITHUB_SHA)
		if len(sha) < COMMIT_HASH_LENGTH {
			return ""
		}
		return sha[:COMMIT_HASH_LENGTH]
	}

	if config.Runner == RUNNER_CLOUD_BUILD {
		sha := os.Getenv(ENV_CLOUD_BUILD_SHA)
		if len(sha) < COMMIT_HASH_LENGTH {
			return ""
		}
		return sha[:COMMIT_HASH_LENGTH]
	}

	return ""
}

func getTagName(cmd *cobra.Command) (string, error) {
	tag, err := cmd.Flags().GetString("tag")
	if err != nil {
		return "", fmt.Errorf("failed to get tag flag: %w", err)
	}

	if tag != "" {
		return tag, nil
	}

	if tag == "" && config.Runner == RUNNER_LOCAL {
		return "", fmt.Errorf("tag flag is required")
	}

	if config.Runner == RUNNER_GITHUB_ACTIONS {
		return os.Getenv(ENV_GITHUB_REF), nil
	}

	if config.Runner == RUNNER_CLOUD_BUILD {
		return os.Getenv(ENV_CLOUD_BUILD_REF), nil
	}

	return "", fmt.Errorf("tag flag is required")
}

func getRevisionName(cmd *cobra.Command) (string, error) {
	rv, err := cmd.Flags().GetString("revision")
	if err != nil {
		return "", fmt.Errorf("failed to get revision flag: %w", err)
	}

	if rv != "" {
		return rv, nil
	}

	if config.Runner == RUNNER_LOCAL && rv == "" {
		return "", fmt.Errorf("revision flag is required")
	}

	if prefix := getCommitHash(cmd); prefix != "" {
		return config.Service + "-" + prefix, nil
	}

	return "", fmt.Errorf("revision flag is required")
}

func validateRunner(runner string) error {
	if !validRunners[runner] {
		return fmt.Errorf("invalid runner type. Valid values: github-actions, cloud-build, local")
	}
	return nil
}
