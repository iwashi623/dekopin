package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	CREATE_TAG_DEFAULT_REVISION = "LATEST"
)

var createTagCmd = &cobra.Command{
	Use:     "create-tag",
	Short:   "Assign a Revision tag to a Cloud Run revision",
	PreRunE: createTagPreRun,
	RunE:    CreateTagCommand,
}

func createTagPreRun(cmd *cobra.Command, args []string) error {
	dekopinCmd, err := GetDekopinCommand(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get dekopin command: %w", err)
	}

	tag, err := dekopinCmd.GetTagByFlag()
	if err != nil {
		return fmt.Errorf("failed to get tag flag: %w", err)
	}

	if tag != "" {
		if err := ValidateTag(tag); err != nil {
			return err
		}
	}

	return nil
}

func CreateTagCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gcloudCmd, err := GetGcloudCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gcloud command: %w", err)
	}

	dekopinCmd, err := GetDekopinCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get dekopin command: %w", err)
	}

	tf, err := dekopinCmd.GetTagByFlag()
	if err != nil {
		return fmt.Errorf("failed to get tag flag: %w", err)
	}

	tag, err := CreateRevisionTagName(ctx, tf)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	rf, err := dekopinCmd.GetRevisionByFlag()
	if err != nil {
		return fmt.Errorf("failed to get revision name: %w", err)
	}

	return createTag(ctx, gcloudCmd, tag, rf)
}

func createTag(ctx context.Context, gc GcloudCommand, tag string, revisionName string) error {
	if revisionName != CREATE_TAG_DEFAULT_REVISION {
		_, err := gc.GetRevision(ctx, revisionName)
		if err != nil {
			return fmt.Errorf("failed to get revision: %w", err)
		}
	}

	return gc.CreateRevisionTag(ctx, tag, revisionName)
}
