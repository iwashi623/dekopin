package dekopin

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var config DekopinConfig

func setConfig(ctx context.Context) error {
	// 設定ファイルを読み込む
	dekopinYaml, err := os.ReadFile("dekopin.yaml")
	if err != nil {
		return fmt.Errorf("dekopin.yamlの読み込みに失敗しました: %w", err)
	}

	// 設定ファイルをパースする
	if err := yaml.Unmarshal(dekopinYaml, &config); err != nil {
		return fmt.Errorf("dekopin.yamlのパースに失敗しました: %w", err)
	}

	return nil
}

type DekopinConfig struct {
	Project string `yaml:"project"`
	Region  string `yaml:"region"`
	Service string `yaml:"service"`
}

func Run(ctx context.Context) (int, error) {
	if err := setConfig(ctx); err != nil {
		fmt.Println(err)
		return 1, err
	}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
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
	Use:   "create-revision",
	Short: "Create a new Cloud Run revision",
	Args:  cobra.NoArgs,
	RunE:  StartCreateRevision,
}

var createTagCmd = &cobra.Command{
	Use:   "create-tag",
	Short: "Assign a tag to a Cloud Run revision",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Assigning tag to revision...")
		// TODO: リビジョン ID とタグ名を引数・フラグから取得し、タグ付与処理を実装
	},
}

var removeTagCmd = &cobra.Command{
	Use:   "remove-tag",
	Short: "Remove a tag from a Cloud Run revision",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Removing tag from revision...")
		// TODO: リビジョン ID と対象タグを指定してタグ削除処理を実装
	},
}

var switchTrafficCmd = &cobra.Command{
	Use:   "switch-traffic",
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
		// TODO: 必要なサブコマンド処理（create-tag, remove-tag, switch-traffic）を内部的に実施
	},
}

func init() {
	rootCmd.AddCommand(createRevisionCmd)
	createRevisionCmd.Flags().StringP("image", "i", "", "コンテナイメージのURL")
	createRevisionCmd.MarkFlagRequired("image")

	rootCmd.AddCommand(createTagCmd)
	rootCmd.AddCommand(removeTagCmd)
	rootCmd.AddCommand(switchTrafficCmd)
	rootCmd.AddCommand(deployCmd)
}
