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
	ctx := cmd.Context()
	dekopinCmd, err := GetDekopinCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get dekopin command: %w", err)
	}

	tag, err := dekopinCmd.GetTagByFlag()
	if err != nil {
		return err
	}

	if tag != "" {
		if err := ValidateTag(tag); err != nil {
			return err
		}
	}

	return nil
}

func RemoveTagCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gcloudCmd, err := GetGcloudCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gcloud command: %w", err)
	}

	dekopinCmd, err := GetDekopinCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get dekopin command: %w", err)
	}

	tag, err := dekopinCmd.GetTagByFlag()
	if err != nil {
		return fmt.Errorf("failed to get tag flag: %w", err)
	}

	tag, err = CreateRevisionTagName(ctx, tag)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	return removeTag(ctx, gcloudCmd, tag)
}

func removeTag(ctx context.Context, gc GcloudCommand, tag string) error {
	return gc.RemoveRevisionTag(ctx, tag)
}
