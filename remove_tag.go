package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var removeTagCmd = &cobra.Command{
	Use:     "remove-tag",
	Short:   "Remove a Revision tag from a Cloud Run revision",
	PreRunE: removeTagPreRun,
	RunE:    RemoveTagCommand,
}

func removeTagPreRun(cmd *cobra.Command, args []string) error {
	tag, err := getTagByFlag(cmd)
	if err != nil {
		return err
	}

	if err := validateTag(tag); err != nil {
		return err
	}

	return nil
}

func RemoveTagCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gcloudCmd, ok := ctx.Value(gcloudCmdKey{}).(GcloudCommand)
	if !ok {
		return fmt.Errorf("gcloud command not found")
	}

	tag, err := getTagByFlag(cmd)
	if err != nil {
		return fmt.Errorf("failed to get tag flag: %w", err)
	}

	tag, err = createRevisionTagName(ctx, tag)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	return removeTag(ctx, gcloudCmd, tag)
}

func removeTag(ctx context.Context, gc GcloudCommand, tag string) error {
	return gc.RemoveRevisionTag(ctx, tag)
}
