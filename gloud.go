package dekopin

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

type GcloudCommand interface {
	CreateTag(ctx context.Context, revisionTag string, revisionName string) error
	RemoveTag(ctx context.Context, revisionTag string) error
	CreateRevision(ctx context.Context, imageName string, commitHash string) error
}

type gcloudCommand struct {
	Stdout io.Writer
	Stderr io.Writer
}

func NewGcloudCommand(stdout io.Writer, stderr io.Writer) GcloudCommand {
	return &gcloudCommand{
		Stdout: stdout,
		Stderr: stderr,
	}
}

func (c *gcloudCommand) CreateTag(ctx context.Context, revisionTag string, revisionName string) error {
	cmd := exec.CommandContext(ctx, "gcloud", "run", "services", "update-traffic", config.Service,
		"--region", config.Region,
		"--project", config.Project,
		"--update-tags", revisionTag+"="+revisionName,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (c *gcloudCommand) RemoveTag(ctx context.Context, revisionTag string) error {
	cmd := exec.CommandContext(ctx, "gcloud", "run", "services", "update-traffic", config.Service,
		"--region", config.Region,
		"--project", config.Project,
		"--remove-tags", revisionTag,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	return nil
}

func (c *gcloudCommand) CreateRevision(ctx context.Context, imageName string, commitHash string) error {
	cmd := exec.CommandContext(ctx, "gcloud", "run", "deploy", config.Service,
		"--image", imageName,
		"--project", config.Project,
		"--region", config.Region,
		"--no-traffic", // Important: Do not route traffic to the new revision
	)

	if commitHash != "" {
		cmd.Args = append(cmd.Args, "--revision-suffix", commitHash)
	}

	// Execute command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy to Cloud Run: %w", err)
	}

	return nil
}
