package dekopin

import (
	"context"
	"fmt"
	"slices"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var stDeployCmd = &cobra.Command{
	Use:     "st-deploy",
	Short:   "Switch Tag Deploy(Assign a Revision tag to a Cloud Run revision)",
	PreRunE: stDeployPreRun,
	RunE:    switchTagDeployCommand,
}

func stDeployPreRun(cmd *cobra.Command, args []string) error {
	dekopinCmd, err := GetDekopinCommand(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get dekopin command: %w", err)
	}

	tag, err := dekopinCmd.GetTagByFlag()
	if err != nil {
		return fmt.Errorf("failed to get tag flag: %w", err)
	}

	if err := ValidateTag(tag); err != nil {
		return err
	}

	return nil
}

func switchTagDeployCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	dekopinCmd, err := GetDekopinCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get dekopin command: %w", err)
	}

	tag, err := dekopinCmd.GetTagByFlag()
	if err != nil {
		return fmt.Errorf("failed to get tag flag: %w", err)
	}

	rt, err := CreateRevisionTagName(ctx, tag)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	gc, err := GetGCloud(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get gcloud command: %w", err)
	}

	removeTags, err := dekopinCmd.GetRemoveTagsByFlag()
	if err != nil {
		return fmt.Errorf("failed to get remove-tags flag: %w", err)
	}

	return switchTagDeploy(cmd.Context(), gc, rt, removeTags)
}

func switchTagDeploy(ctx context.Context, gc GCloud, tag string, removeTags bool) error {
	tags, err := gc.GetActiveRevisionTags(ctx)
	if err != nil {
		return fmt.Errorf("failed to get service: %w", err)
	}

	if !slices.Contains(tags, tag) {
		return fmt.Errorf("active tag %s not found", tag)
	}

	if err := gc.UpdateTrafficToRevisionTag(ctx, tag); err != nil {
		return fmt.Errorf("failed to update traffic to revision tag: %w", err)
	}

	if removeTags {
		filteredTags := lo.Filter(tags, func(t string, _ int) bool {
			return t != tag
		})

		if err := gc.RemoveRevisionTags(ctx, filteredTags); err != nil {
			return fmt.Errorf("failed to remove revision tags: %w", err)
		}
	}

	return nil
}
