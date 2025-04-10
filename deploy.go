package dekopin

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const (
	DEPLOY_DEFAULT_REVISION = "LATEST"
)

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Short:   "Deploy new revision with image",
	PreRunE: deployPreRun,
	RunE:    deployCommand,
}

func deployPreRun(cmd *cobra.Command, args []string) error {
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

type DeployCommandFlags struct {
	Image            string
	Tag              string
	ShouldCreateTag  bool
	ShouldRemoveTags bool
}

func deployCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gc, err := GetGCloud(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gcloud command: %w", err)
	}

	flags, err := getDeployCommandFlags(cmd)
	if err != nil {
		return fmt.Errorf("failed to get deploy command flags: %w", err)
	}

	if flags.Tag == "" && flags.ShouldCreateTag {
		flags.Tag, err = CreateRevisionTagName(ctx, flags.Tag)
		if err != nil {
			return fmt.Errorf("failed to get tag name: %w", err)
		}
	}

	commitHash, err := GetCommitHash(ctx)
	if err != nil {
		if !errors.Is(err, ErrGetCommitHashInLocal) {
			return err
		}
	}

	return deploy(ctx, gc, flags, commitHash)
}

func deploy(
	ctx context.Context,
	gc GCloud,
	flags *DeployCommandFlags,
	commitHash string,
) error {
	if err := gc.DeployWithTraffic(ctx, flags.Image, commitHash); err != nil {
		return fmt.Errorf("failed to deploy to Cloud Run: %w", err)
	}

	if flags.ShouldCreateTag {
		if err := gc.CreateRevisionTag(ctx, flags.Tag, DEPLOY_DEFAULT_REVISION); err != nil {
			return fmt.Errorf("failed to create revision tag: %w", err)
		}
	}

	if flags.ShouldRemoveTags {
		activeRevisionTags, err := gc.GetActiveRevisionTags(ctx)
		if err != nil {
			return fmt.Errorf("failed to get active revision tag: %w", err)
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

func getDeployCommandFlags(cmd *cobra.Command) (*DeployCommandFlags, error) {
	dekopinCmd, err := GetDekopinCommand(cmd.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get dekopin command: %w", err)
	}

	image, err := dekopinCmd.GetImageByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get image flag: %w", err)
	}

	tag, err := dekopinCmd.GetTagByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get tag flag: %w", err)
	}

	createTag, err := dekopinCmd.GetCreateTagByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get create-tag flag: %w", err)
	}

	removeTags, err := dekopinCmd.GetRemoveTagsByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get remove-tags flag: %w", err)
	}

	return &DeployCommandFlags{
		Image:            image,
		Tag:              tag,
		ShouldCreateTag:  createTag,
		ShouldRemoveTags: removeTags,
	}, nil
}
