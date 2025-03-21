package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var srDeployCmd = &cobra.Command{
	Use:   "sr-deploy",
	Short: "Switch Revision Deploy(Deploy new revision with revision name)",
	RunE:  SwitchRevisionDeployCommand,
}

func SwitchRevisionDeployCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gc, err := GetGcloudCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gcloud command: %w", err)
	}

	revision, err := getRevisionByFlag(cmd)
	if err != nil {
		return fmt.Errorf("failed to get revision flag: %w", err)
	}

	return switchRevisionDeploy(ctx, gc, revision)
}

func switchRevisionDeploy(ctx context.Context, gc GcloudCommand, revision string) error {
	_, err := gc.GetRevision(ctx, revision)
	if err != nil {
		return fmt.Errorf("failed to get revision: %w", err)
	}

	if err := gc.UpdateTrafficToRevision(ctx, revision); err != nil {
		return fmt.Errorf("failed to update traffic to latest revision: %w", err)
	}

	return nil
}
