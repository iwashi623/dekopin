package dekopin

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Run(ctx context.Context) (int, error) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		return 1, err
	}
	return 0, nil
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
	RunE:  CreateRevision,
}

var createTagCmd = &cobra.Command{
	Use:   "create-tag",
	Short: "Assign a Revision tag to a Cloud Run revision",
	RunE:  CreateTag,
}

var removeTagCmd = &cobra.Command{
	Use:   "remove-tag",
	Short: "Remove a Revision tag from a Cloud Run revision",
	RunE:  RemoveTag,
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy new revision with tag and traffic management",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deploying new revision with tag management and traffic switching...")
		// TODO: Execute necessary subcommand processes (create-tag, remove-tag, switch-traffic) internally
	},
}

var stDeployCmd = &cobra.Command{
	Use:   "st-deploy",
	Short: "Deploy new revision with tag and traffic management",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deploying new revision with tag management and traffic switching...")
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
	rootCmd.AddCommand(stDeployCmd)
}

func setRootFlags(cmd *cobra.Command) {
	rootCmd.PersistentFlags().String("project", "", "GCP project id")
	rootCmd.PersistentFlags().String("region", "", "region")
	rootCmd.PersistentFlags().String("service", "", "service name")
	rootCmd.PersistentFlags().String("runner", "", "runner type")
	rootCmd.PersistentFlags().StringP("file", "f", "dekopin.yml", "config file name")
}

func getCommitHash(cmd *cobra.Command) string {
	if config.Runner == RUNNER_GITHUB_ACTIONS {
		sha := os.Getenv(ENV_GITHUB_SHA)
		// Get first 7 characters
		return sha[:7]
	}

	if config.Runner == RUNNER_CLOUD_BUILD {
		sha := os.Getenv(ENV_CLOUD_BUILD_SHA)
		// Get first 7 characters
		return sha[:7]
	}

	return ""
}

func getTagName(cmd *cobra.Command) (string, error) {
	tag, err := cmd.Flags().GetString("tag")
	if err != nil {
		return "", fmt.Errorf("failed to get tag flag: %w", err)
	}

	if tag == "" && config.Runner == RUNNER_LOCAL {
		return "", fmt.Errorf("tag flag is required")
	}

	if tag != "" {
		return tag, nil
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

	if config.Runner == RUNNER_LOCAL && rv == "" {
		return "", fmt.Errorf("revision flag is required")
	}

	if rv != "" {
		return rv, nil
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
