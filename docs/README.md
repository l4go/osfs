# l4go/osfs ライブラリ

OSのファイルシステムを、1つのfs.FSとして提供するライブラリです。  
Windowsでも動作するので、OS間でのコード互換性も提供します。

## 標準ライブラリの問題

以下のコードは、Unix系OSではすべてのファイルにアクセスできるfs.FSとして動きますが、Windowsの絶対パスにはドライブ名が必要なので、そのようには動きません。

``` go
os_root := os.DirFS("/")
```

Unix系OSでの`os.DirFS("/")`相当の機能を、Windows<b>でも</b>実現できるようしたのが、このライブラリの`osfs.OsRootFS`です。

## 使い方

使い方は以下の通り。

``` go
os_root := osfs.OsRootFS
```

`os.DirFS()`のように、作成時にディレクトリは指定できません。  
サブディレクトリをfs.FSに変換したいときは、[fs.Sub()](https://pkg.go.dev/io/fs#Sub)を利用してください。

## ドライブ名(ドライブレター)の扱い

`osfs.OsRootFS`では、ルートディレクトリは、有効なドライブをサブディレクトリとして見える仮想的なディレクトリとして動作します。  
WindowsのUNCパス(`\\`で始まるパス)のドライブには対応していません。

### サンプルコード

以下のコードをWindowsで実行すると、現在有効なドライブ名をリストアップします。

``` go
func main() {
	rf, err := osfs.OsRootFS.Open(".")
	if err != nil {
		return
	}
	defer rf.Close()

	rdf, ok := rf.(fs.ReadDirFile)
	if !ok {
		fmt.Println(err)
		return
	}

	dent, err := rdf.ReadDir(-1)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, d := range dent {
        fmt.Println(d.Name())
	}
}
```
