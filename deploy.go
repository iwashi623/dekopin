package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	DEFAULT_REVISION = "LATEST"
)

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Short:   "Deploy new revision with image",
	PreRunE: deployPreRun,
	RunE:    DeployCommand,
}

func deployPreRun(cmd *cobra.Command, args []string) error {
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

type DeployCommandFlags struct {
	Image      string
	Tag        string
	CreateTag  bool
	RemoveTags bool
}

func DeployCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gcloudCmd, ok := ctx.Value(gcloudCmdKey{}).(GcloudCommand)
	if !ok {
		return fmt.Errorf("gcloud command not found")
	}

	flags, err := getDeployCommandFlags(cmd)
	if err != nil {
		return fmt.Errorf("failed to get deploy command flags: %w", err)
	}

	tag, err := createRevisionTagName(ctx, flags.Tag)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}
	commitHash := getCommitHash(opt.Runner)

	tag, err = createRevisionTagName(ctx, tag)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	return deploy(ctx, gcloudCmd, flags, commitHash, tag)
}

func deploy(
	ctx context.Context,
	gcloudCmd GcloudCommand,
	flags *DeployCommandFlags,
	commitHash string,
	newRevisionTag string,
) error {
	if flags.RemoveTags {
		activeRevisionTags, err := gcloudCmd.GetActiveRevisionTags(ctx)
		if err != nil {
			return fmt.Errorf("failed to get active revision tag: %w", err)
		}

		for _, activeRevisionTag := range activeRevisionTags {
			if err := gcloudCmd.RemoveRevisionTag(ctx, activeRevisionTag); err != nil {
				return fmt.Errorf("failed to remove revision tag: %w", err)
			}
		}
	}

	if err := gcloudCmd.DeployWithTraffic(ctx, flags.Image, commitHash); err != nil {
		return fmt.Errorf("failed to deploy to Cloud Run: %w", err)
	}

	if flags.CreateTag {
		if err := gcloudCmd.CreateRevisionTag(ctx, newRevisionTag, DEFAULT_REVISION); err != nil {
			return fmt.Errorf("failed to create revision tag: %w", err)
		}
	}

	return nil
}

func getDeployCommandFlags(cmd *cobra.Command) (*DeployCommandFlags, error) {
	image, err := cmd.Flags().GetString("image")
	if err != nil {
		return nil, fmt.Errorf("failed to get image flag: %w", err)
	}

	tag, err := getTagByFlag(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag flag: %w", err)
	}

	createTag, err := getCreateTagByFlag(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get create-tag flag: %w", err)
	}

	removeTag, err := getRemoveTagByFlag(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get remove-tags flag: %w", err)
	}

	return &DeployCommandFlags{
		Image:      image,
		Tag:        tag,
		CreateTag:  createTag,
		RemoveTags: removeTag,
	}, nil
}
