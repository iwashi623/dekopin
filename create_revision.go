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

	return CreateRevisionExec(ctx, gcloudCmd, image, commitHash)
}

func CreateRevisionExec(ctx context.Context, gc GcloudCommand, image string, commitHash string) error {
	fmt.Println("Starting to create a new Cloud Run revision...")

	if err := gc.CreateRevision(ctx, image, commitHash); err != nil {
		return fmt.Errorf("failed to create revision: %w", err)
	}

	fmt.Println("New revision has been successfully created")

	return nil
}
