# cvpn-go
## 開発環境
docker のおかげで go のインストールも vscode の拡張機能のインストールもしないで開発環境を一発で作ることができます。  
docker のありがたみを知る & 原因不明のバグ等を防ぐため、原則としてこのコンテナ上で開発してください。  
1. VScode の拡張機能 `Remote Development` をインストールする。  
2. `cvpn-go/` ディレクトリを VScode で開き、左下の青いボタンをクリックし、`Remote-Containers: Reopen in Container` を選択する。※ここで環境がホストからコンテナ上に切り替わるので注意。  
3. コマンドパレットを開いて `Go: Install/Update Tools` と入力して、全てのツールをチェックしてインストールする。    
4. (任意) git の ssh 設定？