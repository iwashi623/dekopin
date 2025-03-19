package dekopin

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func CreateTag(cmd *cobra.Command, args []string) error {
	tag, err := getTagName(cmd)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	revision, err := getRevisionName(cmd)
	if err != nil {
		return fmt.Errorf("failed to get revision name: %w", err)
	}

	return createTag(cmd, tag, revision)
}

func createTag(cmd *cobra.Command, tag string, revision string) error {
	ctx := cmd.Context()
	// Convert tag format
	formattedTag := strings.ReplaceAll(tag, ".", "-")

	gcloudCmd := exec.CommandContext(ctx, "gcloud", "run", "services", "update-traffic", config.Service,
		"--region", config.Region,
		"--project", config.Project,
		"--update-tags", "tag-"+formattedTag+"="+revision,
		"--to-revisions", revision+"=0",
	)

	gcloudCmd.Stdout = os.Stdout
	gcloudCmd.Stderr = os.Stderr

	// Execute command
	if err := gcloudCmd.Run(); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}
