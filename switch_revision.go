package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	SWITCH_REVISION_DEFAULT_REVISION = "LATEST"
)

var srDeployCmd = &cobra.Command{
	Use:   "sr-deploy",
	Short: "Switch Revision Deploy(Deploy new revision with revision name)",
	RunE:  switchRevisionDeployCommand,
}

func switchRevisionDeployCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gc, err := GetGCloud(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gcloud command: %w", err)
	}

	dekopinCmd, err := GetDekopinCommand(ctx)
	if err != nil {
		return fmt.Errorf("failed to get dekopin command: %w", err)
	}

	revision, err := dekopinCmd.GetRevisionByFlag()
	if err != nil {
		return fmt.Errorf("failed to get revision flag: %w", err)
	}

	return switchRevisionDeploy(ctx, gc, revision)
}

func switchRevisionDeploy(ctx context.Context, gc GCloud, revision string) error {
	if revision != SWITCH_REVISION_DEFAULT_REVISION {
		_, err := gc.GetRevision(ctx, revision)
		if err != nil {
			return fmt.Errorf("failed to get revision: %w", err)
		}
	}

	if err := gc.UpdateTrafficToRevision(ctx, revision); err != nil {
		return fmt.Errorf("failed to update traffic to latest revision: %w", err)
	}

	return nil
}
