# cvpn-go

 `cvpn` は静岡大学情報学部 VPN サービスのコマンドラインアプリケーションです。VPN サービスをコマンドラインを通じて利用することができます。

## Documentation

### Installation

インストール方法は以下の2種類があります。推奨する方法は `go install` です。

- `go install` によるビルド & インストール 【推奨】
- release ページからビルド済みバイナリをダウンロード

#### インストール方法A: go get 【推奨】

1. go 言語のコンパイラを [ここ](https://golang.org/doc/install) からインストールする

2. パスを通す

```console
# linux
$ echo 'export PATH="$PATH:$HOME/go/bin"' >> $HOME/.profile

# mac
$ echo 'export PATH="$HOME/go/bin:$PATH"' >> $HOME/.profile
```

3. cvpn をインストール

```console
$ go install github.com/Shizuoka-Univ-dev/cvpn/cmd/cvpn@latest
$ source $HOME/.profile
$ cvpn
cvpn is a tool which makes you happy
.
.
```

#### インストール方法B: ビルド済みバイナリのダウンロード

1. [release ページ](https://github.com/szpp-dev-team/cvpn/releases) から適切なファイルをダウンロードし、展開してください。  

2. バイナリファイルを任意のディレクトリに置いて `PATH` を設定してください。推奨のディレクトリのパスは `CVPN_PATH=$HOME/cvpn/bin` です。  
以下に設定手順の例を示します。

##### Windows

TODO

##### Linux & MacOS

```console
$ cd バイナリファイルをダウンロードしたディレクトリ
$ ls
cvpn があることを確認

$ CVPN_PATH=$HOME/cvpn/bin
$ mkdir -p $CVPN_PATH
$ cp cvpn $CVPN_PATH
$ echo 'export PATH=$PATH:$CVPN_PATH' >> $HOME/.profile
$ source $HOME/.profile
$ cvpn
cvpn is a tool which makes you happy
.
.
```

### Usage

> Note: パスやオプションの指定が `{...}` となっていますが、このとき `{}` をつけて入力する必要はありません。

　一番最初に `login` を行ってください。`login` ではユーザー ID とパスワードを各 OS の設定ディレクトリ上に保存します。  
また、途中で作成を確認するメッセージが表示されますが、設定ファイルの作成を許可する場合は `y`, 許可しない場合は `n` を入力してください。

```console
$ cvpn login
username >> cs200xx
password >> your_password
.
.
```

#### List

`cvpn ls` コマンドは指定した `{dir_path}` 上のファイルとディレクトリを一覧表示します。  

```console
$ cvpn ls {dir_path} -v {volume} --path --json

example
$ cvpn ls /report -v fsshare
.
.
```

絵文字を使用したりしているため、ターミナルのフォントを [Nerd Font](https://www.nerdfonts.com/) にすることを推奨します。  
推奨フォントは [JetBrainsMono Nerd Font](https://github.com/ryanoasis/nerd-fonts/releases/download/v2.1.0/JetBrainsMono.zip) です。

- Options
  
  - `--path`: ファイル or ディレクトリのパスを表示します。  
  
  - `-v {volume}`: 参照するファイルが存在するボリュームを指定します(`fsshare`, `fs/{dir}`)。デフォルト値は `fsshare` です。

  - `--json`: json 形式で表示します。

> Note: volume は FSShare や FS などのことを示します。

#### Download

`cvpn download` コマンドは指定した `{target_file_path}` をダウンロードします。

```console
$ cvpn download {target_file_path} -o {save_path} -v {volume}

example
$ cvpn download /cs200xx/I_am_file.txt -o ./univ -v fs/2020
```

- Options

  - `-o {save_path}`: ダウンロードしたファイルの保存先を指定します。**必ずディレクトリのパスを指定してください**(仕様変更予定)
  
  - `-v {volume}`: 参照するファイルが存在するボリュームを指定します(`fsshare`, `fs/{dir}`)。デフォルト値は `fsshare` です。

#### Upload

`cvpn upload` コマンドは `{source_file_path}` ファイルを `{dst_path}` 上にアップロードします。

```console
$ cvpn upload {source_file_path} {dst_path} -v {volume}

example
$ cvpn upload text.txt /cs200xx -v fs/2020
```

- Options
  
  - `-v {volume}`: 参照するファイルが存在するボリュームを指定します(`fsshare`, `fs/{dir}`)。デフォルト値は `fsshare` です。

#### Find

`cvpn find` コマンドは `{starting-directory}` ディレクトリを元にして検索を行います。

```console
$ cvpn find {starting-directory} -v {volume} -name {name_pattern} -path {path_pattern} -r

example
$ cvpn find /cs200xx -v fs/2020 -name hogehoge
```

- Options
  
  - `-v {volume}`: 参照するファイルが存在するボリュームを指定します(`fsshare`, `fs/{dir}`)。デフォルト値は `fsshare` です。

  - `-r`, `--recursive`: 再帰的にファイルを探索します。

  - `--name {name_pattern}`: 検索するディレクトリ or ファイルの **名前** のパターンを指定します。[正規表現](https://github.com/google/re2/wiki/Syntax)が使用可能です。

  - `--path {path_pattern}`: 検索するディレクトリ or ファイルの **パス** のパターンを指定します。[正規表現](https://github.com/google/re2/wiki/Syntax)が使用可能です。

> Note:  
> **名前** は **パス** のうち末尾の部分を指します。
> 例えば、`/class/english/hoge.pdf` において、名前とパスは以下のようになります。

|k|v|
|----|----|
|名前|`hoge.pdf`|
|パス|`/class/english/hoge.pdf`|

### Log

ログは `$USER_CACHE_DIR/cvpn/log/` 配下に保存されます。  
`$USER_CACHE_DIR` は OS によって異なるので、詳しい内容は https://golang.org/pkg/os/#UserCacheDir を参照してください。

### Development

開発者向けの内容です。

#### Build

```console
$ make windows
$ make linux
$ make darwin
```
