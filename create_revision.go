package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var createRevisionCmd = &cobra.Command{
	Use:   "create-revision",
	Short: "Create a new Cloud Run revision",
	RunE:  CreateRevisionCommand,
}

func CreateRevisionCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gcloudCmd, err := GetGcloudCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gcloud command: %w", err)
	}

	image, err := getImageByFlag(cmd)
	if err != nil {
		return fmt.Errorf("failed to get image flag: %w", err)
	}

	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}

	commitHash := getCommitHash(opt.Runner)

	return createRevision(ctx, gcloudCmd, image, commitHash)
}

func createRevision(ctx context.Context, gc GcloudCommand, image string, commitHash string) error {
	if err := gc.CreateRevision(ctx, image, commitHash); err != nil {
		return fmt.Errorf("failed to create revision: %w", err)
	}

	return nil
}
