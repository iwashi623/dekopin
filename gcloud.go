package dekopin

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
)

type gcloudCmdKey struct{}

func SetGcloudCommand(ctx context.Context, gcloudCmd GcloudCommand) context.Context {
	return context.WithValue(ctx, gcloudCmdKey{}, gcloudCmd)
}

func GetGcloudCommand(ctx context.Context) (GcloudCommand, error) {
	gcloudCmd, ok := ctx.Value(gcloudCmdKey{}).(GcloudCommand)
	if !ok {
		return nil, fmt.Errorf("gcloud command not found")
	}
	return gcloudCmd, nil
}

type GcloudCommand interface {
	CreateRevision(ctx context.Context, imageName string, commitHash string) error          // Create a revision
	CreateRevisionTag(ctx context.Context, revisionTag string, revisionName string) error   // Assign a tag to a revision
	RemoveRevisionTag(ctx context.Context, revisionTag string) error                        // Remove a tag from a revision
	Deploy(ctx context.Context, imageName string, commitHash string, useTraffic bool) error // Deploy a revision
	UpdateTrafficToLatestRevision(ctx context.Context) error                                // Update traffic to the latest revision
	UpdateTrafficToRevision(ctx context.Context, revisionName string) error                 // Update traffic to the specified revision
	UpdateTrafficToRevisionTag(ctx context.Context, tag string) error                       // Update traffic to the specified tag
	DeployWithTraffic(ctx context.Context, imageName string, commitHash string) error       // Deploy with traffic
	GetActiveRevisionTags(ctx context.Context) ([]string, error)                            // Get active revision tags
	GetRevision(ctx context.Context, revisionName string) (*runpb.Revision, error)          // Get a revision
}

type gcloudCommand struct {
	ServicesClient  *run.ServicesClient
	RevisionsClient *run.RevisionsClient

	Stdout io.Writer
	Stderr io.Writer
}

var _ GcloudCommand = &gcloudCommand{}

const (
	SERVICE_FULL_NAME_FORMAT  = "projects/%s/locations/%s/services/%s"
	REVISION_FULL_NAME_FORMAT = "projects/%s/locations/%s/services/%s/revisions/%s"
)

func NewGcloudCommand(stdout io.Writer, stderr io.Writer, servicesClient *run.ServicesClient, revisionsClient *run.RevisionsClient) GcloudCommand {
	return &gcloudCommand{
		ServicesClient:  servicesClient,
		RevisionsClient: revisionsClient,
		Stdout:          stdout,
		Stderr:          stderr,
	}
}

func (c *gcloudCommand) GetRevision(ctx context.Context, revisionName string) (*runpb.Revision, error) {
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cmdOption: %w", err)
	}

	fullRevisionName := fmt.Sprintf(REVISION_FULL_NAME_FORMAT, opt.Project, opt.Region, opt.Service, revisionName)
	revision, err := c.RevisionsClient.GetRevision(ctx, &runpb.GetRevisionRequest{
		Name: fullRevisionName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get revision: revisionName: %s is not found, error: %w", fullRevisionName, err)
	}

	return revision, nil
}

func (c *gcloudCommand) GetActiveRevisionTags(ctx context.Context) ([]string, error) {
	tagNames := []string{}
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cmdOption: %w", err)
	}

	service, err := c.ServicesClient.GetService(ctx, &runpb.GetServiceRequest{
		Name: fmt.Sprintf(SERVICE_FULL_NAME_FORMAT, opt.Project, opt.Region, opt.Service),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	for _, tag := range service.Traffic {
		if tag.Tag == "" {
			continue
		}

		tagNames = append(tagNames, tag.Tag)
	}

	return tagNames, nil
}

func (c *gcloudCommand) CreateRevision(ctx context.Context, imageName string, commitHash string) error {
	if err := c.Deploy(ctx, imageName, commitHash, false); err != nil {
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

func (c *gcloudCommand) Deploy(ctx context.Context, imageName string, commitHash string, useTraffic bool) error {
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}

	cmd := runDeployCmd(ctx, opt.Service, opt.Region, opt.Project)
	cmd.Args = append(cmd.Args, "--image", imageName)

	if commitHash != "" {
		cmd.Args = append(cmd.Args, "--revision-suffix", commitHash)
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

func (c *gcloudCommand) UpdateTrafficToRevisionTag(ctx context.Context, tag string) error {
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}

	cmd := updateTrafficCmd(ctx, opt.Service, opt.Region, opt.Project)
	cmd.Args = append(cmd.Args, "--to-tags", tag+"=100")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update traffic to revision tag: %w", err)
	}

	return nil
}

func (c *gcloudCommand) DeployWithTraffic(ctx context.Context, imageName string, commitHash string) error {
	if err := c.Deploy(ctx, imageName, commitHash, true); err != nil {
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
