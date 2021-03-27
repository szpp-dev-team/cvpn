# cvpn-go

 `cvpn` は静岡大学情報学部 VPN サービスのコマンドラインアプリケーションです。VPN サービスをコマンドラインを通じて利用することができます。

## Documentation

### Installation

#### 【Recommend】go get (Windows/Linux/MacOS)

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
$ go get -u github.com/Shizuoka-Univ-dev/cvpn/cmd/cvpn
$ source $HOME/.profile
$ cvpn
cvpn is a tool which makes you happy
.
.
```

#### Binary

1. [release ページ](https://github.com/szpp-dev-team/cvpn/releases) から適切なファイルをダウンロードし、展開してください。  

2. バイナリファイルを任意のディレクトリに置いて `PATH` を設定してください。推奨のディレクトリのパスは `CVPN_PATH=$HOME/cvpn/bin` です。

##### Windows

TODO

##### Linux & MacOS

```console
$ ls
cvpn_linux_amd64 or cvpn_darwin_amd64 があることを確認

$ CVPN_PATH=$HOME/cvpn/bin
$ mkdir -p $CVPN_PATH
$ cp cvpn_linux_amd64 $CVPN_PATH (mac: cp cvpn_darwin_amd64 $HOME/cvpn/bin/cvpn)
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
また、途中で作成を確認するプロンプトが表示されますが、設定ファイルの作成を許可する場合は `y`, 許可しない場合は `n` を入力してください。

```console
$ cvpn login
username >> cs200xx
password >> your_password
.
.
```

#### List

```console
$ cvpn ls {dir_path} -v {volume} --path

example
$ cvpn ls /report
.
.
```

`{dir_path}` 上のファイルとディレクトリを一覧表示します。  
絵文字を使用したりしているため、ターミナルのフォントを [Nerd Font](https://www.nerdfonts.com/) にすることを推奨します。  
推奨フォントは [JetBrainsMono Nerd Font](https://github.com/ryanoasis/nerd-fonts/releases/download/v2.1.0/JetBrainsMono.zip) です。

+ Options
  
  + `--path`: ファイル or ディレクトリのパスを表示します。  
  
  + `-v {volume}`: 参照するファイルが存在するボリュームを指定します(`fsshare`, `fs/{dir}`)。デフォルト値は `fsshare` です。

#### Download

```console
$ cvpn download {target_file_path} -o {save_path} -v {volume}

example
$ cvpn download /report/hoge.txt
```

`{target_file_path}` をダウンロードします。

+ Options

  + `-o {save_path}`: ダウンロードしたファイルの保存先を指定します。**必ずディレクトリのパスを指定してください**(仕様変更予定)
  
  + `-v {volume}`: 参照するファイルが存在するボリュームを指定します(`fsshare`, `fs/{dir}`)。デフォルト値は `fsshare` です。

#### Upload

現在開発中です。

### Development

#### Build

```console
$ make windows
$ make linux
$ make darwin
```
