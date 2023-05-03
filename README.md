# Go + gRPCを使ったサンプルプロジェクト

## 目次

1. [開発環境](#開発環境)
2. [インストール](#インストール)
3. [使用方法](#使用方法)
4. [開発手順](#開発手順)
5. [おまけ](#おまけ)


## 開発環境

- Go: 1.20.3
- protoc-gen-go: 1.30.0
- protoc: 3.21.12
- Docker: 20.10.23
- Docker Compose: v2.15.1
- MySQL:8.0.32

## インストール

### Go

[公式サイト](https://go.dev/dl/)より対象のバージョンをインストールします。

### protoc-gen-go

以下のコマンドを実行してインストールします。

```bash
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
```

### protoc

brewを使用してインストールします。

```bash
$ brew install protobuf
```

#### brewのインストール

[公式サイト](https://brew.sh/index_ja)を参照してインストールします。

### Docker, Docker Compose

[公式サイト](https://docs.docker.com/engine/install/)よりDockerをインストールします。

## 使用方法

### アプリケーションの起動

以下コマンドでイメージのビルドとコンテナを起動します。

```bash
$ docker compose up -d --build
```

次にマイグレーションを実行してテーブルを作成します。

```bash
$ docker compose exec backend go run cmd/main.go migration
```

コンテナの立ち上げとマイグレーションが完了しましたら  
以下のコマンドでアプリケーションを実行します。

```bash
$ docker compose exec backend go run main.go
```

アプリケーションが正常に起動した場合、以下のようなログが確認できます。

```log
2023/05/03 07:07:51 listening server with port:8080
```

アプリケーションを停止するときは`ctrl + c`を実行します。

## 開発手順

開発の流れは大体以下の流れになります。

1. `.proto`ファイルを編集する
2. `.proto`ファイルをコンパイルする
3. 必要に応じて`.go`ファイルを編集する

### 1. `.proto`ファイルを編集する

`./proto/backend.proto`ファイルを編集します。

### 2. `.proto`ファイルをコンパイルする

以下のコマンドを実行すると`.proto`ファイルがコンパイルされます。

```bash
$ protoc -I ./proto \
  --go_out=./pb --go_opt=paths=source_relative \
  --go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
  ./proto/backend.proto
```

毎回このコマンドを書くのは面倒ですので`makefile`に記載しています。
以下のように短縮できます。

```bash
$ make proto
```

### 3. 必要に応じて`.go`ファイルを編集する

必要に応じて`.go`ファイルを編集します。実際には以下のファイルを変更することが多いと思います。

- `./database/querier.go`
- `./database/query.go`
- `./server/server.go`

## おまけ

### マイグレーション

データベースを使った開発ではあらかじめ特定のデータを登録しておきたい場合や
テーブルを再作成したくなる場合があると思います。
Dockerのエコシステムを利用する手もありますが、今回はGoでマイグレーションを実装しました。

`./cmd/migration/migration.go`に全テーブルをリセットする例を用意しました。  

以下のコマンドで全テーブルをリセットできます。

```bash
$ go run cmd/main.go migration
```

### `cmd`ディレクトリ

[cobra](https://github.com/spf13/cobra)を使って`cmd`ディレクトリに`migration`コマンドを実装しています。


**実行例**

```bash
$ go run cmd/main.go ANY_COMMAND
```

LaravelやRailsにあるようなコマンドを自分自身で実装できます。
