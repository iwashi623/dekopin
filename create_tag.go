package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var createTagCmd = &cobra.Command{
	Use:     "create-tag",
	Short:   "Assign a Revision tag to a Cloud Run revision",
	PreRunE: createTagPreRun,
	RunE:    CreateTagCommand,
}

func createTagPreRun(cmd *cobra.Command, args []string) error {
	tag, err := getTagByFlag(cmd)
	if err != nil {
		return err
	}

	if tag != "" {
		if err := validateTag(tag); err != nil {
			return err
		}
	}

	return nil
}

func CreateTagCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gcloudCmd, ok := ctx.Value(gcloudCmdKey{}).(GcloudCommand)
	if !ok {
		return fmt.Errorf("gcloud command not found")
	}

	tf, err := getTagByFlag(cmd)
	if err != nil {
		return fmt.Errorf("failed to get tag flag: %w", err)
	}

	tag, err := createRevisionTagName(ctx, tf)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	rf, err := getRevisionByFlag(cmd)
	if err != nil {
		return fmt.Errorf("failed to get revision name: %w", err)
	}

	return createTag(ctx, gcloudCmd, tag, rf)
}

func createTag(ctx context.Context, gc GcloudCommand, tag string, revisionName string) error {
	if revisionName != DEFAULT_REVISION {
		// revisionが存在するか確認する
		_, err := gc.GetRevision(ctx, revisionName)
		if err != nil {
			return fmt.Errorf("failed to get revision: %w", err)
		}
	}

	return gc.CreateRevisionTag(ctx, tag, revisionName)
}
