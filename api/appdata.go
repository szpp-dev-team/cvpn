package api

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

// ConfigDirPath は cvpn の設定データを格納するユーザディレクトリの絶対パスを返す。
// Linux ならおそらく一般的に ~/.config 。
/*
func configDirPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(configDir, "cvpn"), nil
}
*/

// cacheDirPath は cvpn のキャッシュデータを格納するユーザディレクトリの絶対パスを返す。
// Linux ならおそらく一般的に ~/.cache 。
func cacheDirPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return path.Join(cacheDir, "cvpn"), nil
}

// cookieCacheFilePath はクッキーを保存するファイルの絶対パスを返す。
func cookieFilePath() (string, error) {
	cacheDir, err := cacheDirPath()
	if err != nil {
		return "", err
	}
	return path.Join(cacheDir, "cookies.txt"), nil
}

// saveCookies は cookies を cookieFilePath() で示されるパスのファイルへ書き込む。
// ファイルはまず空になったあとで書き込まれる。
// 保存形式や保存先ファイル名は規定しない。読み出しは loadCookies を用いること。
func saveCookies(cookies []string) error {
	cookieFilePath, err := cookieFilePath()
	if err != nil {
		return err
	}

	// MkdirAll は mkdir -p と同じ効果。ディレクトリが既に存在してもエラーを返さない。
	if err := os.MkdirAll(path.Dir(cookieFilePath), 0700); err != nil {
		return err
	}

	fp, err := os.OpenFile(cookieFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fp.Close()

	for _, cookie := range cookies {
		fmt.Fprintln(fp, cookie)
	}
	logger.Printf("saveCookies(): Saved %d lines of cookie into %q.\n", len(cookies), cookieFilePath)
	return nil
}

// loadCookies は saveCookies() で書き出されたクッキー群を読み出して返す。
func loadCookies() ([]string, error) {
	cookieFilePath, err := cookieFilePath()
	if err != nil {
		return nil, err
	}

	fp, err := os.Open(cookieFilePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	cookies := make([]string, 0, 5)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		cookie := scanner.Text()
		cookies = append(cookies, cookie)
	}

	logger.Printf("loadCookies(): Loaded %d lines of cookie.\n", len(cookies))
	return cookies, nil
}

func deleteCookieFile() error {
	cookieFilePath, err := cookieFilePath()
	if err != nil {
		return err
	}
	logger.Printf("deleteCookieFile(): delete.")
	return os.Remove(cookieFilePath)
}