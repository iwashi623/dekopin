package dekopin

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

type GcloudCommand interface {
	CreateRevision(ctx context.Context, imageName string, revisionName string) error
	CreateRevisionTag(ctx context.Context, revisionTag string, revisionName string) error
	RemoveRevisionTag(ctx context.Context, revisionTag string) error
	Deploy(ctx context.Context, imageName string, revisionName string, useTraffic bool) error
	UpdateTrafficToLatestRevision(ctx context.Context) error
	DeployWithTraffic(ctx context.Context, imageName string, revisionName string) error
}

type gcloudCommand struct {
	Stdout io.Writer
	Stderr io.Writer
}

var _ GcloudCommand = &gcloudCommand{}

func NewGcloudCommand(stdout io.Writer, stderr io.Writer) GcloudCommand {
	return &gcloudCommand{
		Stdout: stdout,
		Stderr: stderr,
	}
}

func (c *gcloudCommand) CreateRevision(ctx context.Context, imageName string, revisionName string) error {
	if err := c.Deploy(ctx, imageName, revisionName, false); err != nil {
		return fmt.Errorf("failed to create revision: %w", err)
	}

	return nil
}

func (c *gcloudCommand) CreateRevisionTag(ctx context.Context, revisionTag string, revisionName string) error {
	cmd := updateTrafficCmd(ctx)
	cmd.Args = append(cmd.Args, "--update-tags", revisionTag+"="+revisionName)

	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (c *gcloudCommand) RemoveRevisionTag(ctx context.Context, revisionTag string) error {
	cmd := updateTrafficCmd(ctx)
	cmd.Args = append(cmd.Args, "--remove-tags", revisionTag)

	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	return nil
}

func (c *gcloudCommand) Deploy(ctx context.Context, imageName string, revisionName string, useTraffic bool) error {
	cmd := runDeployCmd(ctx)
	cmd.Args = append(cmd.Args, "--image", imageName)

	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	if revisionName != "" {
		cmd.Args = append(cmd.Args, "--revision-suffix", revisionName)
	}

	if !useTraffic {
		fmt.Println("Deploying without traffic")
		cmd.Args = append(cmd.Args, "--no-traffic")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy to Cloud Run: %w", err)
	}

	return nil
}

func (c *gcloudCommand) UpdateTrafficToLatestRevision(ctx context.Context) error {
	cmd := updateTrafficCmd(ctx)
	cmd.Args = append(cmd.Args, "--to-latest")

	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update traffic to latest revision: %w", err)
	}

	return nil
}

func (c *gcloudCommand) DeployWithTraffic(ctx context.Context, imageName string, revisionName string) error {
	if err := c.Deploy(ctx, imageName, revisionName, true); err != nil {
		return fmt.Errorf("failed to deploy to Cloud Run: %w", err)
	}

	if err := c.UpdateTrafficToLatestRevision(ctx); err != nil {
		return fmt.Errorf("failed to update traffic to latest revision: %w", err)
	}

	return nil
}

func updateTrafficCmd(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, "gcloud", "run", "services", "update-traffic", config.Service,
		"--region", config.Region,
		"--project", config.Project,
	)
}

func runDeployCmd(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, "gcloud", "run", "deploy", config.Service,
		"--project", config.Project,
		"--region", config.Region,
	)
}
