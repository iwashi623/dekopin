package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var createRevisionCmd = &cobra.Command{
	Use:   "create-revision",
	Short: "Create a new Cloud Run revision",
	RunE:  createRevisionCommand,
}

func createRevisionCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gc, err := GetGCloud(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gcloud command: %w", err)
	}

	dekopinCmd, err := GetDekopinCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get dekopin command: %w", err)
	}

	image, err := dekopinCmd.GetImageByFlag()
	if err != nil {
		return fmt.Errorf("failed to get image flag: %w", err)
	}

	commitHash, err := GetCommitHash(ctx)
	if err != nil {
		return fmt.Errorf("failed to get commit hash: %w", err)
	}

	return createRevision(ctx, gc, image, commitHash)
}

func createRevision(ctx context.Context, gc GCloud, image string, commitHash string) error {
	if err := gc.CreateRevision(ctx, image, commitHash); err != nil {
		return fmt.Errorf("failed to create revision: %w", err)
	}

	return nil
}
