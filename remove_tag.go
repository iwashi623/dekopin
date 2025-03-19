package dekopin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func RemoveTagCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	tag, err := getTagName(cmd)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	return removeTag(ctx, tag)
}

func removeTag(ctx context.Context, tag string) error {
	formattedTag := "tag-" + strings.ReplaceAll(tag, ".", "-")

	gcloudCmd := exec.CommandContext(ctx, "gcloud", "run", "services", "update-traffic", config.Service,
		"--region", config.Region,
		"--project", config.Project,
		"--remove-tags", formattedTag,
	)

	gcloudCmd.Stdout = os.Stdout
	gcloudCmd.Stderr = os.Stderr

	if err := gcloudCmd.Run(); err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	return nil
}
