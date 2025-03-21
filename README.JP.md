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

### 基本コマンド

```bash
# トラフィックありで新しいリビジョンをデプロイ
dekopin deploy --image gcr.io/project/image:tag

# トラフィックなしで新しいリビジョンを作成
dekopin create-revision --image gcr.io/project/image:tag

# リビジョンにタグを割り当てる
dekopin create-tag --tag v1-0-0 --revision service-abcdef

# リビジョンからタグを削除する
dekopin remove-tag --tag v1-0-0

# 特定のリビジョンにトラフィックを切り替える
dekopin sr-deploy --revision service-abcdef

# タグにトラフィックを切り替える
dekopin st-deploy --tag v1-0-0
```

### グローバルフラグ

```
--project    GCPプロジェクトID
--region     GCPリージョン
--service    Cloud Runサービス名
--runner     ランナータイプ (github-actions, cloud-build, local)
--file, -f   設定ファイルのパス (デフォルト: dekopin.yml)
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

## ライセンス

このプロジェクトはMITライセンスの下で公開されています - 詳細はLICENSEファイルをご覧ください。 
