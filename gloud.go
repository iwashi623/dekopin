package dekopin

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type GcloudCommand interface {
	CreateRevision(ctx context.Context, imageName string, revisionName string) error
	CreateRevisionTag(ctx context.Context, revisionTag string, revisionName string) error
	RemoveRevisionTag(ctx context.Context, revisionTag string) error
	Deploy(ctx context.Context, imageName string, revisionName string, useTraffic bool) error
	UpdateTrafficToLatestRevision(ctx context.Context) error
	UpdateTrafficToRevision(ctx context.Context, revisionName string) error
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
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}

	cmd := updateTrafficCmd(ctx, opt.Service, opt.Region, opt.Project)
	cmd.Args = append(cmd.Args, "--update-tags", revisionTag+"="+revisionName)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (c *gcloudCommand) RemoveRevisionTag(ctx context.Context, revisionTag string) error {
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}

	cmd := updateTrafficCmd(ctx, opt.Service, opt.Region, opt.Project)
	cmd.Args = append(cmd.Args, "--remove-tags", revisionTag)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	return nil
}

func (c *gcloudCommand) Deploy(ctx context.Context, imageName string, revisionName string, useTraffic bool) error {
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}

	cmd := runDeployCmd(ctx, opt.Service, opt.Region, opt.Project)
	cmd.Args = append(cmd.Args, "--image", imageName)

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
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}

	cmd := updateTrafficCmd(ctx, opt.Service, opt.Region, opt.Project)
	cmd.Args = append(cmd.Args, "--to-latest")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update traffic to latest revision: %w", err)
	}

	return nil
}

func (c *gcloudCommand) UpdateTrafficToRevision(ctx context.Context, revisionName string) error {
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}

	cmd := updateTrafficCmd(ctx, opt.Service, opt.Region, opt.Project)
	cmd.Args = append(cmd.Args, "--to-revisions", revisionName+"=100")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update traffic to revision: %w", err)
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

func updateTrafficCmd(ctx context.Context, service string, region string, project string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "gcloud", "run", "services", "update-traffic", service,
		"--region", region,
		"--project", project,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func runDeployCmd(ctx context.Context, service string, region string, project string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "gcloud", "run", "deploy", service,
		"--project", project,
		"--region", region,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
