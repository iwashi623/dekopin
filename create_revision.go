package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func CreateRevisionCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gcloudCmd, ok := ctx.Value(gcloudCmdKey{}).(GcloudCommand)
	if !ok {
		return fmt.Errorf("gcloud command not found")
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
