package dekopin

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func CreateRevision(cmd *cobra.Command, args []string) error {
	image, err := cmd.Flags().GetString("image")
	if err != nil {
		return fmt.Errorf("imageフラグの取得に失敗しました: %w", err)
	}
	if image == "" {
		return fmt.Errorf("imageフラグが指定されていません")
	}
	return createRevision(cmd, image)
}

func createRevision(cmd *cobra.Command, image string) error {
	// コンテキストを取得
	ctx := cmd.Context()

	fmt.Println("新しいCloud Runリビジョンの作成を開始します...")

	gcloudCmd := exec.CommandContext(ctx, "gcloud", "run", "deploy", config.Service,
		"--image", image,
		"--project", config.Project,
		"--region", config.Region,
		"--no-traffic", // 重要: 新しいリビジョンにトラフィックを流さない
	)

	gcloudCmd.Stdout = os.Stdout
	gcloudCmd.Stderr = os.Stderr

	// コマンド実行
	if err := gcloudCmd.Run(); err != nil {
		return fmt.Errorf("Cloud Runへのデプロイに失敗しました: %w", err)
	}

	fmt.Println("新しいリビジョンが正常にデプロイされました")

	return nil
}
