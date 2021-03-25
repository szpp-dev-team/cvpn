# cvpn-go

 `cvpn` は静岡大学情報学部 VPN サービスのコマンドラインアプリケーションです。VPN サービスをコマンドラインを通じて利用することができます。

## Documentation

### Installation

1. [release ページ](https://github.com/szpp-dev-team/cvpn/releases) から適切なファイルをダウンロードし、展開してください。  

2. バイナリファイルを任意のディレクトリに置いて `PATH` を設定してください。推奨のディレクトリのパスは `$HOME/cvpn/bin` です。

#### Recommended(Windows/Linux/MacOS)

1. go 言語のコンパイラを [ここ](https://golang.org/doc/install) からインストールする  

2. cvpn をインストール

```console
$ go get github.com/Shizuoka-Univ-dev/cvpn/cmd/cvpn
```

#### Binary

##### Windows

TODO

##### Linux

```console
$ ls
cvpn_linux_amd64 があることを確認

$ mkdir -p $HOME/cvpn/bin
$ cp cvpn_linux_amd64 $HOME/cvpn/bin/cvpn
$ echo 'export PATH="$PATH:$HOME/cvpn/bin"' >> $HOME/.profile
$ source $HOME/.profile
```

##### MacOS

```console
$ ls
cvpn_darwin_amd64 があることを確認

$ mkdir -p $HOME/cvpn/bin
$ cp cvpn_darwin_amd64 $HOME/cvpn/bin/cvpn
$ echo 'export PATH="$PATH:$HOME/cvpn/bin"' >> $HOME/.profile
$ source $HOME/.profile
```

### Usage

　一番最初に `login` を行ってください。`login` ではユーザー ID とパスワードを各 OS の設定ディレクトリ上に保存します。

```console
$ cvpn login
username >> cs200xx
password >> your_password
.
.
```

#### List

```console
$ cvpn ls

example
$ cvpn ls
```

#### Download

```console
$ cvpn download {target_file_path} -o {save_path} -v {volume}

example
$ cvpn download /report/hoge.txt
```

`{volume}` 上の `{target_path}` をダウンロードし、 `{save_path}` に保存します。  
オプションの指定がない場合、`{target_path}` はカレントディレクトリ、`{volume}` は `FSshare` がデフォルト値として設定されます。

#### Upload

現在開発中です。

```console
$ cvpn upload

example
$ cvpn upload
```

### Development

原則として devcontainer 上で開発してください。

#### Requirements

+ [Docker](https://www.docker.com/get-started)
+ [VScode](https://code.visualstudio.com/download)

#### Steps

1. VScode の拡張機能 `Remote Development` をインストールする。  
2. `cvpn-go/` ディレクトリを VScode で開き、左下の青いボタンをクリックし、`Remote-Containers: Reopen in Container` を選択する。※ここで環境がホストからコンテナ上に切り替わるので注意。  
3. コマンドパレットを開いて `Go: Install/Update Tools` と入力して、全てのツールをチェックしてインストールする。  
4. (任意) git の ssh 設定？

#### Build

```console
$ make windows
$ make linux
$ make darwin
```

#### Direcotry Structure

```console
cmd/
  └ cvpn
      └ main.go     # エントリーポイント(main 関数だけ)
  
api/
  ├ common.go       # api の共通部分(構造体とかリクエストとか)
  ├ auth.go         # auth api
  ├ download.go     # download api
  └ list.go         # list api

pkg/
  ├ config/         # config 関係
  |   └ config.go   
  └ util/           # ユーティリティ(入力とか)
      └ input.go
```
