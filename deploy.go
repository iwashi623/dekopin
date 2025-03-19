package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func DeployCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gcloudCmd, ok := ctx.Value(gcloudCmdKey{}).(GcloudCommand)
	if !ok {
		return fmt.Errorf("gcloud 	command not found")
	}

	image, err := cmd.Flags().GetString("image")
	if err != nil {
		return fmt.Errorf("failed to get image flag: %w", err)
	}

	commitHash := getCommitHash()

	return deploy(ctx, gcloudCmd, image, commitHash)
}

func deploy(ctx context.Context, gcloudCmd GcloudCommand, imageName string, commitHash string) error {
	if err := gcloudCmd.DeployWithTraffic(ctx, imageName, commitHash); err != nil {
		return fmt.Errorf("failed to deploy to Cloud Run: %w", err)
	}

	return nil
}
