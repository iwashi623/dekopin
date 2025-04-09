# Dekopin

Dekopinは、Google Cloud Runのデプロイ、リビジョン、トラフィックルーティングを管理するためのコマンドラインツールです。リビジョンタグ付けとトラフィック管理をサポートしたCloud Runサービスのデプロイのための効率的なワークフローを提供します。

## 機能

- トラフィックの有無にかかわらず、新しいCloud Runリビジョンをデプロイ
- リビジョンタグの作成と管理
- リビジョン間のトラフィック切り替え
- 複数のデプロイ環境（ローカル、GitHub Actions、Cloud Build）のサポート
- YAML形式の設定
- 組み込みタイムアウト処理（デフォルト30秒）
- コミットハッシュによる一貫したリビジョン命名

## インストール

```bash
go install github.com/iwashi623/dekopin/cmd/dekopin@latest
```

## 設定

Dekopinはデフォルトで`dekopin.yml`という名前のYAML設定ファイルを使用します。設定ファイルの構造は以下の通りです：

```yaml
project: あなたのGCPプロジェクトID
region: GCPリージョン
service: あなたのCloud Runサービス名
runner: github-actions  # または: cloud-build, local
```

## 使用方法

### グローバルフラグ

以下のフラグは任意のサブコマンドで使用できます：

```
--project    GCPプロジェクトID
--region     GCPリージョン
--service    Cloud Runサービス名
--runner     ランナータイプ (github-actions, cloud-build, local)
--file, -f   設定ファイルのパス (デフォルト: dekopin.yml)
```

### タグの命名規則

Dekopinのタグは以下の規則に従う必要があります：

- 小文字の英数字とハイフン（`a-z`、`0-9`、`-`）のみで構成する必要があります
- 大文字、ピリオド、アンダースコア、スペース、特殊文字は使用できません
- 空のタグはランナータイプによって異なる扱いになります：
  - GitHub ActionsとCloud Build：参照に基づいてタグが自動生成されます
  - ローカルランナー：空のタグは許可されず、エラーになります

有効なタグの例：
- `production`
- `staging`
- `release-v1`
- `v1-0-0`
- `feature-123`

無効なタグの例：
- `Production`（大文字を含む）
- `staging.1`（ピリオドを含む）
- `test_tag`（アンダースコアを含む）
- `tag with spaces`（スペースを含む）

### サブコマンド

#### deploy

トラフィックを向けた新しいリビジョンをデプロイします。

```bash
dekopin deploy --image [イメージURL]
```

オプション：
- `--image, -i`（必須）：コンテナイメージURL（例：gcr.io/project/image:tag）
- `--tag, -t`：新しいリビジョンのタグ名（タグの命名規則に従う必要があります）
- `--create-tag`：デプロイ後にリビジョンタグを作成します
- `--remove-tags`：デプロイ前にすべてのリビジョンタグを削除します

例：
```bash
# 特定のイメージでデプロイ
dekopin deploy --image gcr.io/project/image:latest

# デプロイしてタグを作成
dekopin deploy --image gcr.io/project/image:latest --create-tag --tag release-v1
```

#### create-revision

トラフィックを向けずに新しいリビジョンを作成します。

```bash
dekopin create-revision --image [イメージURL]
```

オプション：
- `--image, -i`（必須）：コンテナイメージURL

例：
```bash
# トラフィックなしで新しいリビジョンを作成
dekopin create-revision --image gcr.io/project/image:v2
```

#### create-tag

既存のリビジョンにタグを割り当てます。タグは上記の命名規則に従う必要があります。

```bash
dekopin create-tag --tag [タグ名] --revision [リビジョン名]
```

オプション：
- `--tag, -t`：作成するタグ名（タグの命名規則に従う必要があります）
- `--revision`：タグ付けするリビジョン名（デフォルトは最新）

例：
```bash
# 最新のリビジョンにタグ付け
dekopin create-tag --tag production

# 特定のリビジョンにタグ付け
dekopin create-tag --tag staging --revision service-abcdef
```

#### remove-tag

リビジョンからタグを削除します。

```bash
dekopin remove-tag --tag [タグ名]
```

オプション：
- `--tag, -t`（必須）：削除するタグ名

例：
```bash
# タグを削除
dekopin remove-tag --tag old-release
```

#### sr-deploy（Switch Revision Deploy）

特定のリビジョンにトラフィックを切り替えます。

```bash
dekopin sr-deploy --revision [リビジョン名]
```

オプション：
- `--revision`：トラフィックを向けるリビジョン名

例：
```bash
# 特定のリビジョンにすべてのトラフィックを向ける
dekopin sr-deploy --revision service-abcdef
```

#### st-deploy（Switch Tag Deploy）

特定のタグを持つリビジョンにトラフィックを切り替えます。タグはすでに存在し、タグの命名規則に従っている必要があります。

```bash
dekopin st-deploy --tag [タグ名]
```

オプション：
- `--tag, -t`（必須）：トラフィックを向けるタグ名
- `--remove-tags`：デプロイ対象のリビジョンタグを除くすべてのリビジョンタグを削除します

例：
```bash
# タグ付けされたリビジョンにすべてのトラフィックを向ける
dekopin st-deploy --tag production

# タグ付けされたリビジョンに切り替えて他のタグをクリーンアップ
dekopin st-deploy --tag production --remove-tags
```

## CI/CD統合

### GitHub Actions

DekopinはGitHub Actions環境を自動的に検出し、タグ名やコミットハッシュに環境変数を使用できます。コミットハッシュの最初の7文字がリビジョン名に使用されます。

ワークフロー例：

```yaml
name: Deploy to Cloud Run

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
          
      - name: Install Dekopin
        run: go install github.com/iwashi623/dekopin/cmd/dekopin@latest
        
      - name: Deploy to Cloud Run
        run: dekopin deploy --image gcr.io/project/image:${{ github.sha }}
```

### Google Cloud Build

Dekopinはビルド環境変数を使用したCloud Build統合もサポートしています。コミットハッシュが7文字以下の場合はそのまま使用され、それ以上の場合は最初の7文字が使用されます。

`cloudbuild.yaml`の例：

```yaml
steps:
  - name: 'golang'
    entrypoint: 'go'
    args: ['install', 'github.com/iwashi623/dekopin/cmd/dekopin@latest']
  
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/image:$COMMIT_SHA', '.']
  
  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/$PROJECT_ID/image:$COMMIT_SHA']
  
  - name: 'golang'
    entrypoint: 'dekopin'
    args: ['deploy', '--image', 'gcr.io/$PROJECT_ID/image:$COMMIT_SHA']
```

## バリデーション

Dekopinには様々な入力値のバリデーションが含まれています：

- タグは小文字の英数字とハイフンで構成される必要があります（例：`release-v1`、`v1-0-0`）
- コマンドには適切な必須フラグがあります
- 入力値は実行前に検証されます

## トラブルシューティング

### 一般的なエラー

- **タイムアウトエラー**：デフォルトでは、Dekopinのタイムアウトは30秒です。長時間実行される操作では、コード内でこの値を増やすことを検討してください。
- **接続エラー**：「client connection is closing」などのエラーが表示される場合は、API呼び出し中にクライアント接続が開いたままになっていることを確認してください。
- **タグフォーマットエラー**：無効なタグフォーマットに関するエラーが発生した場合は、タグが命名規則（小文字の英数字とハイフンのみ）に従っていることを確認してください。

## ライセンス

このプロジェクトはMITライセンスの下で公開されています - 詳細はLICENSEファイルをご覧ください。 
