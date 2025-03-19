package dekopin

import (
	"fmt"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

func StartCreateRevision(cmd *cobra.Command, args []string) error {
	image, err := cmd.Flags().GetString("image")
	if err != nil {
		return fmt.Errorf("imageフラグの取得に失敗しました: %w", err)
	}
	if image == "" {
		return fmt.Errorf("imageフラグが指定されていません")
	}
	return CreateRevision(cmd, image)
}

func CreateRevision(cmd *cobra.Command, image string) error {
	// コンテキストを取得
	ctx := cmd.Context()

	fmt.Println("新しいCloud Runリビジョンの作成を開始します...")

	// Cloud Run クライアントの初期化
	client, err := run.NewServicesRESTClient(ctx, option.WithUserAgent("dekopin"))
	if err != nil {
		return fmt.Errorf("Cloud Run クライアントの初期化に失敗しました: %w", err)
	}
	defer client.Close()

	// サービス名のフルパスを構築
	parent := fmt.Sprintf("projects/%s/locations/%s", config.Project, config.Region)
	name := fmt.Sprintf("%s/services/%s", parent, config.Service)

	// 既存のサービスを取得
	getReq := &runpb.GetServiceRequest{
		Name: name,
	}

	service, err := client.GetService(ctx, getReq)
	if err != nil {
		return fmt.Errorf("サービス情報の取得に失敗しました: %w", err)
	}

	fmt.Printf("サービス情報を取得しました: %s\n", service)

	// // 新しいリビジョンのために必要な変更のみを行う
	// // 既存のテンプレートを更新
	// if len(service.Template.Containers) > 0 {
	// 	// コンテナイメージの更新
	// 	service.Template.Containers[0].Image = image

	// 	// 強制的に新しいリビジョンを作成するための設定
	// 	// クライアントIDを環境変数ではなくアノテーションとして追加
	// 	if service.Template.Annotations == nil {
	// 		service.Template.Annotations = make(map[string]string)
	// 	}
	// 	service.Template.Annotations["client.knative.dev/user-image"] = image
	// 	service.Template.Annotations["deploy-time"] = fmt.Sprintf("%d", time.Now().Unix())
	// } else {
	// 	return fmt.Errorf("サービスにコンテナが定義されていません")
	// }

	// // サービスの更新リクエスト作成
	// updateReq := &runpb.UpdateServiceRequest{
	// 	Service: service,
	// }

	// // サービスを更新して新しいリビジョンを作成
	// operation, err := client.UpdateService(ctx, updateReq)
	// if err != nil {
	// 	return fmt.Errorf("サービスの更新に失敗しました: %w", err)
	// }

	// fmt.Println("リビジョンのデプロイを開始しました。完了を待機しています...")

	// // 操作の完了を待機
	// resp, err := operation.Wait(ctx)
	// if err != nil {
	// 	return fmt.Errorf("デプロイ操作の完了待機中にエラーが発生しました: %w", err)
	// }

	// fmt.Printf("新しいリビジョンが正常にデプロイされました: %s\n", resp.Uri)
	return nil
}
