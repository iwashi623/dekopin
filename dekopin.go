package dekopin

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	run "cloud.google.com/go/run/apiv2"
	"github.com/spf13/cobra"
)

const (
	TIMEOUT = 120 * time.Second
)

func Run(ctx context.Context) int {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()

	sc, err := run.NewServicesClient(ctx)
	if err != nil {
		log.Printf("ERROR: failed to create services client: %s", err)
		return 1
	}
	defer sc.Close()

	rc, err := run.NewRevisionsClient(ctx)
	if err != nil {
		log.Printf("ERROR: failed to create revisions client: %s", err)
		return 1
	}
	defer rc.Close()

	ctx = SetGCloud(ctx, NewGCloud(os.Stdout, os.Stderr, sc, rc))

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Printf("ERROR: %s", err)
		return 1
	}
	return 0
}

var rootCmd = &cobra.Command{
	Use:               "dekopin",
	Short:             "Dekopin is a Cloud Run deployment tool",
	Long:              "Dekopin is a tool to deploy Cloud Run services with tags and traffic management.",
	PersistentPreRunE: prepareAllRun,
}

func init() {
	setRootFlags(rootCmd)

	rootCmd.AddCommand(createRevisionCmd)
	createRevisionCmd.Flags().StringP("image", "i", "", "container image")
	createRevisionCmd.MarkFlagRequired("image")

	rootCmd.AddCommand(createTagCmd)
	createTagCmd.Flags().StringP("tag", "t", "", "tag name")
	createTagCmd.Flags().String("revision", CREATE_TAG_DEFAULT_REVISION, "revision name")
	createTagCmd.Flags().Bool("update-traffic", false, "update traffic to the revision after deploy")
	createTagCmd.Flags().Bool("remove-tags", false, "remove all revision tags before deploy")

	rootCmd.AddCommand(removeTagCmd)
	removeTagCmd.Flags().StringP("tag", "t", "", "tag name")
	removeTagCmd.MarkFlagRequired("tag")

	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringP("image", "i", "", "container image")
	deployCmd.MarkFlagRequired("image")
	deployCmd.Flags().StringP("tag", "t", "", "new revision tag name")
	deployCmd.Flags().Bool("create-tag", false, "create a revision tag after deploy")
	deployCmd.Flags().Bool("remove-tags", false, "remove all revision tags before deploy")

	rootCmd.AddCommand(srDeployCmd)
	srDeployCmd.Flags().String("revision", SWITCH_REVISION_DEFAULT_REVISION, "revision name")

	rootCmd.AddCommand(stDeployCmd)
	stDeployCmd.Flags().StringP("tag", "t", "", "tag name")
	stDeployCmd.MarkFlagRequired("tag")
	stDeployCmd.Flags().Bool("remove-tags", false, "remove all revision tags except the deployment target revision tag")
}

func setRootFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().String("project", "", "GCP project id")
	rootCmd.PersistentFlags().String("region", "", "region")
	rootCmd.PersistentFlags().String("service", "", "service name")
	rootCmd.PersistentFlags().String("runner", "", "runner type")
	rootCmd.PersistentFlags().StringP("file", "f", "dekopin.yml", "config file name")
}

func prepareAllRun(cmd *cobra.Command, args []string) error {
	dekopinCmd := NewDekopinCommand(cmd)
	ctx := SetDekopinCommand(cmd.Context(), dekopinCmd)
	cmd.SetContext(ctx)

	var config *DekopinConfig
	fileName, err := dekopinCmd.GetFileByFlag()
	if err != nil {
		return err
	}

	if fileName != "" {
		config, err = getConfig(fileName)
		if err != nil {
			return err
		}
	}

	cmdOption, err := NewCmdOption(ctx, config, cmd)
	if err != nil {
		return err
	}

	ctx = SetCmdOption(ctx, cmdOption)
	cmd.SetContext(ctx)
	return nil
}

const (
	COMMIT_HASH_LENGTH = 7
)

func GetCommitHash(ctx context.Context) (string, error) {
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get cmdOption: %w", err)
	}

	if opt.Runner == RUNNER_GITHUB_ACTIONS {
		sha := os.Getenv(ENV_GITHUB_SHA)
		if len(sha) == 0 {
			return "", fmt.Errorf("ref name is required")
		}
		if len(sha) <= COMMIT_HASH_LENGTH {
			return sha, nil
		}
		return sha[:COMMIT_HASH_LENGTH], nil
	}

	if opt.Runner == RUNNER_CLOUD_BUILD {
		sha := os.Getenv(ENV_CLOUD_BUILD_SHA)
		if len(sha) == 0 {
			return "", fmt.Errorf("ref name is required")
		}
		if len(sha) <= COMMIT_HASH_LENGTH {
			return sha, nil
		}
		return sha[:COMMIT_HASH_LENGTH], nil
	}

	return "", ErrGetCommitHashInLocal
}

func GetRunnerRef(ctx context.Context) (string, error) {
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get cmdOption: %w", err)
	}

	if opt.Runner == RUNNER_GITHUB_ACTIONS {
		return os.Getenv(ENV_GITHUB_REF), nil
	}

	if opt.Runner == RUNNER_CLOUD_BUILD {
		return os.Getenv(ENV_CLOUD_BUILD_REF), nil
	}

	return "", fmt.Errorf("ref name is required")
}

func CreateRevisionTagName(ctx context.Context, tag string) (string, error) {
	if tag != "" {
		return tag, nil
	}

	opt, err := GetCmdOption(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get cmdOption: %w", err)
	}

	if tag == "" && opt.Runner == RUNNER_LOCAL {
		return "", fmt.Errorf("local execution requires the tag flag")
	}

	ref, err := GetRunnerRef(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get runner ref: %w", err)
	}

	reg := regexp.MustCompile(`[./: _]`)

	return "tag-" + reg.ReplaceAllString(ref, "-"), nil
}

func ValidateTag(tag string) error {
	reg := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !reg.MatchString(tag) && tag != "" {
		return fmt.Errorf("invalid tag name. Valid values: lowercase alphanumeric, numbers, hyphen")
	}
	return nil
}
