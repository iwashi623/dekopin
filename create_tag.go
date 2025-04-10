package dekopin

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const (
	CREATE_TAG_DEFAULT_REVISION = "LATEST"
)

var createTagCmd = &cobra.Command{
	Use:     "create-tag",
	Short:   "Assign a Revision tag to a Cloud Run revision",
	PreRunE: createTagPreRun,
	RunE:    createTagCommand,
}

type createTagCommandFlags struct {
	Tag                 string
	Revision            string
	ShouldRemoveTags    bool
	ShouldUpdateTraffic bool
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

func newCreateTagCommandFlags(ctx context.Context, cmd DekopinCommand) (*createTagCommandFlags, error) {
	tagFlag, err := cmd.GetTagByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get tag flag: %w", err)
	}

	tagName, err := CreateRevisionTagName(ctx, tagFlag)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag name: %w", err)
	}

	revisionName, err := cmd.GetRevisionByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get revision name: %w", err)
	}

	removeTags, err := cmd.GetRemoveTagsByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get remove tags flag: %w", err)
	}

	updateTraffic, err := cmd.GetUpdateTrafficByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get update traffic flag: %w", err)
	}

	return &createTagCommandFlags{
		Tag:                 tagName,
		Revision:            revisionName,
		ShouldRemoveTags:    removeTags,
		ShouldUpdateTraffic: updateTraffic,
	}, nil
}

func createTagCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gc, err := GetGCloud(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gcloud command: %w", err)
	}

	dekopinCmd, err := GetDekopinCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get dekopin command: %w", err)
	}

	flags, err := newCreateTagCommandFlags(ctx, dekopinCmd)
	if err != nil {
		return fmt.Errorf("failed to get create tag command flags: %w", err)
	}

	return createTag(ctx, gc, flags)
}

func createTag(ctx context.Context, gc GCloud, flags *createTagCommandFlags) error {
	if flags.Revision != CREATE_TAG_DEFAULT_REVISION {
		_, err := gc.GetRevision(ctx, flags.Revision)
		if err != nil {
			return fmt.Errorf("failed to get revision: %w", err)
		}
	}

	if err := gc.CreateRevisionTag(ctx, flags.Tag, flags.Revision); err != nil {
		return fmt.Errorf("failed to create revision tag: %w", err)
	}

	if flags.ShouldUpdateTraffic {
		if err := gc.UpdateTrafficToRevisionTag(ctx, flags.Tag); err != nil {
			return fmt.Errorf("failed to update traffic to revision tag: %w", err)
		}
	}

	if flags.ShouldRemoveTags {
		activeRevisionTags, err := gc.GetActiveRevisionTags(ctx)
		if err != nil {
			return fmt.Errorf("failed to get active revision tags: %w", err)
		}

		filteredTags := lo.Filter(activeRevisionTags, func(tag string, _ int) bool {
			return tag != flags.Tag
		})

		if err := gc.RemoveRevisionTags(ctx, filteredTags); err != nil {
			return fmt.Errorf("failed to remove revision tags: %w", err)
		}
	}

	return nil
}
