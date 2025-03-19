package dekopin

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func CreateTag(cmd *cobra.Command, args []string) error {
	tag, err := getTagName(cmd)
	if err != nil {
		return fmt.Errorf("tag名の取得に失敗しました: %w", err)
	}

	revision, err := getRevisionName(cmd)
	if err != nil {
		return fmt.Errorf("revision名の取得に失敗しました: %w", err)
	}

	return createTag(cmd, tag, revision)
}

func createTag(cmd *cobra.Command, tag string, revision string) error {
	ctx := cmd.Context()
	// tagのフォーマットを変換
	formattedTag := strings.ReplaceAll(tag, ".", "-")

	gcloudCmd := exec.CommandContext(ctx, "gcloud", "run", "services", "update-traffic", config.Service,
		"--region", config.Region,
		"--project", config.Project,
		"--update-tags", "tag-"+formattedTag+"="+revision,
		"--to-revisions", revision+"=0",
	)

	gcloudCmd.Stdout = os.Stdout
	gcloudCmd.Stderr = os.Stderr

	// コマンド実行
	if err := gcloudCmd.Run(); err != nil {
		return fmt.Errorf("タグの作成に失敗しました: %w", err)
	}
	return nil
}
