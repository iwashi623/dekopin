package dekopin

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func RemoveTagCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	gcloudCmd, ok := ctx.Value(gcloudCmdKey{}).(GcloudCommand)
	if !ok {
		return fmt.Errorf("gcloud command not found")
	}

	tag, err := getTagName(cmd)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	return removeTag(ctx, gcloudCmd, tag)
}

func removeTag(ctx context.Context, gc GcloudCommand, tag string) error {
	formattedTag := "tag-" + strings.ReplaceAll(tag, ".", "-")
	return gc.RemoveRevisionTag(ctx, formattedTag)
}
