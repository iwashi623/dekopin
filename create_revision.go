package dekopin

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func CreateRevision(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	image, err := cmd.Flags().GetString("image")
	if err != nil {
		return fmt.Errorf("failed to get image flag: %w", err)
	}

	if image == "" {
		return fmt.Errorf("image flag is required")
	}

	commitHash := getCommitHash(cmd)

	return createRevision(ctx, image, commitHash)
}

func createRevision(ctx context.Context, image string, commitHash string) error {
	fmt.Println("Starting to create a new Cloud Run revision...")

	gcloudCmd := exec.CommandContext(ctx, "gcloud", "run", "deploy", config.Service,
		"--image", image,
		"--project", config.Project,
		"--region", config.Region,
		"--no-traffic", // Important: Do not route traffic to the new revision
	)

	if commitHash != "" {
		gcloudCmd.Args = append(gcloudCmd.Args, "--revision-suffix", commitHash)
	}

	gcloudCmd.Stdout = os.Stdout
	gcloudCmd.Stderr = os.Stderr

	// Execute command
	if err := gcloudCmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy to Cloud Run: %w", err)
	}

	fmt.Println("New revision has been successfully deployed")

	return nil
}
