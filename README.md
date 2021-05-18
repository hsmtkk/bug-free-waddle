en-photo thumbnails

# Usage

```shell
export ENPHOTO_USERNAME=xxxx
export ENPHOTO_PASSWORD=xxxx
go run main.go html https://en-photo.net/albums/xxxx > album.html
go run main.go datasrc album.html > datasrc.txt
go run main.go resolve datasrc.txt > thumbnail.txt
go run main.go thumbnail thumbnail.txt
```

# 解説

* html

アルバムページのHTMLを取得する。

* datasrc

imgタグ、data-src属性を抽出する。

* resolve

上記で抽出したパスにアクセスし、サムネイルの実URLを得る。

* thumbnail

サムネイル画像をダウンロードする。
