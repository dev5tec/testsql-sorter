# testsql-sorter

テストデータ作成用SQLファイルのINSERT文を任意の順番にソートするツール

指定したディレクトリにある Test(.*).sql にマッチするファイルを再帰的に検索し、INSERT文を抽出・並べ替え・上書きします（既存のテキストデータは消去）

## Installation

```bash
go get -u github.com/dev5tec/testsql-sorter
```

## Usage

1. testsql-sorter.yml にテーブルの順番を設定する
2. 下記のコマンドを実行する

```bash
cd hoge # testsql-sorter.yml があるディレクトリに移動
testsql-sorter [directory]
```
