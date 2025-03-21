package dekopin

import (
	"context"
	"fmt"
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
	gcloudCmd, ok := ctx.Value(gcloudCmdKey{}).(GcloudCommand)
	if !ok {
		return fmt.Errorf("gcloud command not found")
	}

	tf, err := getTagByFlag(cmd)
	if err != nil {
		return fmt.Errorf("failed to get tag flag: %w", err)
	}

	tag, err := createRevisionTagName(ctx, tf)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	rf, err := getRevisionByFlag(cmd)
	if err != nil {
		return fmt.Errorf("failed to get revision name: %w", err)
	}

	revision, err := createRevisionName(ctx, rf)
	if err != nil {
		return fmt.Errorf("failed to create revision name: %w", err)
	}

	return createTag(ctx, gcloudCmd, tag, revision)
}

func createTag(ctx context.Context, gc GcloudCommand, tag string, revisionName string) error {
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}

	// revisionが存在するか確認する
	client, err := run.NewRevisionsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create run client: %w", err)
	}
	defer client.Close()

	fullRevisionName := fmt.Sprintf(REVISION_FULL_NAME_FORMAT, opt.Project, opt.Region, opt.Service, revisionName)
	_, err = client.GetRevision(ctx, &runpb.GetRevisionRequest{
		Name: fullRevisionName,
	})
	if err != nil {
		return fmt.Errorf("failed to get revision: revisionName: %s is not found, error: %w", fullRevisionName, err)
	}

	formattedTag := "tag-" + strings.ReplaceAll(tag, ".", "-")
	return gc.CreateRevisionTag(ctx, formattedTag, revisionName)
}
