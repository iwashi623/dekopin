package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func Run(ctx context.Context) (int, error) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return 1, err
	}
	return 0, nil
}

var rootCmd = &cobra.Command{
	Use:   "dekopin",
	Short: "Dekopin is a Cloud Run deployment tool",
	Long:  "Dekopin is a tool to deploy Cloud Run services with tags and traffic management.",
}

var createRevisionCmd = &cobra.Command{
	Use:   "create_revision",
	Short: "Create a new Cloud Run revision",
	Run: func(cmd *cobra.Command, args []string) {
		// Cloud Run API を呼び出してリビジョン作成処理を実装
		fmt.Println("Creating a new revision...")
		// TODO: ビルド、プッシュ、デプロイ処理
	},
}

var createTagCmd = &cobra.Command{
	Use:   "create_tag",
	Short: "Assign a tag to a Cloud Run revision",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Assigning tag to revision...")
		// TODO: リビジョン ID とタグ名を引数・フラグから取得し、タグ付与処理を実装
	},
}

var removeTagCmd = &cobra.Command{
	Use:   "remove_tag",
	Short: "Remove a tag from a Cloud Run revision",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Removing tag from revision...")
		// TODO: リビジョン ID と対象タグを指定してタグ削除処理を実装
	},
}

var switchTrafficCmd = &cobra.Command{
	Use:   "switch_traffic",
	Short: "Switch traffic to a specified Cloud Run revision",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Switching traffic...")
		// TODO: 対象リビジョンとトラフィック割合を引数・フラグから取得し、トラフィック制御を実施
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy new revision with tag and traffic management",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deploying new revision with tag management and traffic switching...")
		// TODO: 必要なサブコマンド処理（create_tag, remove_tag, switch_traffic）を内部的に実施
	},
}

func init() {
	rootCmd.AddCommand(createRevisionCmd)
	rootCmd.AddCommand(createTagCmd)
	rootCmd.AddCommand(removeTagCmd)
	rootCmd.AddCommand(switchTrafficCmd)
	rootCmd.AddCommand(deployCmd)
}
