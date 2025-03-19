package dekopin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/spf13/cobra"
)

const (
	REVISION_FULL_NAME_FORMAT = "projects/%s/locations/%s/services/%s/revisions/%s"
)

func CreateTagCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	tag, err := getTagName(cmd)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	revision, err := getRevisionName(cmd)
	if err != nil {
		return fmt.Errorf("failed to get revision name: %w", err)
	}

	return createTag(ctx, tag, revision)
}

func createTag(ctx context.Context, tag string, revisionName string) error {
	// revisionが存在するか確認する
	client, err := run.NewRevisionsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create run client: %w", err)
	}
	defer client.Close()

	fullRevisionName := fmt.Sprintf(REVISION_FULL_NAME_FORMAT, config.Project, config.Region, config.Service, revisionName)
	_, err = client.GetRevision(ctx, &runpb.GetRevisionRequest{
		Name: fullRevisionName,
	})
	if err != nil {
		return fmt.Errorf("failed to get revision: revisionName: %s is not found, error: %w", fullRevisionName, err)
	}

	// Convert tag format
	formattedTag := "tag-" + strings.ReplaceAll(tag, ".", "-")
	gcloudCmd := exec.CommandContext(ctx, "gcloud", "run", "services", "update-traffic", config.Service,
		"--region", config.Region,
		"--project", config.Project,
		"--update-tags", formattedTag+"="+revisionName,
	)

	gcloudCmd.Stdout = os.Stdout
	gcloudCmd.Stderr = os.Stderr

	// Execute command
	if err := gcloudCmd.Run(); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}
