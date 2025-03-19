package dekopin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func CreateTag(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	tag, err := getTagName(cmd)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	revision, err := getRevisionName(cmd)
	if err != nil {
		return fmt.Errorf("failed to get revision name: %w", err)
	}

	return createTag(ctx, tag, revision)
}

func createTag(ctx context.Context, tag string, revision string) error {
	// Convert tag format
	formattedTag := "tag-" + strings.ReplaceAll(tag, ".", "-")

	gcloudCmd := exec.CommandContext(ctx, "gcloud", "run", "services", "update-traffic", config.Service,
		"--region", config.Region,
		"--project", config.Project,
		"--update-tags", formattedTag+"="+revision,
	)

	gcloudCmd.Stdout = os.Stdout
	gcloudCmd.Stderr = os.Stderr

	// Execute command
	if err := gcloudCmd.Run(); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}
